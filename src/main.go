// Package main provides various examples of Fyne API capabilities.
package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"app/src/conf"
	"app/src/main_window"
	"app/src/rcd"
	"app/src/settings/auth_settings"
	"app/src/settings/more_settings"
	"app/src/settings/theme_settings"
	"app/src/state"
	"app/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var topWindow fyne.Window

var connectItem *fyne.MenuItem
var disconnectItem *fyne.MenuItem

func main() {
	os.Setenv("MD_PID", base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(os.Getpid()))))

	initLog()
	rcd.Rcd.InitConfig()
	rcd.Rcd.RefreshSyncingData()
	rcd.Rcd.RefreshHistoryData()

	more_settings.LoadSettingsValues()
	auth_settings.LoadAccessKey()

	a := app.NewWithID(conf.AppName)
	a.SetIcon(theme.FyneLogo())

	logLifecycle(a)

	w := a.NewWindow(conf.AppName)
	w.CenterOnScreen()
	w.SetCloseIntercept(func() {
		//w.Hide()
		w.Close()
	})
	w.SetOnClosed(func() {
		rcd.Rcd.Quit()
	})

	topWindow = w

	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	app := container.NewBorder(
		nil, nil, nil, nil, main_window.MakeMainWindow(w))

	w.SetContent(app)
	w.Resize(fyne.NewSize(460, 740))
	/*
		go func() {
			systray.Run(func() {
				systray.SetIcon(icon.Data)
				systray.SetTitle(conf.AppName)
				systray.SetTooltip(conf.AppDesc)

				mToggleMainWindow := systray.AddMenuItem("Show UI", "Show main window")
				// Sets the icon of a menu item. Only available on Mac and Windows.
				mToggleMainWindow.SetIcon(icon.Data)
				go func() {
					<-mToggleMainWindow.ClickedCh
					w.Show()
					w.RequestFocus()
				}()

				mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
				// Sets the icon of a menu item. Only available on Mac and Windows.
				mQuit.SetIcon(icon.Data)
				go func() {
					<-mQuit.ClickedCh
					systray.Quit()
					w.Close()
					a.Quit()
				}()
			}, func() {

			})
		}()
	*/
	w.ShowAndRun()
}

func initLog() {
	_ = os.Remove(conf.LogFilename)

	f, err := os.OpenFile(conf.LogFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		//log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		//log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		//log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		//log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Lifecycle: Exited Foreground")
	})
}

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	connectItem = fyne.NewMenuItem("Connect", func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Connect context menu clicked")

		go func() {
			keyStatus := rcd.Rcd.StartRCD(false)

			if keyStatus {
				connectItem.Disabled = true
				disconnectItem.Disabled = false

				state.ConnectionStatus.Status = state.Connected
				state.ConnectionStatus.BindTitle.Set("Drives mounted successfuly")

				var list string
				for _, elem := range state.ActivePaths {
					rw := "read-only"
					if elem.RW {
						rw = "full-access"
					}
					list += elem.Name + "(" + rw + ") (" + elem.Letter + ":)\n"
				}
				state.ConnectionStatus.BindDescription.Set("Mounted drives should be accessible now:\n" + list)

				rcd.Rcd.CoreTransferring()

			} else {
				connectItem.Disabled = false
				disconnectItem.Disabled = true

				go func() {
					state.ConnectionStatus.Status = state.Idle
					rcd.Rcd.UnmountAll()
				}()
			}
		}()

	})

	disconnectItem = fyne.NewMenuItem("Disconnect", func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Disconnect context menu clicked")

		go func() {
			state.ConnectionStatus.BindTitle.Set("Unmounting drives...")
			state.ConnectionStatus.BindDescription.Set("Unmounting drives is in process. Please wait")

			rcd.Rcd.UnmountAll()
			state.ActivePaths = nil

			state.ConnectionStatus.BindTitle.Set("Drives unmounted successfully")
			state.ConnectionStatus.BindDescription.Set("Drives unmounted successfully")

			connectItem.Disabled = false
			disconnectItem.Disabled = true
		}()
	})
	disconnectItem.Disabled = true

	themeSettingsItem := fyne.NewMenuItem("Theme", func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Theme context menu clicked")

		w := a.NewWindow("Theme settings")
		w.SetContent(theme_settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(460, 460))
		w.CenterOnScreen()
		w.Show()
	})

	autorizationSettingsItem := fyne.NewMenuItem("Authorization", func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Authorization settings context menu clicked")

		w := a.NewWindow("Authorization settings")
		w.SetContent(auth_settings.NewSettings().LoadAccessKeyScreen(w))
		w.Resize(fyne.NewSize(460, 80))
		w.CenterOnScreen()
		w.Show()
	})

	advancedSettingsItem := fyne.NewMenuItem("More...", func() {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Advanced settings context menu clicked")

		w := a.NewWindow("Advanced settings")
		w.SetContent(more_settings.NewSettings().LoadScreen(w))
		w.Resize(fyne.NewSize(660, 450))
		w.CenterOnScreen()
		w.Show()
	})

	/*
		tutorialsItem := fyne.NewMenuItem("Tutorials", func() {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Tutorials context menu clicked")

			w := a.NewWindow("Tutorials")
			w.SetContent(tutorials.NewTutorials().LoadScreen(w))
			w.Resize(fyne.NewSize(800, 600))
			w.CenterOnScreen()
			w.Show()
		})
	*/

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Documentation context menu clicked")

			u, _ := url.Parse(conf.AppHomepageURL)
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("GitHub", func() {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("GitHub context menu clicked")

			u, _ := url.Parse(conf.AuthorURL)
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("About", func() {
			log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("About context menu clicked")

			content := conf.AppDesc + "\n" + conf.AppRcloneVersion + "\n" + "\n" + conf.AppAuthor + "\n" + conf.AuthorEmail + "\n"

			dialog.ShowInformation("About "+conf.AppName, content, w)
		}),
	)

	// Connection loop
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	interval := conf.WebhookInterval
	switch state.SettingsValues.Remote.ReconnectRate {
	case more_settings.Every1min:
		interval = "1m"
	case more_settings.Every3min:
		interval = "3m"
	case more_settings.Every5min:
		interval = "5m"
	}

	c.AddFunc("@every "+interval, func() {
		status, _ := state.ConnectionStatus.BindTitle.Get()
		if status == string(state.Connected) {

			auth_settings.LoadAccessKey()
			keyStatus := rcd.Rcd.StartRCD(true)

			if keyStatus {
				connectItem.Disabled = true
				disconnectItem.Disabled = false

				state.ConnectionStatus.Status = state.Connected
				state.ConnectionStatus.BindTitle.Set("Drives mounted successfuly")
				state.ConnectionStatus.BindDescription.Set("Mounted drives should be accessible now")

				rcd.Rcd.CoreTransferring()

			} else {
				connectItem.Disabled = false
				disconnectItem.Disabled = true

				state.ConnectionStatus.Status = state.Error
				state.ConnectionStatus.BindTitle.Set("Error! Mounting has failed")
				state.ConnectionStatus.BindDescription.Set("Error! Mounting has failed")

				go func() {
					state.ConnectionStatus.Status = state.Idle
					rcd.Rcd.UnmountAll()
				}()
			}
		}
	})
	go c.Run()

	return fyne.NewMainMenu(
		fyne.NewMenu("Connect", connectItem, disconnectItem),
		fyne.NewMenu("Settings", autorizationSettingsItem, themeSettingsItem, fyne.NewMenuItemSeparator(), advancedSettingsItem), //, tutorialsItem),
		helpMenu,
	)
}
