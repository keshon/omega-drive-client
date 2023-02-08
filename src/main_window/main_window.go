package main_window

import (
	"app/src/rcd"
	"app/src/state"
	"app/src/utils"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/go-cmp/cmp"
	"github.com/robfig/cron/v3"
)

/*
	This package contains main window user interface
*/

// Main window
func MakeMainWindow(w fyne.Window) fyne.CanvasObject {
	// Testing
	//test_gui.MakeFakeSyncingData()
	//test_gui.MakeFakeHistoryData()

	// Toolbar
	toolbar := container.NewHBox()
	toolbarScroll := container.NewHScroll(toolbar)

	// Syncing and history
	syncing := container.NewVBox()
	history := container.NewVBox()
	transferScroll := container.NewScroll(container.NewVBox(syncing, history))

	// Footer
	statusDesc := widget.NewMultiLineEntry()
	statusDesc.Hide()
	statusDesc.Bind(state.ConnectionStatus.BindDescription)
	statusDesc.Wrapping = fyne.TextWrapWord

	toggle := true
	showMoreLabel := "More"
	showMore := widget.NewButton(showMoreLabel, func() {
		if toggle {
			statusDesc.Show()
			toggle = false
		} else {
			statusDesc.Hide()
			toggle = true
		}
	})

	state.ConnectionStatus.BindTitle.Set("Idle")
	statusTitle := widget.NewLabelWithData(state.ConnectionStatus.BindTitle)
	statusbar := container.NewBorder(nil, nil, statusTitle, showMore)

	connection := container.NewVBox(statusDesc, statusbar)
	connectionScroll := container.NewHScroll(connection)

	// Final layout
	finalLayout := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbarScroll, connectionScroll, nil, nil), toolbarScroll, transferScroll, connectionScroll)

	// Update loop
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	var toolbarBtns []fyne.CanvasObject

	var canvasSyncList []fyne.CanvasObject
	var canvasHistoryList []fyne.CanvasObject

	var prevHistoryStamps []string
	var prevHistoryData []state.HistoryDataStruct

	c.AddFunc("@every 1s", func() {

		// Toolbar
		if state.ConnectionStatus.Status == state.Connected {
			// Get array of drive names
			var availPaths []string
			var activPaths []string

			for _, elem := range state.Response.AvailablePaths {
				availPaths = append(availPaths, elem.Name)
			}
			for _, elem := range state.ActivePaths {
				activPaths = append(activPaths, elem.Name)
			}

			//fmt.Println(availPaths)
			//fmt.Println(activPaths)
			//fmt.Println(len(toolbarBtns))

			if len(activPaths) > 0 {
				// Fill-in toolbar
				if len(toolbarBtns) == 0 {
					if (len(toolbarBtns) == 0 && cmp.Equal(availPaths, activPaths)) || (len(toolbarBtns) > 0 && !cmp.Equal(availPaths, activPaths)) {
						for i := 0; i < len(state.ActivePaths); i++ {
							path := state.ActivePaths[i].Letter + ":\\"
							elem := widget.NewButtonWithIcon(state.ActivePaths[i].Name, theme.FolderOpenIcon(), func() {
								rcd.Rcd.Refresh()
								cmd := exec.Command("explorer", path)
								cmd.Run()
							})
							toolbarBtns = append(toolbarBtns, elem)
							toolbar.Add(elem)
						}

						toolbar.Refresh()

						fmt.Println(cmp.Equal(state.Response.AvailablePaths, state.ActivePaths))
						fmt.Println("Toolbar changed")
					}
				}
			} else {
				// Clear tooldbar
				if len(toolbarBtns) > 0 {
					for _, elem := range toolbarBtns {
						toolbar.Remove(elem)
					}
					toolbarBtns = nil
				}
			}
		}

		// Syncing list
		/*
			fmt.Println(len(canvasSyncList))
			fmt.Println(len(state.SyncingData))
			fmt.Println("-----")
		*/
		if len(canvasSyncList) > 0 && len(state.SyncingData) == 0 {
			fmt.Println("syncing.Refresh()")
			for _, elem := range canvasSyncList {
				syncing.Remove(elem)
			}
			canvasSyncList = nil
			syncing.Refresh()
		}

		if len(state.SyncingData) > 0 {
			canvasSyncList = nil
			// fill in
			for i := 0; i < len(state.SyncingData); i++ {
				elem := makeSyncingProgress(state.SyncingData[i].Label, state.SyncingData[i].Status, state.SyncingData[i].Progress, w)
				canvasSyncList = append(canvasSyncList, elem)
				//syncing = append(syncing, elem)
				//uploads.Add(elem)
			}
			syncing.Objects = canvasSyncList

			fmt.Println("syncing.Refresh()")
			syncing.Refresh()
		}

		// History list
		if len(canvasHistoryList) > 0 {
			history.Refresh()
		}
		canvasHistoryList = nil

		var currentHistoryStamps []string

		for _, elem := range state.HistoryData {
			currentHistoryStamps = append(currentHistoryStamps, elem.Label)
		}

		diff := utils.StringArrayDifference(currentHistoryStamps, prevHistoryStamps)
		if len(diff) > 0 {
			// fill-in
			for k, elem := range diff {
				prevHistoryStamps = append(prevHistoryStamps, elem)
				prevHistoryData = append(prevHistoryData, state.HistoryData[k])
			}

			// reverse
			reversedData := []state.HistoryDataStruct{}
			for i := len(prevHistoryData) - 1; i >= 0; i-- {
				reversedData = append(reversedData, prevHistoryData[i])
			}

			// fill in

			for i := 0; i < len(reversedData); i++ {
				elem := makeHistoryStats(reversedData[i].Label, reversedData[i].Status, w)
				canvasHistoryList = append(canvasHistoryList, elem)
			}
			history.Objects = canvasHistoryList

			fmt.Println("history.Refresh()")
			history.Refresh()
		}

		/*
			fmt.Println("diff")
			fmt.Println(diff)
			fmt.Println("currentHistoryStamps")
			fmt.Println(currentHistoryStamps)
			fmt.Println("prevHistoryStamps")
			fmt.Println(prevHistoryStamps)
			fmt.Println("prevHistoryData")
			fmt.Println(prevHistoryData)
			fmt.Println("-----")
		*/

	})
	go c.Run()

	return finalLayout
}

