package gui_ng

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"app/src/conf"
	"app/src/external/nanogui"
	"app/src/external/nanovgo"
	"app/src/rclone"
	"app/src/settings"
	"app/src/states"
	"app/src/utils"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/robfig/cron/v3"
)

var (
	mainloopActive bool = false

	// Mouse cursor position
	MousePosX int
	MousePosY int

	// Main window state
	ShowMainWindow     bool = false
	MainWindowIsHidden bool = true

	// NG windows instances
	MainNGWindow           *nanogui.Window
	EmptyTransfersNGWindow *nanogui.Window
	TransfersNGWindow      *nanogui.Window
	SettingsNKWindow       *nanogui.Window
	AboutNKWindow          *nanogui.Window

	info      string
	infoLabel *nanogui.Label
	c         *cron.Cron
)

const (
	// Main window (incl. NG windows) fixed size
	winWidth  = 400
	winHeight = 500

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024

	fontSize   = 19
	elemHeight = 31
)

var app Application

type Application struct {
	screen *nanogui.Screen
}

type Transfer struct {
	Widget      *nanogui.Widget
	Label       *nanogui.Label
	ProgressBar *nanogui.ProgressBar
}

func (a *Application) init() {
	go func() {
		http.ListenAndServe(":5555", http.DefaultServeMux)
	}()

	glfw.WindowHint(glfw.Samples, 4)

	a.screen = nanogui.NewScreen(winWidth, winHeight, conf.AppName+" ("+conf.AppVersion+")", true, false)

	// Call NG windows here
	gfxMain(a.screen)
	gfxSettings(a.screen)
	gfxAbout(a.screen)
	//MiscWidgetsDemo(a.screen)

	a.screen.PerformLayout()
	a.screen.DebugPrint()

	/* All NanoGUI widgets are initialized at this point. Now
	create an OpenGL shader to draw the main window contents.
	*/
}

func initGlfw() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	/*
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	*/
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)

	nanogui.StartTime = time.Now()
}

func mainLoop() {
	mainloopActive = true

	var wg sync.WaitGroup

	/* If there are no mouse/keyboard events, try to refresh the
	view roughly every 50 ms; this is to support animations
	such as progress bars while keeping the system load
	reasonably low */
	wg.Add(1)
	go func() {
		for mainloopActive {
			time.Sleep(50 * time.Millisecond)
			glfw.PostEmptyEvent()
		}
		wg.Done()
	}()
	for mainloopActive {
		// Nanogui windows
		haveActiveScreen := false
		for _, screen := range nanogui.NanoguiScreens {
			if !screen.Visible() {
				continue
			} else if screen.GLFWWindow().ShouldClose() {
				screen.SetVisible(false)
				continue
			}
			//screen.DebugPrint()
			screen.DrawAll()
			haveActiveScreen = true
		}
		if !haveActiveScreen {
			mainloopActive = false
			break
		}

		// Main window control
		if !ShowMainWindow {
			if !MainWindowIsHidden {
				MainWindowIsHidden = true
				app.screen.Window.Hide()
			}
		} else {
			if MainWindowIsHidden {
				MainWindowIsHidden = false

				// Set window position
				scrWidth, _ := GetMonitorWidth()
				winPosX := scrWidth - (winWidth / 2) - (scrWidth - MousePosX)
				winPosY := MousePosY - winHeight - 28

				fmt.Println(MousePosX)
				fmt.Println(MousePosY)

				app.screen.Window.SetPos(winPosX, winPosY)
				app.screen.Window.Show()
			}
		}

		glfw.WaitEvents()
	}

	wg.Wait()
}

func Main() {
	initGlfw()

	if conf.IsDev {
		nanogui.SetDebug(true)
	}

	app.init()
	app.screen.DrawAll()
	app.screen.SetVisible(true)
	mainLoop()
}

