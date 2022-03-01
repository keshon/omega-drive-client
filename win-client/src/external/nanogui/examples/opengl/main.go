// +build !js

package main

import (
	"app/src/external/nanogui"
	"app/src/external/nanogui/examples/opengl/ui"

	"github.com/go-gl/mathgl/mgl64"
	glu "github.com/maxfish/gl_utils/gl_utils"
)

type Application struct {
	screen *nanogui.Screen
	camera *glu.Camera2D
	quad   *glu.Primitive2D
}

var app Application

func (app *Application) init() {
	app.screen = nanogui.NewScreen(1024, 768, "Mixing OpenGL Test", true, false)

	ui.GridPanel(app.screen, 20, 20)

	app.screen.SetDrawContentsCallback(drawContent)
	app.screen.PerformLayout()
	app.screen.DebugPrint()

	// Content
	app.camera = glu.NewCamera2D(app.screen.Width(), app.screen.Height(), 1)
	app.camera.SetCentered(true)
	app.quad = glu.NewRegularPolygonPrimitive(mgl64.Vec3{0, 0, 0.5}, 200, 9, true)
	app.quad.SetAnchorToCenter()
	app.quad.SetColor(glu.Color{1, 1, 1, 1})
}

func drawContent() {
	app.quad.SetAngle(float64(nanogui.GetTime()) / 4)
	app.quad.Draw(app.camera.ProjectionMatrix32())
}

func main() {
	nanogui.Init()
	//nanogui.SetDebug(true)

	app = Application{}
	app.init()
	app.screen.DrawAll()
	app.screen.SetVisible(true)
	nanogui.MainLoop()
}
