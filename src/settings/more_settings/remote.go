package more_settings

import (
	"app/src/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	Every1min = "Every minute"
	Every3min = "Every 3 minutes"
	Every5min = "Every 5 minutes"
)

func remoteWindowScreen(w fyne.Window) fyne.CanvasObject {

	selectEntry := widget.NewSelectEntry([]string{Every1min, Every3min, Every5min})
	selectEntry.PlaceHolder = ""
	switch state.SettingsValues.Remote.ReconnectRate {
	case Every1min:
		selectEntry.SetText(Every1min)
	case Every3min:
		selectEntry.SetText(Every3min)
	case Every5min:
		selectEntry.SetText(Every5min)
	}
	/*
		randPort := widget.NewCheck("Random", func(bool) {})
		if state.SettingsValues.Remote.RandServerPort {
			randPort.SetChecked(true)
		}
	*/
	sep := widget.NewSeparator()
	sep.Hide()

	inputs := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Reconnection rate", Widget: selectEntry, HintText: "how often re-validate access key"},
			/*{Text: "Rclone server port", Widget: randPort, HintText: "default rclone port or random"},*/
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
				Content: "Remote settings saved succesfully",
			})
		},
	}
	return container.NewBorder(container.NewVBox(inputs), container.NewBorder(nil, nil, nil, footer), nil, nil)
}
