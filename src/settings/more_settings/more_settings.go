package more_settings

import (
	"app/src/conf"
	"app/src/settings/auth_settings"
	"app/src/state"
	"app/src/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MoreSettingsStruct struct {
	Title string
	View  func(w fyne.Window) fyne.CanvasObject
}

var (
	preferenceCurrent = "cache"

	// Define the metadata for each tutorial
	MoreSettings = map[string]MoreSettingsStruct{
		"general": {"General", generalWindowScreen},
		"cache":   {"Cache", cacheWindowScreen},
		"remote":  {"Remote", remoteWindowScreen},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	MoreSettingsIndex = map[string][]string{
		"": {"general", "cache", "remote"},
	}
)

// NewSettings returns a new settings instance
func NewSettings() *MoreSettingsStruct {
	s := &MoreSettingsStruct{}

	return s
}

// LoadScreen creates a new settings screen to handle appearance configuration
func (s *MoreSettingsStruct) LoadScreen(w fyne.Window) fyne.CanvasObject {

	content := container.NewMax()

	title := widget.NewLabel("Component name")

	setMoreSettings := func(as MoreSettingsStruct) {

		title.SetText(as.Title)

		content.Objects = []fyne.CanvasObject{as.View(w)}
		content.Refresh()
	}

	settings := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)

	split := container.NewHSplit(makeNav(setMoreSettings, true), settings)
	split.Offset = 0.2

	w.SetOnClosed(func() {
		LoadSettingsValues()
	})

	return split
}

func makeNav(setSettings func(settings MoreSettingsStruct), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return MoreSettingsIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := MoreSettingsIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			as, ok := MoreSettings[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(as.Title)
		},
		OnSelected: func(uid string) {
			if as, ok := MoreSettings[uid]; ok {
				a.Preferences().SetString(preferenceCurrent, uid)
				setSettings(as)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrent, "cache")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func LoadSettingsValues() {
	// Set default values
	state.SettingsValues.General.StoreLocationPortable = true
	state.SettingsValues.General.AppdataPath = createAppdataDir()
	state.SettingsValues.Cache.DefaultPath = createAppdataDir() + "cache"
	state.SettingsValues.Cache.OverridePath = ""
	state.SettingsValues.Remote.ReconnectRate = Every1min

	// Find config file location and override the default settings
	appdataConf := state.SettingsValues.General.AppdataPath + "conf.json"

	localConfExist, err := utils.PathExist("conf.json")
	if err != nil {
		panic("err reading conf")
	}

	if localConfExist {
		confData := utils.ReadFile("conf.json")
		json.Unmarshal(confData, &state.SettingsValues)
	} else {
		appdataConfExists, err := utils.PathExist(appdataConf)
		if err != nil {
			panic("err reading appdata conf")
		}

		if appdataConfExists {
			confData := utils.ReadFile(appdataConf)
			json.Unmarshal(confData, &state.SettingsValues)
		} else {
			//createDefaultSettings()
			saveToFile(state.SettingsValues, "")
		}
	}

	// Move conf file according to it's setting
	saveConfOnConfLocation()
}

func saveConfOnConfLocation() {
	// work with key
	auth_settings.LoadAccessKey()
	if state.SettingsValues.General.StoreLocationPortable {
		e := os.Remove(state.SettingsValues.General.AppdataPath + "/conf.json")
		if e != nil {
			fmt.Println("file conf.json not found")
		}
		state.SettingsValues.General.StoreLocationPortable = true
		saveToFile(state.SettingsValues, "")
	} else {
		e := os.Remove("conf.json")
		if e != nil {
			fmt.Println("file conf.json not found")
		}
		state.SettingsValues.General.StoreLocationPortable = false
		saveToFile(state.SettingsValues, state.SettingsValues.General.AppdataPath)
	}
	auth_settings.SaveAccessKeyOnConfLocation()
}

func saveToFile(s state.SettingsValuesStruct, dest string) {
	file, _ := json.MarshalIndent(s, "", " ")

	err := ioutil.WriteFile(dest+"conf.json", file, 0644)
	if err != nil {
		panic("cant write conf")
	}
}

func createAppdataDir() string {
	appdataPath, err := os.UserConfigDir()
	if err != nil {
		panic("appdata path is invalid")
	}

	path := filepath.Join(appdataPath+"/", conf.AppName)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic("err creating new appdata path")
	}

	return path + "/"
}
