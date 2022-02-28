package rclone

import (
	"app/src/conf"
	"app/src/states"
	"app/src/utils"
	"app/src/webhook"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func deleteCache() {
	log.Println("[rclone/main.go][deleteCache] Deleting Rclone cache has started")
	// Check for exiting config file and remove it
	path, err := os.UserCacheDir()
	if err != nil {
		log.Println("[rclone/main.go][deleteCache] AppData/Local (UserCacheDir) dir was not found which is odd")
	}

	path = path + "\\rclone\\"
	err = os.RemoveAll(path)
	if err != nil {
		log.Println("[rclone/main.go][deleteCache] Rclone cache dir was not found. Details:")
		log.Println(err)
	} else {
		log.Println("[rclone/main.go][deleteCache] Rclone cache dir was found and deleted")
	}
	log.Println("[rclone/main.go][deleteCache] Deleting Rclone cache has finished")
}

func Main(onCron bool) bool {
	log.Println("[rclone/main.go][Main] has called")

	// Get permissions from Airtable (via webhook)
	accessRights := webhook.PostToWebhook()

	if len(accessRights.ID) == 0 {
		log.Println("[rclone/main.go][Main] Permissions not found. Exiting..")

		states.AcceskeyIsValid = false

		// Delete Rclone config
		if conf.RcAutodeleteConfig {
			DeleteConfig()
		}

		return false
	}

	log.Println("[rclone/main.go][Main] Permissions found")
	log.Println("[rclone/main.go][Main] Permissions found - update drives according to create/remove lists")

	states.AcceskeyIsValid = true

	if !onCron {
		log.Println("[Main] Mount drives on first run")

		createList := accessRights.AvailPaths

		// Create Action
		if len(createList) > 0 {
			log.Println("[rclone/main.go][Main] Disks to mount:")
			log.Println(createList)

			for _, elem := range createList {
				log.Println(elem.Letter)
				Mount(elem.Name, elem.Letter, elem.RW) // or ReadWrite
			}

		} else {
			log.Println("[rclone/main.go][Main] Nothing to mount")
		}

	} else {
		log.Println("[rclone/main.go][Main] Mount drives on Cron")

		// Dismount old drives, connect new
		oldPaths := webhook.ActivePaths
		newPaths := accessRights.AvailPaths

		log.Println("[rclone/main.go][Main] Old permissions:")
		log.Println(oldPaths)

		log.Println("[rclone/main.go][Main] New permissions:")
		log.Println(newPaths)

		//log.Println(accessRights)

		// Create remove and create drive(s) list
		var validList []string
		var oldList []string
		var newList []string

		for _, oldElem := range oldPaths {
			oldList = append(oldList, oldElem.Name)
			for _, newElem := range newPaths {
				newList = append(newList, newElem.Name)
				if oldElem.Name == newElem.Name {
					validList = append(validList, oldElem.Name)
				}
			}
		}

		removeList := utils.Difference(oldList, validList) // list of disks to unmount
		createList := utils.Difference(newList, validList) // list of disks to mount

		// Remove action
		if len(removeList) > 0 {

			log.Println("[rclone/main.go][Main] Disks to unmount:")
			log.Println(removeList)

			for _, elem := range webhook.ActivePaths {
				for _, toRemoveElem := range removeList {
					if elem.Name == toRemoveElem {
						Unmout(elem.Letter)
					}
				}
			}

		} else {
			log.Println("[rclone/main.go][Main] Nothing to unmount")
		}

		// Create Action
		if len(createList) > 0 {

			log.Println("[rclone/main.go][Main] Disks to mount:")
			log.Println(createList)

			for _, elem := range newPaths {
				for _, toCreateElem := range createList {
					if elem.Name == toCreateElem {
						Mount(elem.Name, elem.Letter, elem.RW)
					}
				}
			}

		} else {
			log.Println("[rclone/main.go][Main] Nothing to mount")
		}

	}

	return true
}

// Start Rclone server, prepare config file
func InitConfig() {
	log.Println("[rclone/main.go][Init] Rclone config init has started")

	// Getting user config path
	log.Println("[rclone/main.go][Init] Checking for AppData/Roaming (UserConfigDir) dir")
	confPath, err := os.UserConfigDir()
	if err != nil {
		log.Println("[rclone/main.go][Init] AppData/Roaming (UserConfigDir) dir was not found which is odd")
	}
	confPath = confPath + "\\rclone\\rclone.conf"

	// Delete old cofig
	if conf.RcAutodeleteConfig {
		DeleteConfig()
	}

	// Start Rclone Server as a separate process
	go func() {
		log.Println("[rclone/main.go][Init] Starting Rclone server")
		process := exec.Command("rcd.exe")
		process.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := process.Start()
		if err != nil {
			//log.Println("[Init] Error starting Rclone server. Details:")
			//log.Print(err)
			// TODO: Should terminate whole app
			InitConfig()
		}
	}()

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		if conf.RcCreateConfig {

			// Create read and wirte (full access)
			log.Println("[rclone/main.go][Init] Creating config for ReadWrite")
			APIConfigCreate(`name=ReadWrite&parameters={"provider":"Wasabi","type":"s3","access_key_id":"` + conf.S3ReadWriteAccessKey + `","secret_access_key":"` + conf.S3ReadWriteSecretKey + `","region":"` + conf.S3Region + `","acl":"public-read-write","endpoint":"` + conf.S3Endpoint + `","no_check_bucket":true}&type=s3&opt={"nonInteractive":true,"obscure":true}`)

			// Create read only
			log.Println("[rclone/main.go][Init] Creating config for ReadOnly")
			APIConfigCreate(`name=ReadOnly&parameters={"provider":"Wasabi","type":"s3","access_key_id":"` + conf.S3ReadOnlyAccessKey + `","secret_access_key":"` + conf.S3ReadOnlySecretKey + `","region":"` + conf.S3Region + `","acl":"public-read","endpoint":"` + conf.S3Endpoint + `","no_check_bucket":false}&type=s3&opt={"nonInteractive":true,"obscure":true}`)

			if conf.RcAutodeleteConfig {
				DeleteConfig()
			}
		}
	}

	log.Println("[rclone/main.go][Init] Rclone config init has finished")
}

