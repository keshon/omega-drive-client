package state

import (
	"app/src/conf"
	"encoding/base64"

	"fyne.io/fyne/v2/data/binding"
)

/*
	State package contains structs and vars that temporary store information during app run
*/

type Status string

const (
	Idle       Status = "Idle"
	Connecting Status = "Connecting"
	Connected  Status = "Connected"
	Syncing    Status = "Syncing"
	Error      Status = "Error"
)

type ConnectionStatusStruct struct {
	Status          Status
	BindTitle       binding.String
	BindDescription binding.String
}

// Files and folder are currently syncing
type SyncingDataStruct struct {
	Label    string
	Progress float64
	Status   string
}

// Files and folders already synced
type HistoryDataStruct struct {
	Label     string
	Status    string
	Timestamo string `json:"store_location_portable,omitempty"`
}

// Settings values
type SettingsValuesStruct struct {
	General struct {
		StoreLocationPortable bool   `json:"store_location_portable"`
		AppdataPath           string `json:"app_data_path"`
	}
	Cache struct {
		DefaultPath  string `json:"default_path"`
		OverridePath string `json:"override_path"`
		Disabled     bool   `json:"disabled"`
	}
	Remote struct {
		ReconnectRate string `json:"reconnect_rate"`
	}
}

// Drive path currently allocated
type ActivePathsStruct struct {
	Name   string `json:"Name,omitempty"`
	Letter string `json:"Letter,omitempty"`
	RW     bool   `json:"RW,omitempty"`
}

// Authenticated successfull response
type AvailPathsStruct struct {
	Name   string `json:"Name,omitempty"`
	Letter string `json:"Letter,omitempty"`
	RW     bool   `json:"RW,omitempty"`
}

type ResponseStruct struct {
	ID             string `json:"id,omitempty"`
	Fullname       string `json:"Fullname,omitempty"`
	AccessKey      string `json:"AccessKey,omitempty"`
	AvailablePaths []ActivePathsStruct
}

var (
	// App current connection status
	ConnectionStatus = ConnectionStatusStruct{BindTitle: binding.NewString(), BindDescription: binding.NewString()}

	AccessKey      string               // access key
	SettingsValues SettingsValuesStruct // settings
	ActivePaths    []ActivePathsStruct  // drive paths currently allocated
	SyncingData    []SyncingDataStruct  // rcd (rclone) response for currently syncing object(s)
	HistoryData    []HistoryDataStruct  // rcd (rclone) response for already synced object(s)
	Response       ResponseStruct       // webhook successfull response with available path list to mount

	// Encoded `username:password` pair for basic auths
	RcAuthEncoded  = base64.StdEncoding.EncodeToString([]byte(conf.RcUsername + ":" + conf.RcPassword))   // rcd
	N8nAuthEncoded = base64.StdEncoding.EncodeToString([]byte(conf.N8nUsername + ":" + conf.N8nPassword)) // webhook
)
