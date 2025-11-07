package internal

import (
	"os"
	"path"
	"runtime"
	"strconv"
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
			return "undefined"
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
	return (c.PackageName + "/" + c.FileName + " " + c.FuncName + ":" + strconv.Itoa(c.Line))
}

func CallSource(PC uintptr, host string) CallerInfo {
	fs := runtime.CallersFrames([]uintptr{PC})
	f, _ := fs.Next()

	fFile := strings.LastIndex(f.File, "/")
	fName := runtime.FuncForPC(PC).Name()
	index := strings.LastIndex(fName, ".")

	return CallerInfo{
		Host:        host,
		PackageName: fName[0:index],
		FileName:    f.File[fFile+1:],
		FuncName:    fName[index:],
		Line:        f.Line,
	}
}
