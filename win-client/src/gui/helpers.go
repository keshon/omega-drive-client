package gui

import (
	"app/src/external/nanogui"
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

func Notify(n *nanogui.Label, msg string) {
	n.SetCaption(msg)
}