// Mount mountPath to letter
func Mount(mountPath, letter string, rw bool) {

	log.Println("[rclone/main.go][Mount] Mounting to " + letter + " has started")

	// Mount paths
	if len(conf.AllowedKeywords) != 0 {
		if len(mountPath) > 0 {
			// Validate list of accessible buckets
			log.Println("[rclone/main.go][Mount] Checking if bucket root name is in allowed list from conf:")
			var match int

			for _, elem := range conf.AllowedKeywords {
				res := strings.Contains(mountPath, elem)
				log.Println(res) // true
				if res {
					match++
					log.Println("[rclone/main.go][Mount] Mount path " + mountPath + " is allowed")
				}
			}
			if match == 0 {
				log.Println("[rclone/main.go][Mount] Mount path " + mountPath + " is NOT allowed. Skipping..")
				return
			}
		} else {
			log.Println("[rclone/main.go][Mount] Error! Mount path is empty. Specify root path at least.")
		}
	}

	// Letter
	// Deal with predefined drive letter or create new automatically
	if len(letter) == 0 {
		availLetters := utils.GetAvailDriveLetters()

		log.Println("[rclone/main.go][Mount] Posting available letters:")
		log.Println(availLetters)

		if len(availLetters) > 0 {
			letter = availLetters[len(availLetters)-1]
		} else {
			log.Println("[rclone/main.go][Mount] Error! No available letters not specified letter were found. Quitting")
			Quit()
			return
		}
	}

	// Read and write or read only
	var cacheMode string
	var name string
	var provider string
	if rw == true {
		provider = "ReadWrite"
		name = "full-access"
		cacheMode = "2"
	} else {
		provider = "ReadOnly"
		name = "read-only"
		cacheMode = "0"
	}

	// Api call
	APIMountMount(`mountPoint=` + letter + `:&` + `fs=` + provider + `:` + mountPath + `&mountType=cmount&vfsOpt={"CacheMode":` + cacheMode + `}&mountOpt={"AllowOther":true,"NetworkMode":true,"VolumeName":"\\\\` + name + `\\` + mountPath + `\\"}`)

	// Update available path
	var current webhook.AvailPaths
	current.Name = mountPath
	current.Letter = letter
	webhook.ActivePaths = append(webhook.ActivePaths, current)

	log.Println("[rclone/main.go][Mount] Mounting to " + letter + " has finished. Path is active")
}