func gfxMain(screen *nanogui.Screen) {
	splitterHeight := 120
	if conf.IsDev {
		splitterHeight = 180
	}

	MainNGWindow = nanogui.NewWindow(screen, conf.AppName)

	MainNGWindow.SetPosition(0, 0)
	MainNGWindow.SetFixedWidth(winWidth)
	MainNGWindow.SetFixedHeight(splitterHeight)
	MainNGWindow.SetLayout(nanogui.NewGroupLayout(0, 5, 0))
	MainNGWindow.SetVisible(true)

	menuWidget := nanogui.NewWidget(MainNGWindow)
	//menu.SetFixedWidth(winWidth)
	//menu.SetWidth(winWidth)
	//menu.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Maximum, 5, 10))
	menuWidget.SetLayout(nanogui.NewGridLayout(nanogui.Horizontal, 6, nanogui.Fill, 10, 5))

	actionButton := nanogui.NewPopupButton(menuWidget, "Diconnected")
	actionButton.SetIcon(nanogui.IconFlash)
	popup := actionButton.Popup()
	popup.SetLayout(nanogui.NewGroupLayout())

	// Action buttons
	//nanogui.NewLabel(popup, "Select action below:")
	//nanogui.NewCheckBox(popup, "A check box")
	connectButton := nanogui.NewButton(popup, "Connect")
	refreshButton := nanogui.NewButton(popup, "Refresh")
	reconnectButton := nanogui.NewButton(popup, "Reconnect")
	disconnectButton := nanogui.NewButton(popup, "Disconnect")

	connectButton.SetEnabled(true)
	refreshButton.SetEnabled(false)
	reconnectButton.SetEnabled(false)
	disconnectButton.SetEnabled(false)

	connectButton.SetCallback(func() {
		actionButton.SetEnabled(false)
		popup.SetVisible(false)

		// Get permissions
		log.Println("[gui_ng/main.go][gfxMain] connectButton")

		actionButton.SetCaption("Connecting..  ")
		info = "Validating access key.."
		Notify(infoLabel, info)
		states.Status = states.StatusWait

		go func() {
			keyStatus := rclone.Main(false)

			if keyStatus {
				actionButton.SetEnabled(true)
				actionButton.SetCaption("Connected")
				info = "Successfully connected"
				Notify(infoLabel, info)

				states.Status = states.StatusSuccess
				states.AcceskeyIsValid = true

				connectButton.SetEnabled(false)
				refreshButton.SetEnabled(true)
				reconnectButton.SetEnabled(true)
				disconnectButton.SetEnabled(true)

				rclone.CoreTransferring()

			} else {
				actionButton.SetEnabled(true)
				actionButton.SetCaption("Failed")
				info = "Authentication failed - wrong access key"
				Notify(infoLabel, info)

				states.Status = states.StatusError
				states.AcceskeyIsValid = false

				go func() {
					rclone.UnmountAll()
				}()

				connectButton.SetEnabled(true)
				refreshButton.SetEnabled(false)
				reconnectButton.SetEnabled(false)
				disconnectButton.SetEnabled(false)
			}
		}()
	})

	refreshButton.SetCallback(func() {
		popup.SetVisible(false)
		actionButton.SetEnabled(false)

		go func() {
			connectButton.SetEnabled(false)
			refreshButton.SetEnabled(false)
			reconnectButton.SetEnabled(false)
			disconnectButton.SetEnabled(false)

			actionButton.SetCaption("Refreshing..")
			info = "Refreshing local folders..."
			Notify(infoLabel, info)

			states.Status = states.StatusWait

			rclone.Refresh()

			actionButton.SetCaption("Connected")
			info = "Refreshing successfull"
			Notify(infoLabel, info)

			refreshButton.SetEnabled(true)
			reconnectButton.SetEnabled(true)
			disconnectButton.SetEnabled(true)

			actionButton.SetEnabled(true)
		}()
	})

	reconnectButton.SetCallback(func() {
		actionButton.SetEnabled(false)
		popup.SetVisible(false)

		go func() {
			connectButton.SetEnabled(false)
			refreshButton.SetEnabled(false)
			reconnectButton.SetEnabled(false)
			disconnectButton.SetEnabled(false)

			info = "Reconnecting in progress.."
			Notify(infoLabel, info)
			actionButton.SetCaption("Reconnecting..")
			states.Status = states.StatusWait

			rclone.UnmountAll()

			keyStatus := rclone.Main(false)
			//log.Println(keyStatus)

			if keyStatus {

				actionButton.SetEnabled(true)
				actionButton.SetCaption("Connected")
				info = "Successfully connected"
				Notify(infoLabel, info)
				states.Status = states.StatusSuccess
				states.AcceskeyIsValid = true

				connectButton.SetEnabled(false)
				refreshButton.SetEnabled(true)
				reconnectButton.SetEnabled(true)
				disconnectButton.SetEnabled(true)
			} else {
				actionButton.SetEnabled(true)
				actionButton.SetCaption("Failed")
				info = "Authentication failed - wrong access key"
				Notify(infoLabel, info)
				states.Status = states.StatusError
				states.AcceskeyIsValid = false

				go func() {
					rclone.UnmountAll()
				}()

				connectButton.SetEnabled(true)
				refreshButton.SetEnabled(false)
				reconnectButton.SetEnabled(false)
				disconnectButton.SetEnabled(false)
			}

			rclone.Refresh()

			actionButton.SetEnabled(true)
		}()

	})

	disconnectButton.SetCallback(func() {
		actionButton.SetEnabled(false)
		popup.SetVisible(false)

		go func() {
			connectButton.SetEnabled(true)
			refreshButton.SetEnabled(false)
			reconnectButton.SetEnabled(false)
			disconnectButton.SetEnabled(false)

			// Get permissions
			log.Println("[gfxMain] connectButton")

			actionButton.SetCaption("Disconnecting..  ")
			info = "Trying to unmount drives.."
			Notify(infoLabel, info)
			states.Status = states.StatusWait

			rclone.UnmountAll()

			actionButton.SetCaption("Diconnected")
			info = "Successfully disconnected"
			Notify(infoLabel, info)
			states.Status = states.StatusDefault

			actionButton.SetEnabled(true)
		}()
	})

	nanogui.NewLabel(menuWidget, "")
	nanogui.NewLabel(menuWidget, "")

	settings := nanogui.NewButton(menuWidget, "Settings")
	settings.SetCallback(func() {
		MainNGWindow.SetVisible(false)
		TransfersNGWindow.SetVisible(false)
		EmptyTransfersNGWindow.SetVisible(false)
		SettingsNKWindow.SetVisible(true)
	})

	about := nanogui.NewButton(menuWidget, "About")
	about.SetCallback(func() {
		MainNGWindow.SetVisible(false)
		TransfersNGWindow.SetVisible(false)
		EmptyTransfersNGWindow.SetVisible(false)
		AboutNKWindow.SetVisible(true)
	})

	quit := nanogui.NewButton(menuWidget, "Quit")
	quit.SetCallback(func() {
		log.Println("[gui_ng/main.go][gfxMain] Quit button pressed!")
		rclone.Quit()
		os.Setenv("MD_PID", "")
		os.Exit(1)
	})

	infoWidget := nanogui.NewWidget(MainNGWindow)
	infoWidget.SetLayout(nanogui.NewBoxLayout(nanogui.Vertical, nanogui.Fill, 12, 0))
	//infoLabel = nanogui.NewTextBox(infoWidget, info)
	infoLabel = nanogui.NewLabel(infoWidget, info)

	infoLabel.SetClampWidth(true)
	infoLabel.SetFontSize(fontSize)
	infoLabel.SetColor(nanovgo.Color{1, 1, 1, 1})

	EmptyTransfersNGWindow = nanogui.NewWindow(screen, "Transferings")

	EmptyTransfersNGWindow.SetPosition(0, splitterHeight)
	EmptyTransfersNGWindow.SetFixedWidth(winWidth)
	EmptyTransfersNGWindow.SetFixedHeight(winHeight - splitterHeight)
	EmptyTransfersNGWindow.SetLayout(nanogui.NewGroupLayout(0, 0, 0))
	EmptyTransfersNGWindow.SetVisible(true)
	EmptyTransfersNGWindow.SetDraggable(false)

	EmptyLabelWidget := nanogui.NewWidget(EmptyTransfersNGWindow)
	EmptyLabelWidget.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Fill, 10, 10))
	EmptyLabelWidget.SetPosition(0, 100)
	nanogui.NewLabel(EmptyLabelWidget, "No activity")

	TransfersNGWindow = nanogui.NewWindow(screen, "Transferings")

	TransfersNGWindow.SetPosition(0, splitterHeight)
	TransfersNGWindow.SetFixedWidth(winWidth)
	TransfersNGWindow.SetFixedHeight(winHeight - splitterHeight)
	TransfersNGWindow.SetLayout(nanogui.NewGroupLayout(0, 5, 0))
	TransfersNGWindow.SetVisible(false)
	TransfersNGWindow.SetDraggable(false)
	TransfersNGWindow.SetFocused(false)

	// To create scrollable widget we need:
	// 1. Parent widget
	// 	  2. Scroll widget
	// 		3. Child widget
	//           .. add elements here
	//           .. add elements here
	//           .. add elements here

	progressWidget := nanogui.NewWidget(TransfersNGWindow)
	progressWidget.SetLayout(nanogui.NewBoxLayout(nanogui.Vertical, nanogui.Fill, 0, 10))

	vscrollWidget := nanogui.NewVScrollPanel(progressWidget)
	vscrollLayout := nanogui.NewBoxLayout(nanogui.Vertical, nanogui.Fill, 5, 10)
	vscrollWidget.SetLayout(vscrollLayout)
	vscrollWidget.SetFixedHeight(winHeight)
	vscrollWidget.SetFixedWidth(winWidth)

	listWidget := nanogui.NewWidget(vscrollWidget)
	listWidget.SetLayout(nanogui.NewGridLayout(nanogui.Horizontal, 1, nanogui.Fill, 10, 10))

	// TODO: test data should be separated
	testArr := make([]rclone.Transferring, 6)
	if conf.IsDev {
		// Test data
		testData1 := rclone.Transferring{"Test 1", 2195316919, 592445440, 26, 1.196727107058912e+07, 1.153178700901072e+07}
		testData2 := rclone.Transferring{"Test 2", 1514794462, 26214400, 41, 1.590787626902579e+07, 4.196312355093148e+06}
		testData3 := rclone.Transferring{"Test 3", 2195316919, 592445440, 26, 1.196727107058912e+07, 1.153178700901072e+07}
		testData4 := rclone.Transferring{"Test 4", 1514794462, 26214400, 41, 1.590787626902579e+07, 4.196312355093148e+06}
		testData5 := rclone.Transferring{"Test 5", 2195316919, 592445440, 26, 1.196727107058912e+07, 1.153178700901072e+07}
		testData6 := rclone.Transferring{"Test 6", 1514794462, 26214400, 41, 1.590787626902579e+07, 4.196312355093148e+06}
		testArr[0] = testData1
		testArr[1] = testData2
		testArr[2] = testData3
		testArr[3] = testData4
		testArr[4] = testData5
		testArr[5] = testData6
	}

	// Create list of transfers
	maxLimit := 100
	tr := make([]Transfer, maxLimit)

	for i := 0; i < maxLimit; i++ {
		tr[i].Label = nanogui.NewLabel(listWidget, "Placeholder")
		tr[i].ProgressBar = nanogui.NewProgressBar(listWidget)
	}

	// Update key (Cron)
	// Setup CRON on Schedule
	log.Println("[gui_ng/main.go][gfxMain] Create Cron for keyaccess validation")
	cr := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	cr.AddFunc("@every "+conf.WebhookInterval, func() {
		if states.Status == states.StatusSuccess && states.AcceskeyIsValid {
			log.Println("[gui_ng/main.go][gfxMain] (re) Execute Cron for keyaccess validation")
			//rclone.Main(true)
			keyStatus := rclone.Main(true)
			//log.Println(keyStatus)

			if keyStatus {
				actionButton.SetEnabled(true)
				actionButton.SetCaption("Connected")
				info = "Successfully connected"
				Notify(infoLabel, info)

				states.Status = states.StatusSuccess
				states.AcceskeyIsValid = true

				connectButton.SetEnabled(false)
				refreshButton.SetEnabled(true)
				reconnectButton.SetEnabled(true)
				disconnectButton.SetEnabled(true)

				rclone.CoreTransferring()

			} else {
				actionButton.SetEnabled(true)
				actionButton.SetCaption("Failed")
				info = "Authentication failed - wrong access key"
				Notify(infoLabel, info)
				states.Status = states.StatusError
				states.AcceskeyIsValid = false

				go func() {
					rclone.UnmountAll()
				}()

				connectButton.SetEnabled(true)
				refreshButton.SetEnabled(false)
				reconnectButton.SetEnabled(false)
				disconnectButton.SetEnabled(false)
			}
		}
	})
	log.Println("[gui_ng/main.go][gfxMain] Run Cron for keyaccess validation")
	go cr.Run()

	// Transferings (Cron)
	// Cron update
	log.Println("[gui_ng/main.go][gfxMain] Create Cron for transfers update")
	c = cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	c.AddFunc("@every 2s", func() {
		log.Println("[gui_ng/main.go][gfxMain] (re) Execute Cron for transfers update")
		rawResp := rclone.CoreTransferring().Each

		resp := make([]rclone.Transferring, len(rawResp))

		// Prepare response (остановился тут)
		for i := 0; i < len(rawResp); i++ {
			if rawResp[i].Percentage > 0 {
				resp[i].Bytes = rawResp[i].Bytes
				resp[i].Name = rawResp[i].Name
				resp[i].Percentage = rawResp[i].Percentage
				resp[i].Size = rawResp[i].Size
				resp[i].Speed = rawResp[i].Speed
				resp[i].SpeedAvg = rawResp[i].SpeedAvg
			}
		}

		// Process the response
		if len(resp) > 0 {
			TransfersNGWindow.SetVisible(true)
			EmptyTransfersNGWindow.SetVisible(false)
			states.Status = states.StatusWait

			// Update values
			for i := 0; i < len(resp); i++ {
				tr[i].Label.SetCaption(resp[i].Name)
				//tr[i].Label.SetWidth(winWidth - 45)

				var floatPercent float32
				floatPercent = float32(resp[i].Percentage) / 100
				tr[i].ProgressBar.SetValue(floatPercent)
				//tr[i].ProgressBar.SetWidth(winWidth - 45)

			}

			// Hide all except new
			for i := len(resp); i < maxLimit; i++ {
				tr[i].Label.SetVisible(false)
				tr[i].ProgressBar.SetVisible(false)
			}

			// Show new only
			for i := 0; i < len(resp); i++ {
				tr[i].Label.SetVisible(true)
				tr[i].ProgressBar.SetVisible(true)
			}
			fmt.Println(len(resp))

		} else {
			TransfersNGWindow.SetVisible(false)
			if MainNGWindow.Visible() {
				EmptyTransfersNGWindow.SetVisible(true)
			}

			if states.AcceskeyIsValid {
				states.Status = states.StatusSuccess
			}
		}
	})
	log.Println("[gui_ng/main.go][gfxMain] Run Cron for transfers update")
	go c.Run()

	// TODO: test data should be separated
	if conf.IsDev {
		// Manual Update
		update := nanogui.NewButton(infoWidget, "Update Real Data")
		//update.SetFixedWidth(100)
		update.SetVisible(true)
		update.SetCallback(func() {
			resp := rclone.CoreTransferring().Each

			if len(resp) > 0 {
				TransfersNGWindow.SetVisible(true)
				EmptyTransfersNGWindow.SetVisible(false)
				states.Status = states.StatusWait

				// Update values
				for i := 0; i < len(resp); i++ {
					tr[i].Label.SetCaption(resp[i].Name)
					//tr[i].Label.SetWidth(winWidth - 45)

					var floatPercent float32
					floatPercent = float32(resp[i].Percentage) / 100
					tr[i].ProgressBar.SetValue(floatPercent)
					//tr[i].ProgressBar.SetWidth(winWidth - 45)
				}

				// Hide all except new
				for i := len(resp); i < maxLimit; i++ {
					tr[i].Label.SetVisible(false)
					tr[i].ProgressBar.SetVisible(false)
				}

				// Show new only
				for i := 0; i < len(resp); i++ {
					tr[i].Label.SetVisible(true)
					tr[i].ProgressBar.SetVisible(true)
				}
				fmt.Println(len(resp))

			} else {
				TransfersNGWindow.SetVisible(false)
				if MainNGWindow.Visible() {
					EmptyTransfersNGWindow.SetVisible(true)
				}

				if states.AcceskeyIsValid {
					states.Status = states.StatusSuccess
				}
			}
		})
		// Manual Update
		test := nanogui.NewButton(infoWidget, "Update Test Data")
		//update.SetFixedWidth(100)
		test.SetVisible(true)
		test.SetCallback(func() {
			resp := testArr

			if len(resp) > 0 {
				TransfersNGWindow.SetVisible(true)
				EmptyTransfersNGWindow.SetVisible(false)
				states.Status = states.StatusWait

				// Update values
				for i := 0; i < len(resp); i++ {
					tr[i].Label.SetCaption(resp[i].Name)
					//tr[i].Label.SetWidth(winWidth - 45)

					var floatPercent float32
					floatPercent = float32(resp[i].Percentage) / 100
					tr[i].ProgressBar.SetValue(floatPercent)
					//tr[i].ProgressBar.SetWidth(winWidth - 45)
				}

				// Hide all except new
				for i := len(resp); i < maxLimit; i++ {
					tr[i].Label.SetVisible(false)
					tr[i].ProgressBar.SetVisible(false)
				}

				// Show new only
				for i := 0; i < len(resp); i++ {
					tr[i].Label.SetVisible(true)
					tr[i].ProgressBar.SetVisible(true)
				}
				fmt.Println(len(resp))

			} else {
				TransfersNGWindow.SetVisible(false)
				if MainNGWindow.Visible() {
					EmptyTransfersNGWindow.SetVisible(true)
				}

				if states.AcceskeyIsValid {
					states.Status = states.StatusSuccess
				}
			}
		})
	}

}

