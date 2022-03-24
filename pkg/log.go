package catcher

import (
	"fmt"
	"github.com/ArseniySavin/catcher/pkg/catcher/internal"
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
			Payload: msg + ", " + spotMsg,
		}))
	}
}

// LogFatal stop app use os.Exit(1)
func LogFatal(err error) {
	log.Fatalln(internal.Marshal(&internal.LogMsg{
		Level:   "FATAL",
		Host:    internal.GetHost(),
		Payload: err.Error(),
	}))
}

// LogError print error with call data
func LogError(err error) {
	call := internal.CallInfo(2)
	log.Println(internal.Marshal(&internal.LogMsg{
		Level:   "ERROR",
		Host:    internal.GetHost(),
		Payload: err.Error() + ", " + internal.MarshalStruct(call),
	}))
}
