package utils

import (
	"fmt"
	"runtime"
)

/*
	Runtime callers methods allow to find parent func name and a file name it was called from
	The implementation is taken from here: https://stackoverflow.com/a/38551362/10352443

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-29
*/

// Output example: main.main
func CallerFuncLoc() string {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("called from %s\n", details.Name())
	}
	return ""
}

// Output example: /tmp/sandbox269058180/prog.go#16
func CallerFileLoc() string {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s#%d\n", file, no)
	}
	return ""
}
