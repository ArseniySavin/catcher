package pkg

import (
	"io"
	"log/slog"
	"testing"
)

// go test -timeout 30s -run ^Test_slog_Info$ github.com/ArseniySavin/catcher/pkg -v
func Test_slog_Info(t *testing.T) {
	l := slog.New(NewCatcherHandler(nil))
	l.Info("Exsample!", "user", "TEST")

	// 2025/01/08 13:59:57 {"Level":"INFO","Host":"nuc","Message":"Exsample!","Payload":["user:TEST"]}
}

func Test_slog_Error(t *testing.T) {
	l := slog.New(NewCatcherHandler(&slog.HandlerOptions{Level: slog.LevelWarn}))

	l.Warn("Exsample!", "test", struct{ m string }{m: "Test!"}, "test2", struct{ m string }{m: "Test2!"})
	l.Error("Exsample!", ErrorKey, io.ErrNoProgress)

	// 2025/01/08 14:01:22 {"Level":"WARN","Host":"nuc","Message":"Exsample!","Payload":["test:{Test!}","test2:{Test2!}"]}
	// 2025/01/08 14:01:22 {"Level":"ERROR","Host":"nuc","Message":"Exsample!","Payload":["Source:github.com/ArseniySavin/catcher/pkg/slog_test.go .Test_slog_Error:28","Error:multiple Read calls return no data or error"]}
}

func Test_slog_Debug(t *testing.T) {
	l := slog.New(NewCatcherHandler(&slog.HandlerOptions{Level: slog.LevelDebug}))

	l.Debug("Exsample!", "test_data", struct {
		name string
		id   int
	}{name: "User", id: 112})

	// 2025/01/08 14:05:01 {"Level":"DEBUG","Host":"nuc","Message":"Exsample!","Payload":["test_data:{User 112}"]}
}

func Test_slog_WithAttr_Info(t *testing.T) {
	l := slog.New(NewCatcherHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
	l = l.With("HOST", "localhost")
	l = l.With("ID", "112")
	l.Info("Exsample!", "req", "Request")
	l.Info("Exsample!", "user", "test")

	// 2025/01/08 18:17:16 {"Level":"INFO","Host":"nuc","Message":"Exsample!","Payload":["HOST:localhost","ID:112","req:Request"]}
	// 2025/01/08 18:17:16 {"Level":"INFO","Host":"nuc","Message":"Exsample!","Payload":["HOST:localhost","ID:112","user:test"]}
}

func Test_slog_WithAttr_WithGroup(t *testing.T) {
	l := slog.New(NewCatcherHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
	l = l.WithGroup("req")
	l = l.With("HOST", "localhost")
	l = l.With("ID", "112")
	l.Info("Exsample!", "par1", "data")
	l.Info("Exsample!", "par2", "test-data")

	// 2025/01/08 18:19:32 {"Level":"INFO","Host":"nuc","Message":"Exsample!","Payload":["req.HOST:localhost","req.ID:112","req.par1:data"]}
	// 2025/01/08 18:19:32 {"Level":"INFO","Host":"nuc","Message":"Exsample!","Payload":["req.HOST:localhost","req.ID:112","req.par2:test-data"]}
}

// INFO - 43604	     	31734 ns/op	    2114 B/op	      20 allocs/op - OLD
// INFO - 24566     	71679 ns/op     856 B/op          13 allocs/op

// INFO+WARN+ERROR - 32660 		33437 ns/op	    	3779 B/op	      64 allocs/op - OLD
// INFO+WARN+ERROR - 10000 		209826 ns/op		2667 B/op         46 allocs/op
// go test -timeout 30s -bench=^Benchmark_slog_handler$ github.com/ArseniySavin/catcher/pkg -benchmem -memprofile=mem.out
// go tool pprof mem.out
func Benchmark_slog_handler(b *testing.B) {
	l := slog.New(NewCatcherHandler(&slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))

	for i := 0; i < b.N; i++ {

		l.Info("Exsample Info!", "user", "TEST")
		l.Warn("Exsample Warn!", "data1", struct{ m string }{m: "Test1!"}, "data2", struct{ m string }{m: "Test2!"})
		l.Error("Exsample Error!", ErrorKey, io.ErrNoProgress)

	}
	b.ReportAllocs()
}