func gfxSettings(screen *nanogui.Screen) {
	SettingsNKWindow = nanogui.NewWindow(screen, "Settings")
	SettingsNKWindow.SetPosition(0, 0)
	SettingsNKWindow.SetFixedWidth(winWidth)
	SettingsNKWindow.SetFixedHeight(winHeight)
	SettingsNKWindow.SetLayout(nanogui.NewGroupLayout(0, 5, 5))
	SettingsNKWindow.SetVisible(false)

	contentWidget := nanogui.NewWidget(SettingsNKWindow)
	contentLayout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Fill, 10, 5)
	contentWidget.SetLayout(contentLayout)

	// Access key
	accesskeyLabel := nanogui.NewLabel(contentWidget, "Access key")
	accesskeyLabel.SetFontSize(fontSize)

	var accesskeyValue string
	accesskeyInput := nanogui.NewTextBox(contentWidget, "")
	accesskeyInput.SetEditable(true)
	accesskeyInput.SetCallback(func(input string) bool {
		fmt.Println(input)
		accesskeyValue = input
		return true
	})

	/*
		// Connect at startup
		conatstartupLabel := nanogui.NewLabel(contentWidget, "Mount at program start")
		conatstartupLabel.SetFontSize(fontSize)

		conatstartupInput := nanogui.NewCheckBox(contentWidget, "")
		conatstartupInput.SetCallback(func(checked bool) {
			fmt.Println("Check box 1 state:", checked)
		})

		// Portable mode
		portableLabel := nanogui.NewLabel(contentWidget, "Portable mode")
		portableLabel.SetFontSize(fontSize)

		portableInput := nanogui.NewCheckBox(contentWidget, "")
		portableInput.SetCallback(func(checked bool) {
			fmt.Println("Check box 1 state:", checked)
		})
	*/

	// Footer
	buttonsWidget := nanogui.NewWidget(SettingsNKWindow)
	buttonsWidget.SetLayout(nanogui.NewGridLayout(nanogui.Horizontal, 7, nanogui.Fill, 10, 5))

	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")

	save := nanogui.NewButton(buttonsWidget, "Save")
	save.SetCallback(func() {
		//fmt.Println("Access key is " + accesskeyValue)
		if len(accesskeyValue) > 0 {
			settings.WriteAccessKey(accesskeyValue)
			accesskeyInput.SetValue("")

			//info = "access key saved. Try to connect"
			//Notify(infoLabel, info)

			SettingsNKWindow.SetVisible(false)
			MainNGWindow.SetVisible(true)
			TransfersNGWindow.SetVisible(false)
			EmptyTransfersNGWindow.SetVisible(true)
		} else {
			//info = "Warning: can't save empty access key"
			//Notify(infoLabel, info)
		}
	})

	cancel := nanogui.NewButton(buttonsWidget, "Cancel")
	cancel.SetCallback(func() {
		SettingsNKWindow.SetVisible(false)
		MainNGWindow.SetVisible(true)
		TransfersNGWindow.SetVisible(false)
		EmptyTransfersNGWindow.SetVisible(true)
	})

}

