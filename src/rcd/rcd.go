package rcd

import (
	"app/src/auth_webhook"
	"app/src/conf"
	"app/src/settings/auth_settings"
	"app/src/state"
	"app/src/utils"
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type RcdStruct struct {
}

var Rcd RcdStruct

func NewRCD() *RcdStruct {
	r := &RcdStruct{}
	r.init()
	return r
}

func (r RcdStruct) init() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Init RCD has started")

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Init RCD has finished")
}

func (r RcdStruct) StartRCD(onCron bool) bool {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Start RCD has started")

	// Get permissions from Airtable (via webhook)
	//accessRights := auth_webhook.PostToWebhook()
	auth_webhook.Authenticate()

	if len(state.Response.ID) == 0 {
		state.AccessKey = ""

		if conf.RcAutodeleteConfig {
			r.DeleteConfig()
		}

		return false
	}

	auth_settings.LoadAccessKey()

	if !onCron {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Start RCD on app launch")

		createList := state.Response.AvailablePaths

		// Mount list of drives
		if len(createList) > 0 {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Disks to mount:")
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(createList)

			for _, elem := range createList {
				log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(elem.Letter)
				r.Mount(elem.Name, elem.Letter, elem.RW) // or ReadWrite
			}

		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Drives/path to mount not found")
		}

	} else {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Start RCD on Cron event")

		// Dismount old drives, connect new
		oldPaths := state.ActivePaths
		newPaths := state.Response.AvailablePaths

		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Old permissions:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(oldPaths)

		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("New permissions:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(newPaths)

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

		removeList := utils.StringArrayDifference(oldList, validList) // list of disks to unmount
		createList := utils.StringArrayDifference(newList, validList) // list of disks to mount

		// Unmount list of drives
		if len(removeList) > 0 {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Drives/mount path to unmount:")
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(removeList)

			for _, elem := range state.ActivePaths {
				for _, toRemoveElem := range removeList {
					if elem.Name == toRemoveElem {
						r.Unmout(elem.Letter)
					}
				}
			}
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Drives/path to unmount not found")
		}

		// Mount list of drives
		if len(createList) > 0 {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Drives/mount path to mount:")
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(createList)

			for _, elem := range newPaths {
				for _, toCreateElem := range createList {
					if elem.Name == toCreateElem {
						r.Mount(elem.Name, elem.Letter, elem.RW)
					}
				}
			}
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Drives/path to mount not found")
		}

	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Start RCD has finished")

	return true
}

func (r RcdStruct) RefreshSyncingData() {
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every 1s", func() {
		raw := r.CoreTransferring().Each

		state.SyncingData = nil

		// Add to state var
		for i := 0; i < len(raw); i++ {
			if raw[i].Percentage > 0 && raw[i].Percentage <= 100 {
				el := state.SyncingDataStruct{}
				el.Label = raw[i].Name

				el.Progress = float64(raw[i].Percentage) / 100

				humanSpeed := utils.ByteCountDecimal(int64(raw[i].Speed))
				humanTotalSize := utils.ByteCountDecimal(int64(raw[i].Size))
				humanUploadedSize := utils.ByteCountDecimal(int64(raw[i].Bytes))

				el.Status = "Uploading at " + humanSpeed + "/sec - " + humanUploadedSize + " of " + humanTotalSize

				state.SyncingData = append(state.SyncingData, el)
			}
		}

		//fmt.Println(state.SyncingData)
	})
	go c.Run()
}

func (r RcdStruct) RefreshHistoryData() {
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every 1s", func() {
		raw := r.CoreTransfered().Each

		state.HistoryData = nil

		// Add to state var
		for i := 0; i < len(raw); i++ {

			el := state.HistoryDataStruct{}
			el.Label = raw[i].Name

			/*
				var size int64
				var err error
				if raw[i].Size != "" {
					size, err = strconv.ParseInt(raw[i].Size, 10, 64)
					if err != nil {
						panic(err)
					}
				}
				humanTotalSize := utils.ByteCountDecimal(size)
			*/

			el.Status = "Synced successfully"

			state.HistoryData = append([]state.HistoryDataStruct{el}, state.HistoryData...)
			//state.HistoryData = append(state.HistoryData, el)

		}

		//fmt.Println(state.HistoryData)
	})
	go c.Run()
}

// Start Rclone server, prepare config file
func (r RcdStruct) InitConfig() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Rclone config init has started")

	// Getting user config path
	confPath, err := os.UserConfigDir()
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("AppData/Roaming (UserConfigDir) dir was not found which is odd")
	}
	confPath = confPath + "\\rclone\\rclone.conf"

	// Delete old cofig
	if conf.RcAutodeleteConfig {
		r.DeleteConfig()
	}

	// Start Rclone Server as a separate process
	go func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Starting Rclone server")

		process := exec.Command("rcd.exe")
		process.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := process.Start()

		if err != nil {
			r.InitConfig()
		}
	}()

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		if conf.RcCreateConfig {

			apiConfigCreate(`name=ReadWrite&parameters={"provider":"` + conf.S3Name + `","type":"s3","access_key_id":"` + conf.S3ReadWriteAccessKey + `","secret_access_key":"` + conf.S3ReadWriteSecretKey + `","region":"` + conf.S3Region + `","acl":"public-read-write","endpoint":"` + conf.S3Endpoint + `","no_check_bucket":true}&type=s3&opt={"nonInteractive":true,"obscure":true}`)
			apiConfigCreate(`name=ReadOnly&parameters={"provider":"` + conf.S3Name + `","type":"s3","access_key_id":"` + conf.S3ReadOnlyAccessKey + `","secret_access_key":"` + conf.S3ReadOnlySecretKey + `","region":"` + conf.S3Region + `","acl":"public-read","endpoint":"` + conf.S3Endpoint + `","no_check_bucket":false}&type=s3&opt={"nonInteractive":true,"obscure":true}`)

			if conf.RcAutodeleteConfig {
				r.DeleteConfig()
			}
		}
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Rclone config init has finished")
}

