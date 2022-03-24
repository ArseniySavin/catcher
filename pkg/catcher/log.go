package catcher

import (
	"fmt"
	"github.com/ArseniySavin/gocatcher/pkg/catcher/internal"
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
	log.Println(internal2.Marshal(&internal2.LogMsg{
		Level:   "INFO",
		Host:    internal2.GetHost(),
		Payload: msg,
	}))
}

// LogTrace for trace data
func LogTrace(msg string, spot interface{}) {
	if Mode {
		spotMsg := ""
		switch spot.(type) {
		case interface{}:
			spotMsg = internal2.MarshalStruct(spot)
		case int:
			spotMsg = fmt.Sprintf("%+d", spot)
		default:
			spotMsg = fmt.Sprintf("%s", spot)
		}

		log.Println(internal2.Marshal(&internal2.LogMsg{
			Level:   "TRACE",
			Host:    internal2.GetHost(),
			Payload: msg + ", " + spotMsg,
		}))
	}
}

// LogFatal stop app use os.Exit(1)
func LogFatal(err error) {
	log.Fatalln(internal2.Marshal(&internal2.LogMsg{
		Level:   "FATAL",
		Host:    internal2.GetHost(),
		Payload: err.Error(),
	}))
}

// LogError print error with call data
func LogError(err error) {
	call := internal2.CallInfo(2)
	log.Println(internal2.Marshal(&internal2.LogMsg{
		Level:   "ERROR",
		Host:    internal2.GetHost(),
		Payload: err.Error() + ", " + internal2.MarshalStruct(call),
	}))
}