// Unmount mountPath by letter
func Unmout(letter string) {
	log.Println("[rclone/main.go][Unmout] Unmounting letter " + letter + " has started")

	APIMountUnmount(`mount/unmount?mountPoint=` + letter + `:`)

	for i := 0; i < len(webhook.ActivePaths); i++ {
		if webhook.ActivePaths[i].Letter == letter {
			webhook.ActivePaths = append(webhook.ActivePaths[:i], webhook.ActivePaths[i+1:]...)
		}
	}

	log.Println("[rclone/main.go][Unmout] Unmounting letter " + letter + " has finished")
}

// Unmount all mounts
func UnmountAll() {
	log.Println("[rclone/main.go][Reconnect] Reconnect (unmount all) has started")

	states.AcceskeyIsValid = false

	APIMountUnmountall()

	log.Println("[rclone/main.go][Reconnect] Reconnect (unmount all) has finished")
}

// Refresh Dirs
func Refresh() {
	log.Println("[rclone/main.go][Refresh] Refreshing mounted paths has started")

	resp := APIVfsList()
	if len(resp.VFSES) == 0 {
		log.Println("[Refresh] Error getting vfs list via 'vfs/list'")
		log.Println(resp)
		return
	}

	for _, elem := range resp.VFSES {
		APIVfsRefresh(elem)
	}

	log.Println("[rclone/main.go][Refresh] Refreshing mounted paths has finished")
}

// Quit Rclone Server
func Quit() {
	log.Println("[rclone/main.go][Quit] Quting Rclone has started")

	states.AcceskeyIsValid = false

	APIMountUnmountall()

	APIFscacheClear()

	APICoreQuit()

	log.Println("[rclone/main.go][Quit] Quting Rclone has finished")
}

func DeleteConfig() {
	log.Println("[rclone/main.go][DeleteConfig] Deleting Rclone config has started")
	log.Println("[rclone/main.go][DeleteConfig] Checking for AppData/Roaming (UserConfigDir) dir")
	path, err := os.UserConfigDir()
	if err != nil {
		log.Println("[rclone/main.go][DeleteConfig] AppData/Roaming (UserConfigDir) dir was not found which is odd")
	}

	path = path + "\\rclone\\rclone.conf"
	err = os.Remove(path)
	if err != nil {
		log.Println("[rclone/main.go][DeleteConfig] Rclone config file was not found. Details:")
		log.Println(err)
	} else {
		log.Println("[rclone/main.go][DeleteConfig] Rclone config file was found and deleted")
	}

	log.Println("[rclone/main.go][DeleteConfig] Deleting Rclone config has finished")

	deleteCache()
}

func CoreTransfered() {
	resp := APICoreTransfered()
	fmt.Println(resp.Each)
	for _, elem := range resp.Each {
		fmt.Println(elem)
	}
}

func CoreTransferring() CoreStatsResponse {

	resp := APICoreStats()
	fmt.Println(resp.Each)
	for _, elem := range resp.Each {
		fmt.Println(elem)
	}
	return resp
}

/*
func JobsStatus() {
	// Setup CRON on Schedule
	log.Println("[gfxMain] Create Cron")
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	log.Println("[JobsStatus] Running cron")
	list := APIJobList()
	if len(list.Each) > 0 {
		for _, elem := range list.Each {
			resp := APIJobStatus(elem)
			fmt.Println(resp)
		}
	}

	c.AddFunc("@every 3s", func() {

	})
	log.Println("[gfxMain] Run Cron")
	c.Run()
}
*/