func makeSyncingProgress(label, summary string, progress float64, w fyne.Window) fyne.CanvasObject {
	var scale float32
	scale = 0.7

	// File icon
	icon := canvas.NewImageFromResource(theme.FileIcon())
	icon.SetMinSize(fyne.NewSize(40, 40))

	// Filename (e.g. "filename.dat")
	file := filepath.Base(label)
	filename := canvas.NewText(file, theme.ForegroundColor())
	filename.Alignment = fyne.TextAlignLeading
	filename.TextSize = theme.TextSize() * (scale + 0.2)
	filename.TextStyle = fyne.TextStyle{Bold: true}

	// Uploading progress as mini progress bar
	// using rectangle
	_, fraction := math.Modf(progress)
	uploaded := float32(fraction) * 100 * w.Canvas().Size().Width / 100

	progressBar := canvas.NewRectangle(theme.PrimaryColor())
	progressBar.SetMinSize(fyne.NewSize(1, 1))
	progressBar.Resize(fyne.NewSize(uploaded, 1))

	barContainer := container.NewWithoutLayout(progressBar)

	status := canvas.NewText(summary, theme.ForegroundColor())
	status.Alignment = fyne.TextAlignLeading
	status.TextSize = theme.TextSize() * scale

	info := container.NewVBox(filename, barContainer, status)

	return container.NewVBox(
		container.NewBorder(nil, nil, icon, nil, info),
		makePadding(), makePadding(),
	)
}

func makeHistoryStats(label, summary string, w fyne.Window) fyne.CanvasObject {
	var scale float32
	scale = 0.7

	// File icon
	icon := canvas.NewImageFromResource(theme.FileIcon())
	icon.SetMinSize(fyne.NewSize(40, 40))

	// Filename (e.g. "filename.dat")
	file := filepath.Base(label)
	filename := canvas.NewText(file, theme.ForegroundColor())
	filename.Alignment = fyne.TextAlignLeading
	filename.TextSize = theme.TextSize() * (scale + 0.2)
	filename.TextStyle = fyne.TextStyle{Bold: true}

	progressBar := canvas.NewRectangle(theme.BackgroundColor())
	progressBar.SetMinSize(fyne.NewSize(0, 0))
	progressBar.Resize(fyne.NewSize(0, 0))

	barContainer := container.NewWithoutLayout(progressBar)

	status := canvas.NewText(summary, theme.ForegroundColor())
	status.Alignment = fyne.TextAlignLeading
	status.TextSize = theme.TextSize() * scale

	info := container.NewVBox(filename, barContainer, status)

	return container.NewVBox(
		container.NewBorder(nil, nil, icon, nil, info),
		makePadding(), makePadding(),
	)
}

// Padding
func makePadding() fyne.CanvasObject {
	rect := canvas.NewRectangle(theme.BackgroundColor())
	//rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 255})
	rect.SetMinSize(fyne.NewSize(2, 2))
	return rect
}

// Custom context menu for a button
type contextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

func (b *contextMenuButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuButton(label string, icon fyne.Resource, menu *fyne.Menu) *contextMenuButton {
	b := &contextMenuButton{menu: menu}
	b.Text = label
	b.Icon = icon
	b.ExtendBaseWidget(b)
	return b
}
