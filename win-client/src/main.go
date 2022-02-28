package main

// #cgo CFLAGS: -I c:/Temp/trayhost/platform/

import (
	"app/src/conf"
	"app/src/tray"
	"app/src/utils"

	"encoding/base64"
	"log"
	"os"
	"strconv"
)

var (
	theme string
)

func main() {

	// Create log
	err := os.Remove("logs.txt")
	if err != nil {
		log.Println(err)
	}
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// Store PID to Env variable
	log.Println("[main.go][main] App has started, PID is " + strconv.Itoa(os.Getpid()))
	os.Setenv("MD_PID", base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(os.Getpid()))))

	// App expiry timer
	log.Println("[main.go][main] Expiration set to " + conf.ValidTo)
	utils.ExpiryTimer(conf.ValidTo)

	tray.Main()
}
