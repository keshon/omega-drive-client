package more_settings

import (
	"app/src/state"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func generalWindowScreen(w fyne.Window) fyne.CanvasObject {
	var (
		programFolder = "Program folder"
		userFolder    = "User folder"
	)

	radio := widget.NewRadioGroup([]string{programFolder, userFolder}, func(s string) {
		fmt.Println("selected", s)
		if s == programFolder {
			state.SettingsValues.General.StoreLocationPortable = true
		} else {
			state.SettingsValues.General.StoreLocationPortable = false
		}

	})

	if state.SettingsValues.General.StoreLocationPortable {
		radio.SetSelected(programFolder)
	} else {
		radio.SetSelected(userFolder)
	}

	radio.Horizontal = false

	sep := widget.NewSeparator()
	sep.Hide()

	inputs := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Store settings at", Widget: radio, HintText: "where to store settings"},
			{Widget: sep},
		},
	}

	footer := &widget.Form{
		CancelText: "Close",
		OnCancel: func() {
			LoadSettingsValues()
			w.Close()
		},
		SubmitText: "Save",
		OnSubmit: func() {
			saveConfOnConfLocation()

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Success",
				Content: "Config file saved at " + radio.Selected,
			})
		},
	}
	return container.NewBorder(container.NewVBox(inputs), container.NewBorder(nil, nil, nil, footer), nil, nil)
}
