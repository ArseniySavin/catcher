package pkg

import (
	"errors"
	"fmt"
	"github.com/ArseniySavin/catcher/pkg/internal"
	"log"
)

type mode bool

// Global
var (
	Mode mode = Regular
)

// Global
const (
	Tracing mode = true
	Regular mode = false
)

// LogInfo for positive events
func LogInfo(msg string) {
	log.Println(internal.Marshal(&internal.LogMsg{
		Level:   "INFO",
		Host:    internal.GetHost(),
		Message: msg,
		Payload: msg,
	}))
}

// LogTrace for trace data
func LogTrace(msg string, spot interface{}) {
	if Mode {
		spotMsg := ""
		switch spot.(type) {
		case interface{}:
			spotMsg = internal.MarshalStruct(spot)
		case int:
			spotMsg = fmt.Sprintf("%+d", spot)
		default:
			spotMsg = fmt.Sprintf("%s", spot)
		}

		log.Println(internal.Marshal(&internal.LogMsg{
			Level:   "TRACE",
			Host:    internal.GetHost(),
			Message: msg,
			Payload: msg + ", " + spotMsg,
		}))
	}
}

// LogFatal stop app use os.Exit(1)
func LogFatal(err error) {
	errStr := err.Error()
	log.Fatalln(internal.Marshal(&internal.LogMsg{
		Level:   "FATAL",
		Host:    internal.GetHost(),
		Message: errStr,
		Payload: errStr,
	}))
}

// LogError print error with call data
func LogError(err error) {
	errStr := err.Error()

	if errors.Is(err, BaseError) {
		baseErr := BaseError
		errors.As(err, &baseErr)
		errStr = fmt.Sprintf("%s %s", baseErr.Error(), baseErr.Stack())
	}

	call := internal.CallInfo(2)
	log.Println(internal.Marshal(&internal.LogMsg{
		Level:   "ERROR",
		Host:    internal.GetHost(),
		Message: err.Error(),
		Payload: errStr + ", " + internal.MarshalStruct(call),
	}))
}
