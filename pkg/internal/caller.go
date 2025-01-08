package internal

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

type CallerInfo struct {
	Host        string
	PackageName string
	FileName    string
	FuncName    string
	Line        int
}

func GetHost() string {
	host := os.Getenv("HOST")

	if host == "" {
		host, err := os.Hostname()

		if err != nil {
			return fmt.Sprintf("Host undefined, %s", err.Error())
		}

		return host
	}

	return host
}

func CallInfo(caller int) *CallerInfo {
	host := GetHost()
	pc, file, line, _ := runtime.Caller(caller)

	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2] != "" {
		if parts[pl-2][0] == '(' {
			funcName = parts[pl-2] + "." + funcName
			packageName = strings.Join(parts[0:pl-2], ".")
		} else {
			packageName = strings.Join(parts[0:pl-1], ".")
		}
	} else {
		packageName = runtime.FuncForPC(pc).Name()
	}

	return &CallerInfo{
		PackageName: packageName,
		FileName:    fileName,
		FuncName:    funcName,
		Line:        line,
		Host:        host,
	}
}

func (c CallerInfo) String() string {
	return fmt.Sprintf("%s/%s %s:%d", c.PackageName, c.FileName, c.FuncName, c.Line)
}

func CallSource(PC uintptr) *CallerInfo {
	fs := runtime.CallersFrames([]uintptr{PC})
	f, _ := fs.Next()

	_, fileName := path.Split(f.File)

	fName := runtime.FuncForPC(PC).Name()
	index := strings.LastIndex(fName, ".")
	packageName := fName[0:index]
	funcName := fName[index:]

	return &CallerInfo{
		Host:        GetHost(),
		PackageName: packageName,
		FileName:    fileName,
		FuncName:    funcName,
		Line:        f.Line,
	}
}
