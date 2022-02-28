package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println("[utils/open_browser.go][Openbrowser] Unsupported platform")
		log.Fatal(err)
	}

}
