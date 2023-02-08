package auth_settings

import (
	"app/src/conf"
	"app/src/state"
	"app/src/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/widget"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}

	return s
}

// LoadAppearanceScreen creates a new settings screen to handle appearance configuration
func (s *Settings) LoadAccessKeyScreen(w fyne.Window) fyne.CanvasObject {

	key := widget.NewPasswordEntry()
	key.SetPlaceHolder("Enter access key")

	sep := widget.NewSeparator()
	sep.Hide()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Widget: sep},
			{Text: "Access key", Widget: key, HintText: "You can get your access key from a superior officer"},
			{Widget: sep},
		},
		CancelText: "Close",
		OnCancel: func() {
			key.SetText("")
			w.Close()
		},
		SubmitText: "Save",
		OnSubmit: func() {
			if len(key.Text) == 0 {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Failed",
					Content: "Access key is empty",
				})
				return
			}

			state.AccessKey = key.Text
			SaveAccessKeyOnConfLocation()
			key.SetText("")
			state.AccessKey = ""
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Success",
				Content: "Access key saved succesfully",
			})
		},
	}

	return form
}

func LoadAccessKey() {
	// Set default values
	state.AccessKey = ""

	// Find config file location and override the default settings
	appdataKey := state.SettingsValues.General.AppdataPath + "key"

	localFileExist, err := utils.PathExist("key")
	if err != nil {
		panic("err reading key")
	}

	if localFileExist {
		keyData := utils.ReadFile("key")
		if len(keyData) > 0 {
			decrypted, err := utils.Decrypt([]byte(conf.EncryptKey), keyData)
			if err != nil {
				panic(err)
			}
			state.AccessKey = string(decrypted)
		}
	} else {
		appdataKeyExists, err := utils.PathExist(appdataKey)
		if err != nil {
			panic("err reading appdata conf")
		}

		if appdataKeyExists {
			keyData := utils.ReadFile(appdataKey)
			if len(keyData) > 0 {
				decrypted, err := utils.Decrypt([]byte(conf.EncryptKey), keyData)
				if err != nil {
					panic(err)
				}
				state.AccessKey = string(decrypted)
			}
		} else {
			//createDefaultSettings()
			saveToFile("", "")
		}
	}

	// Move conf file according to it's setting
	SaveAccessKeyOnConfLocation()
}

func loadFromFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[utils/read_file.go][ReadFile] Error reading file. Details:")
		log.Println(err)
		return nil, err
	}
	return b, nil
}

func SaveAccessKeyOnConfLocation() {
	if state.SettingsValues.General.StoreLocationPortable {
		e := os.Remove(state.SettingsValues.General.AppdataPath + "/key")
		if e != nil {
			fmt.Println("file key not found")
		}
		saveToFile(state.AccessKey, "")
	} else {
		e := os.Remove("key")
		if e != nil {
			fmt.Println("file key not found")
		}
		saveToFile(state.AccessKey, state.SettingsValues.General.AppdataPath)
	}
}

func saveToFile(accessKey string, dest string) {
	// Simple validate
	accessKey = strings.TrimSuffix(string(accessKey), "\r\n")
	if len(accessKey) == 0 {
		log.Println("[settings/settings.go][WriteAccessKey] Access key is empty")
	}

	// Encrypting
	key := []byte(conf.EncryptKey) // 32 bytes
	plaintext := []byte(accessKey)

	ciphertext, err := utils.Encrypt(key, plaintext)
	if err != nil {
		log.Println("[settings/settings.go][WriteAccessKey] Error encrypting access key. Details:")
		log.Println(err)
	}

	// write the whole body at once
	err = ioutil.WriteFile(dest+"key", []byte(ciphertext), 0644)
	if err != nil {
		log.Println("[settings/settings.go][WriteAccessKey] Error write access key to file key. Details:")
		log.Println(err)
	}

}
