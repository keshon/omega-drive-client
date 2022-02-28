package utils

import "syscall"

func SendCtrlBreak(pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		panic(e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		panic(e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		panic(e)
	}
}
