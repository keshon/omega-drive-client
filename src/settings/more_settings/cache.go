package more_settings

import (
	"app/src/state"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func cacheWindowScreen(w fyne.Window) fyne.CanvasObject {
	pathLabel := widget.NewLabel(state.SettingsValues.Cache.DefaultPath)
	if state.SettingsValues.Cache.OverridePath != "" {
		pathLabel.SetText(state.SettingsValues.Cache.OverridePath)
	}

	selectButton := widget.NewButton("Select folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}

			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			state.SettingsValues.Cache.OverridePath = strings.ReplaceAll(list.String(), "file://", "")
			pathLabel.SetText(strings.ReplaceAll(list.String(), "file://", ""))
		}, w)
	})

	resetButton := widget.NewButton("Reset to default", func() {
		cnf := dialog.NewConfirm("Confirmation", "Reset cache path to default?", func(b bool) {
			if b {
				pathLabel.SetText(state.SettingsValues.Cache.DefaultPath)
				state.SettingsValues.Cache.OverridePath = ""
			}
		}, w)
		cnf.SetConfirmText("Yes")
		cnf.SetDismissText("No")
		cnf.Show()

	})

	disabled := widget.NewCheck("Disable cache", func(b bool) {
		if b {
			state.SettingsValues.Cache.Disabled = true
		} else {
			state.SettingsValues.Cache.Disabled = false
		}
	})
	if state.SettingsValues.Cache.Disabled {
		disabled.SetChecked(true)
		selectButton.Disable()
		resetButton.Disable()
	} else {
		disabled.SetChecked(false)
		selectButton.Enable()
		resetButton.Enable()
	}

	buttons := container.NewHBox(selectButton, resetButton)

	sep := widget.NewSeparator()
	sep.Hide()

	inputs := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Cache location", Widget: buttons, HintText: ""},
			{Text: "", Widget: pathLabel, HintText: ""},
			{Text: "Disabled", Widget: disabled, HintText: "dont use cache for uploading"},
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

			if state.SettingsValues.Cache.Disabled {
				disabled.SetChecked(true)
				selectButton.Disable()
				resetButton.Disable()

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Cache location disabled",
					Content: "",
				})
			} else {
				disabled.SetChecked(false)
				selectButton.Enable()
				resetButton.Enable()

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Cache location is set to:",
					Content: pathLabel.Text,
				})
			}
			selectButton.Refresh()
			resetButton.Refresh()
		},
	}
	return container.NewBorder(container.NewVBox(inputs), container.NewBorder(nil, nil, nil, footer), nil, nil)
}
