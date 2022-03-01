package main

import (
	"app/src/conf"
	"app/src/structs"
	"app/src/utils"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/s3"
	"github.com/rclone/rclone/cmd"
	_ "github.com/rclone/rclone/cmd/all"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/process"
)

func main() {
	// Create log
	err := os.Remove("logs-server.txt")
	if err != nil {
		log.Println(err)
	}
	file, err := os.OpenFile("logs-server.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	log.Println("[main.go][main] App has started, PID is " + strconv.Itoa(os.Getpid()))

	// Setup CRON on schedule
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every "+conf.CheckParentInterval, func() {
		if !conf.SkipPIDValidation {
			resp := verifyPID()
			if !resp {
				Quit()
			}
		}
	})
	go c.Run()

	if !conf.SkipPIDValidation {
		resp := verifyPID()
		if !resp {
			os.Exit(1)
		}
	}

	// Exec rclone with arguments
	os.Args = []string{"rclone", "rcd", "--rc-user", conf.RcUsername, "--rc-pass", conf.RcPassword, "--rc-addr", strings.Replace(strings.Replace(conf.RcHost, "http://", "", -1), "/", "", -1)}
	cmd.Main()

}

// Verify PID
func verifyPID() bool {
	validResp := true

	// Verify parent PID with reference PID taken fron env var
	rpid, err := base64.StdEncoding.DecodeString(os.Getenv("MD_PID"))
	if err != nil {
		log.Println("[main.go][verifyPID] Can't decode reference PID. Exiting..")
		validResp = false
	}

	ppid := os.Getppid()

	isValid, err := process.PidExists(int32(ppid))

	if err != nil {
		log.Println("[main.go][verifyPID] Parent PID does not exist. Details:")
		log.Println(err)
		log.Println("[main.go][verifyPID] Exiting via os.Exit..")
		validResp = false
	}

	if !isValid {
		log.Println("[main.go][verifyPID] Parent PID is not valid. Exiting..")
		validResp = false
	}

	if string(rpid) != strconv.Itoa(ppid) {
		log.Println("[main.go][verifyPID] Parent PID does not match with reference. Exiting..")
		validResp = false
	}

	return validResp
}

// Quit Rclone Server
func Quit() {
	log.Println("[Quit] Quting Rclone has started")

	var resp structs.RcloneResponse
	path := "mount/unmountall"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.Println("[main.go][Quit] Error unmounting all letters via 'mount/unmountall'. Details:")
		log.Println(resp)
	}

	path = "fscache/clear"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.Println("[main.go][Quit] Error clearing cache via 'fscache/clear'. Details:")
		log.Println(resp)
	}

	path = "core/quit"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.Println("[main.go][Quit] Error quitting Rclone via 'core/quit'. Details:")
		log.Println(resp)
	}

	log.Println("[main.go][Quit] Quting Rclone has finished")
}
