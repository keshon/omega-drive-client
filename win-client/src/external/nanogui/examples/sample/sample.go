// +build !js

package main

import (
	"app/src/external/nanogui"
	"app/src/external/nanogui/examples/sample/demo"
	"app/src/external/nanovgo"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	_ "net/http/pprof"
	"path"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Application struct {
	screen   *nanogui.Screen
	progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

func (a *Application) init() {
	go func() {
		http.ListenAndServe(":5555", http.DefaultServeMux)
	}()

	glfw.WindowHint(glfw.Samples, 4)

	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)
	a.screen.NVGContext().CreateFont("japanese", "font/GenShinGothic-P-Regular.ttf")

	demo.ButtonDemo(a.screen)
	images := loadImageDirectory(a.screen.NVGContext(), "icons")
	imageButton, imagePanel, progressBar := demo.BasicWidgetsDemo(a.screen, images)
	a.progress = progressBar
	demo.SelectedImageDemo(a.screen, imageButton, imagePanel)
	demo.MiscWidgetsDemo(a.screen)
	demo.GridDemo(a.screen)

	a.screen.SetDrawContentsCallback(func() {
		a.progress.SetValue(float32(math.Mod(float64(nanogui.GetTime())/10, 1.0)))
	})

	a.screen.PerformLayout()
	a.screen.DebugPrint()

	/* All NanoGUI widgets are initialized at this point. Now
	create an OpenGL shader to draw the main window contents.
	*/
}

func main() {
	nanogui.Init()

	//nanogui.SetDebug(true)

	app := Application{}
	app.init()
	app.screen.DrawAll()
	app.screen.SetVisible(true)
	nanogui.MainLoop()
}

func loadImageDirectory(ctx *nanovgo.Context, dir string) []nanogui.Image {
	var images []nanogui.Image
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("loadImageDirectory: read error %v\n", err))
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := path.Ext(file.Name())
		if ext != ".png" {
			continue
		}
		fullPath := path.Join(dir, file.Name())
		handle, image := ctx.CreateImage(fullPath, 0)
		if handle == 0 {
			panic("Could not open image data!")
		}
		images = append(images, nanogui.Image{
			ImageID:   handle,
			Name:      fullPath[:len(fullPath)-4],
			ImageData: &image,
		})
	}
	return images
}
