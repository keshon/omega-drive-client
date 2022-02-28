package nanogui

import (
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var mainloopActive bool = false
var StartTime time.Time
var debugFlag bool

func Init() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	StartTime = time.Now()

	// -fix
	/*
		if err := gl.Init(); err != nil {
			panic(err)
		}
	*/
}

func GetTime() float32 {
	return float32(time.Now().Sub(StartTime)/time.Millisecond) * 0.001
}

func MainLoop() {
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
		haveActiveScreen := false
		for _, screen := range NanoguiScreens {
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
		glfw.WaitEvents()
	}

	wg.Wait()
}

func SetDebug(d bool) {
	debugFlag = d
}

func InitWidget(child, parent Widget) {
	//w.cursor = Arrow
	if parent != nil {
		parent.AddChild(parent, child)
		child.SetTheme(parent.Theme())
	}
	child.SetVisible(true)
	child.SetEnabled(true)
	child.SetFontSize(-1)
}