func gfxAbout(screen *nanogui.Screen) {
	content := conf.AppName + ` - mount S3 buckets as network drives`
	version := "version " + conf.AppVersion + " (rclone " + conf.RcloneVersion + ")"
	homepageURL := "https://sites.google.com/mirball.com/mirball-drives"
	author := "© 2021 Innokentiy Sokolov"
	authorURL := "GitHub: https://github.com/keshon"
	authorEmail := "Email: keshon@zoho.com"

	AboutNKWindow = nanogui.NewWindow(screen, "About")
	AboutNKWindow.SetPosition(0, 0)
	AboutNKWindow.SetFixedWidth(winWidth)
	AboutNKWindow.SetFixedHeight(winHeight)
	AboutNKWindow.SetLayout(nanogui.NewGroupLayout(0, 5, 5))
	AboutNKWindow.SetVisible(false)

	// Content
	contentWidget := nanogui.NewWidget(AboutNKWindow)
	contentLayout := nanogui.NewGridLayout(nanogui.Horizontal, 1, nanogui.Fill, 10, 5)
	contentWidget.SetLayout(contentLayout)

	contentLabel := nanogui.NewLabel(contentWidget, content)
	contentLabel.SetFontSize(fontSize)

	versionLabel := nanogui.NewLabel(contentWidget, version)

	nanogui.NewLabel(contentWidget, " ")

	authorLabel := nanogui.NewLabel(contentWidget, author)
	authorURLLabel := nanogui.NewLabel(contentWidget, authorURL)
	authorEmailLabel := nanogui.NewLabel(contentWidget, authorEmail)

	versionLabel.SetFontSize(fontSize)
	authorLabel.SetFontSize(fontSize)
	authorURLLabel.SetFontSize(fontSize)
	authorEmailLabel.SetFontSize(fontSize)

	// Footer
	buttonsWidget := nanogui.NewWidget(AboutNKWindow)
	buttonsWidget.SetLayout(nanogui.NewGridLayout(nanogui.Horizontal, 7, nanogui.Fill, 10, 5))

	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")
	nanogui.NewLabel(buttonsWidget, "")

	www := nanogui.NewButton(buttonsWidget, "Open website")
	www.SetCallback(func() {
		utils.Openbrowser(homepageURL)
		ShowMainWindow = false
		MainWindowIsHidden = true
		AboutNKWindow.SetVisible(false)
		MainNGWindow.SetVisible(true)
		TransfersNGWindow.SetVisible(false)
		EmptyTransfersNGWindow.SetVisible(true)
		app.screen.Window.Hide()
	})

	ok := nanogui.NewButton(buttonsWidget, "OK")
	ok.SetCallback(func() {
		AboutNKWindow.SetVisible(false)
		MainNGWindow.SetVisible(true)
		TransfersNGWindow.SetVisible(false)
		EmptyTransfersNGWindow.SetVisible(true)
	})
}
