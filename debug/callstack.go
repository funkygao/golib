package debug

import (
	"runtime"
)

type callStack struct {
	Func   string
	File   string
	LineNo int
}

func Callstack(skipFrames int) callStack {
	pc, file, lineno, _ := runtime.Caller(skipFrames)
	f := runtime.FuncForPC(pc)
	return callStack{f.Name(), file, lineno}
}
