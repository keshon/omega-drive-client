package gui_nk

import (
	"fmt"

	"github.com/lxn/win"
)

func GetMonitorWidth() (int, int) {
	width := int(win.GetSystemMetrics(win.SM_CXSCREEN))
	height := int(win.GetSystemMetrics(win.SM_CYSCREEN))

	fmt.Printf("%dx%d\n", width, height)

	return width, height
}

func GetTasbarHeight() int {
	return 0
}

func OnError(code int32, msg string) {
	fmt.Printf("[glfw ERR]: error %d: %s", code, msg)
}

func B(v int32) bool {
	return v == 1
}

func BoolToInt(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