func (r RcdStruct) Mount(mountPath, letter string, rw bool) {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Mounting to " + letter + " has started")

	// Mount paths
	if len(conf.BucketsMustContainName) != 0 {
		if len(mountPath) > 0 {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Checking if the bucket name is in the allowed list from conf")

			var match int

			for _, elem := range conf.BucketsMustContainName {
				res := strings.Contains(mountPath, elem)
				if res {
					match++
					log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Mount path " + mountPath + " is allowed. Proceed")
				}
			}
			if match == 0 {
				log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Mount path " + mountPath + " is NOT allowed. Skipping")
				return
			}
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Mount path is empty. Specify root path at least")
		}
	}

	// Letter
	// Use predefined drive letter or get any available letter
	if len(letter) == 0 {
		availLetters := utils.AvalableDriveLetters()

		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Available letters found: ")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(availLetters)

		if len(availLetters) > 0 {
			letter = availLetters[len(availLetters)-1]
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("No available letters for specified letter were found")
			r.Quit()
			return
		}
	}

	// Read-and-write or read-only
	var cacheMode string
	var name string
	var provider string
	if rw == true {
		provider = "ReadWrite"
		name = "full-access"
		cacheMode = "2"
		if state.SettingsValues.Cache.Disabled {
			cacheMode = "0"
		}
	} else {
		provider = "ReadOnly"
		name = "read-only"
		cacheMode = "0"
	}

	// API call
	apiMountMount(`mountPoint=` + letter + `:&` + `fs=` + provider + `:` + mountPath + `&mountType=cmount&vfsOpt={"CacheMode":` + cacheMode + `}&mountOpt={"AllowOther":true,"NetworkMode":true,"VolumeName":"\\\\` + name + `\\` + mountPath + `\\"}`)

	// Update available paths
	var current state.ActivePathsStruct
	current.Name = mountPath
	current.Letter = letter
	state.ActivePaths = append(state.ActivePaths, current)

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Mounting to " + letter + " has finished. Path should be active")
}

func (r RcdStruct) Unmout(letter string) {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Unmounting letter " + letter + " has started")

	apiMountUnmount(`mount/unmount?mountPoint=` + letter + `:`)

	for i := 0; i < len(state.ActivePaths); i++ {
		if state.ActivePaths[i].Letter == letter {
			state.ActivePaths = append(state.ActivePaths[:i], state.ActivePaths[i+1:]...)
		}
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Unmounting letter " + letter + " has finished")
}

func (r RcdStruct) UnmountAll() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Reconnect (unmount all) has started")

	state.AccessKey = ""
	apiMountUnmountall()

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Reconnect (unmount all) has finished")
}

func (r RcdStruct) Refresh() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Refreshing mounted paths has started")

	resp := apiVfsList()
	if len(resp.VFSES) == 0 {
		log.Println("[Refresh] Error getting vfs list via 'vfs/list'")
		log.Println(resp)
		return
	}

	for _, elem := range resp.VFSES {
		apiVfsRefresh(elem)
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Refreshing mounted paths has finished")
}

func (r RcdStruct) Quit() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Quiting Rclone rcd has started")

	state.AccessKey = ""

	apiMountUnmountall()
	apiFscacheClear()
	apiCoreQuit()
	r.DeleteCache()

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Quiting Rclone rcd has finished")
}

func (r RcdStruct) DeleteConfig() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Deleting Rclone config has started")

	path, err := os.UserConfigDir()
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("AppData/Roaming (UserConfigDir) dir was not found")
	}

	path = path + "\\rclone\\rclone.conf"
	err = os.Remove(path)
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Rclone config file was not found: " + err.Error())
	} else {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Rclone config file was found and deleted")
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Deleting Rclone config has finished")

	//r.DeleteCache()
}

func (r RcdStruct) DeleteCache() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Deleting Rclone cache has started")

	path, err := os.UserCacheDir()
	if err != nil {
		log.Println("[rclone/main.go][deleteCache] AppData/Local (UserCacheDir) dir was not found which is odd")
	}

	path = path + "\\rclone\\"

	err = os.RemoveAll(path)
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Rclone cache dir was not found: " + err.Error())
	} else {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Rclone cache dir was found and deleted")
	}

	if state.SettingsValues.Cache.DefaultPath != "" {
		err = os.RemoveAll(state.SettingsValues.Cache.DefaultPath)
		if err != nil {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Dir " + state.SettingsValues.Cache.DefaultPath + "  was not found: " + err.Error())
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Dir " + state.SettingsValues.Cache.DefaultPath + " was found and deleted")
		}
	}

	if state.SettingsValues.Cache.OverridePath != "" {
		err = os.RemoveAll(state.SettingsValues.Cache.OverridePath)
		if err != nil {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Dir " + state.SettingsValues.Cache.OverridePath + "  was not found: " + err.Error())
		} else {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Dir " + state.SettingsValues.Cache.OverridePath + " was found and deleted")
		}
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Deleting Rclone cache has finished")
}

func (r RcdStruct) CoreTransfered() CoreTransferedResponse {
	resp := apiCoreTransfered()
	/*
		fmt.Println(resp.Each)
		for _, elem := range resp.Each {
			fmt.Println(elem)
		}
	*/
	return resp
}

func (r RcdStruct) CoreTransferring() CoreStatsResponse {

	resp := apiCoreStats()
	/*
		//fmt.Println(resp.Each)
		for _, elem := range resp.Each {
			fmt.Println(elem)
		}
	*/
	return resp
}
