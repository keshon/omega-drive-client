package tray

import "golang.org/x/sys/windows"

func LoadIcon(path string) uintptr {
	icon, err := LoadImage(
		0,
		windows.StringToUTF16Ptr(path),
		IMAGE_ICON,
		0,
		0,
		LR_DEFAULTSIZE|LR_LOADFROMFILE)
	if err != nil {
		panic(err)
	}

	return icon
}
