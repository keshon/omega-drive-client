package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/maxfish/vg4go-gl4"
	"github.com/maxfish/vg4go-gl4/perfgraph"
	"github.com/maxfish/vg4go-gl4/sample/demo"
	"log"
	"runtime"
)

func initOpenGL() {
	runtime.LockOSThread()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(1000, 600, "NanoVGo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(key)
	window.MakeContextCurrent()

	initOpenGL()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	ctx, err := nanovgo.NewContext(nanovgo.Debug /*nanovgo.AntiAlias | nanovgo.StencilStrokes | nanovgo.Debug*/)
	defer ctx.Delete()

	if err != nil {
		panic(err)
	}

	demoData := LoadDemo(ctx)

	glfw.SwapInterval(0)

	fps := perfgraph.NewPerfGraph("Frame Time", "sans")

	for !window.ShouldClose() {
		t, _ := fps.UpdateGraph()

		fbWidth, fbHeight := window.GetFramebufferSize()
		winWidth, winHeight := window.GetSize()
		mx, my := window.GetCursorPos()

		pixelRatio := float32(fbWidth) / float32(winWidth)
		gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
		if premult {
			gl.ClearColor(0, 0, 0, 0)
		} else {
			gl.ClearColor(0.3, 0.3, 0.32, 1.0)
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.Enable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)

		ctx.BeginFrame(winWidth, winHeight, pixelRatio)

		demo.RenderDemo(ctx, float32(mx), float32(my), float32(winWidth), float32(winHeight), t, blowup, demoData)
		fps.RenderGraph(ctx, 5, 5)

		ctx.EndFrame()

		gl.Enable(gl.DEPTH_TEST)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	demoData.FreeData(ctx)
}

var blowup bool
var premult bool

func key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	} else if key == glfw.KeySpace && action == glfw.Press {
		blowup = !blowup
	} else if key == glfw.KeyP && action == glfw.Press {
		premult = !premult
	}
}

func LoadDemo(ctx *nanovgo.Context) *demo.DemoData {
	d := &demo.DemoData{}
	for i := 0; i < 12; i++ {
		path := fmt.Sprintf("images/image%d.jpg", i+1)
		handle, _ := ctx.CreateImage(path, 0)
		d.Images = append(d.Images, handle)
		if d.Images[i] == 0 {
			log.Fatalf("Could not load %s", path)
		}
	}

	d.FontIcons = ctx.CreateFont("icons", "entypo.ttf")
	if d.FontIcons == -1 {
		log.Fatalln("Could not add font icons.")
	}
	d.FontNormal = ctx.CreateFont("sans", "Roboto-Regular.ttf")
	if d.FontNormal == -1 {
		log.Fatalln("Could not add font italic.")
	}
	d.FontBold = ctx.CreateFont("sans-bold", "Roboto-Bold.ttf")
	if d.FontBold == -1 {
		log.Fatalln("Could not add font bold.")
	}
	return d
}
