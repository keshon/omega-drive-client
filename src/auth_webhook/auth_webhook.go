package auth_webhook

import (
	"app/src/conf"
	"app/src/settings/auth_settings"
	"app/src/state"
	"app/src/utils"
	"fmt"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func Authenticate() {
	state.ConnectionStatus.Status = state.Connecting
	state.ConnectionStatus.BindTitle.Set("Connecting...")
	state.ConnectionStatus.BindDescription.Set("Connection to authentication server has started.")
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Connection to authentication server has started")

	// Prepare payload
	type Payload struct {
		AccessKey string
	}
	var payload Payload

	// Read access key
	auth_settings.LoadAccessKey()
	payload.AccessKey = state.AccessKey

	// Access key is missing in `key` file
	if len(payload.AccessKey) == 0 {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Warning! Access key is missing")
		state.ConnectionStatus.BindDescription.Set("Access key is missing. Check authorization settings")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Access key is missing. Check authorization settings")
		state.Response = state.ResponseStruct{}
		return
	}

	// Which URL to use. N8n offers two of them
	var path string
	if conf.IsDev {
		path = conf.N8nDevURL // Development webhook
	} else {
		path = conf.N8nURL // Production webhook
	}

	// Make request
	response := []state.ResponseStruct{}
	err := utils.Request(state.N8nAuthEncoded, http.MethodPost, path, payload, &response)
	if err != nil {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Error! Could not connect to authentication server")
		state.ConnectionStatus.BindDescription.Set(err.Error())
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Could not connect to authentication server: " + err.Error())
		return
	}

	// Access key not found or auth server is not reachable
	if len(response) <= 0 {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Warning! Access key not found / auth server is not reachable")
		state.ConnectionStatus.BindDescription.Set("The authentication server could no be reached or it did not find the specified key.\nContact administrator")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("The authentication server did not find the specified key. Contact administrator")
		state.Response = state.ResponseStruct{}
		return
	}

	state.ConnectionStatus.Status = state.Connecting
	state.ConnectionStatus.BindTitle.Set("Access key was found")
	state.ConnectionStatus.BindDescription.Set("Access key was found.\nChecking for available network drives/folders...")
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Access key was found")
	state.Response = response[0]
}

// Send access key to webhook
func PostToWebhook() state.ResponseStruct {
	state.ConnectionStatus.Status = state.Connecting
	state.ConnectionStatus.BindTitle.Set("Connecting...")
	state.ConnectionStatus.BindDescription.Set("Connection to authentication server has started.")
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Connection to authentication server has started")

	// Prepare payload
	type Payload struct {
		AccessKey string
	}
	var payload Payload

	// Read access key
	auth_settings.LoadAccessKey()
	payload.AccessKey = state.AccessKey

	// Access key is missing in `key` file
	if len(payload.AccessKey) == 0 {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Warning! Access key is missing")
		state.ConnectionStatus.BindDescription.Set("Access key is missing. Check authorization settings")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Access key is missing. Check authorization settings")
		return state.ResponseStruct{}
	}

	// Which URL to use. N8n offers two of them
	var path string
	if conf.IsDev {
		path = conf.N8nDevURL // Development webhook
	} else {
		path = conf.N8nURL // Production webhook
	}

	// Make request
	response := []state.ResponseStruct{}
	err := utils.Request(state.N8nAuthEncoded, http.MethodPost, path, payload, &response)
	if err != nil {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Error! Could not connect to authentication server")
		state.ConnectionStatus.BindDescription.Set(err.Error())
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Could not connect to authentication server: " + err.Error())
	}

	fmt.Println(response)

	// Access key not found
	if len(response) <= 0 {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Warning! Provided access key was not found")
		state.ConnectionStatus.BindDescription.Set("The authentication server did not find the specified key.\nContact administrator")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("The authentication server did not find the specified key. Contact administrator")
		return state.ResponseStruct{}
	}

	// TODO: not working yet - need to fix n8n workflow to return response with empty path list
	// Drive list is empty
	if len(response[0].AvailablePaths) <= 0 {
		state.ConnectionStatus.Status = state.Error
		state.ConnectionStatus.BindTitle.Set("Warning! No drives are assinged")
		state.ConnectionStatus.BindDescription.Set("Access key is valid but no drives/paths are assigned to it.\nContact administrator")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Access key is valid but no drives/paths are assigned to it. Contact administrator")
		return state.ResponseStruct{}
	}

	state.ConnectionStatus.Status = state.Connecting
	state.ConnectionStatus.BindTitle.Set("Access key was found")
	state.ConnectionStatus.BindDescription.Set("Access key was found.\nChecking for available network drives/folders...")
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Access key was found")
	return response[0]
}
