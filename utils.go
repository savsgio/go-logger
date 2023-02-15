package logger

import (
	"runtime"
)

func getFileCaller(calldepth int) (frame runtime.Frame) {
	pc := make([]uintptr, 1)

	numFrames := runtime.Callers(calldepth, pc)
	if numFrames < 1 {
		frame.File = "???"
		frame.Line = 0
	} else {
		frame, _ = runtime.CallersFrames(pc).Next()
	}

	return frame
}
