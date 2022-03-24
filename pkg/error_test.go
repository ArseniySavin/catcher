package pkg

import (
	"errors"
	"testing"
)

func Test_Info(t *testing.T) {
	LogInfo("Exsample!")
}

func Test_Trace(t *testing.T) {
	Mode = Tracing
	LogTrace("Debug. It is debug info", struct {
		Id   int
		Name string
	}{
		Id:   100,
		Name: "test",
	})
}

type Exsample struct {
	Id   int
	Name string
}

func Test_Trace_struct(t *testing.T) {
	Mode = Tracing
	LogTrace("Debug. It is debug info", Exsample{
		Id:   102,
		Name: "Debug test",
	})
}

func Test_Trace_error_struct(t *testing.T) {
	Mode = Tracing
	LogTrace("Trace base error", BaseError.NewCode("-999").Throw("Base error"))
}

func Test_Error(t *testing.T) {
	Mode = Tracing
	LogError(BaseError.NewCode("500").Throw("Example error"))
}

func Test_Error_Simple(t *testing.T) {
	Mode = Tracing
	LogError(errors.New("Simple error"))
}

//func Test_Fatal(t *testing.T) {
//	Mode = Tracing
//	LogFatal(errors.New("Simple error"))
//}
