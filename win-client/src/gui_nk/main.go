package gui_nk

import (
	"C"

	"app/src/settings"
	"app/src/utils"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
)

// Nuklear - alternative GUI
// Not finished

var (
	ShowMainWindow     bool
	MainWindowIsHidden bool

	ActiveNKWindow string

	MousePosX int
	MousePosY int
)

const (
	winWidth  = 500
	winHeight = 600

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024

	fontSize   = 19
	elemHeight = 31

	appName       = "Mirball Drives"
	appVersion    = "220216"
	rcloneVersion = "1.53"
)

type Option uint8

type State struct {
	bgColor   nk.Color
	prop      int32
	opt       Option
	accessKey nk.TextEdit
}

func init() {
	runtime.LockOSThread()
}

func Main() *glfw.Window {
	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	//glfw.WindowHint(glfw.Maximized, glfw.True)

	win, err := glfw.CreateWindow(winWidth, winHeight, "Nuklear Demo", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}
	win.MakeContextCurrent()

	width, height := win.GetSize()
	fmt.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	// sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
	//config := nk.NkFontConfig(14)
	//config.SetOversample(1, 1)
	//config.SetRange(nk.NkFontChineseGlyphRanges())
	//simsunFont := nk.NkFontAtlasAddFromFile(atlas, "c:\\Windows\\Fonts\\segoeuib.ttf", fontSize, &config)
	nk.NkFontStashEnd()
	/*
		if simsunFont != nil {
			nk.NkStyleSetFont(ctx, simsunFont.Handle())
		}
	*/

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	state := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}
	nk.NkTexteditInitDefault(&state.accessKey)

	fpsTicker := time.NewTicker(time.Second / 30)
	ActiveNKWindow = "main"
	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)
			return &glfw.Window{}
		case <-fpsTicker.C:
			if win.ShouldClose() {
				close(exitC)
				continue
			}
			// Toggle main window show/hide state
			if !ShowMainWindow {
				// We are in a loop so it should be one time toggle between states
				if !MainWindowIsHidden {
					MainWindowIsHidden = true
					win.Hide()
				}
			} else {
				if MainWindowIsHidden {
					MainWindowIsHidden = false
					win.Show()
					// Set window position
					scrWidth, _ := GetMonitorWidth()
					width, height := win.GetSize()
					winPosX := scrWidth - (width / 2) - (scrWidth - MousePosX)
					winPosY := MousePosY - height - 32
					fmt.Println(MousePosX)
					fmt.Println(MousePosY)
					win.SetPos(winPosX, winPosY)
				}
			}
			glfw.PollEvents()

			switch ActiveNKWindow {
			case "main":
				gfxMain(win, ctx, state)
			case "settings":
				gfxSettings(win, ctx, state)
			case "about":
				gfxAbout(win, ctx, state)
			default:
				gfxMain(win, ctx, state)
			}
		}
	}

}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(0, 0, 500, 600)
	update := nk.NkBegin(ctx, "Mirball Drives", bounds,
		nk.WindowBorder|nk.WindowTitle)

	if update > 0 {
		//nk.NkLayoutRowStatic(ctx, 30, 80, 1)
		nk.NkLayoutRowDynamic(ctx, elemHeight, 5)
		list := []string{"Connect", "Refresj", "Reconnect", "Dicsonnect"}
		{
			if nk.NkCombo(ctx, list, 4, 0, 40, nk.NkVec2(100, 300)) > 0 {
				fmt.Println("[gfxMain] Settings button pressed!")
				//ActiveNKWindow = "settings"
			}

			nk.NkSpacing(ctx, 1)

			if nk.NkButtonLabel(ctx, "Settings") > 0 {
				fmt.Println("[gfxMain] Settings button pressed!")
				ActiveNKWindow = "settings"
			}
			if nk.NkButtonLabel(ctx, "About") > 0 {
				fmt.Println("[gfxMain] About button pressed!")
				ActiveNKWindow = "about"
			}
			if nk.NkButtonLabel(ctx, "Quit") > 0 {
				fmt.Println("[gfxMain] Quit button pressed!")
				os.Exit(1)
			}
		}
		/*
			if (nk_contextual_begin(ctx, NK_WINDOW_NO_SCROLLBAR, nk_vec2(150, 300), nk_window_get_bounds(ctx))) {
				nk_layout_row_dynamic(ctx, 30, 1);
				if (nk_contextual_item_image_label(ctx, media->copy, "Clone", NK_TEXT_RIGHT))
					fprintf(stdout, "pressed clone!\n");
				if (nk_contextual_item_image_label(ctx, media->del, "Delete", NK_TEXT_RIGHT))
					fprintf(stdout, "pressed delete!\n");
				if (nk_contextual_item_image_label(ctx, media->convert, "Convert", NK_TEXT_RIGHT))
					fprintf(stdout, "pressed convert!\n");
				if (nk_contextual_item_image_label(ctx, media->edit, "Edit", NK_TEXT_RIGHT))
					fprintf(stdout, "pressed edit!\n");
				nk_contextual_end(ctx);
			}
		*/
		nk.NkLayoutRowDynamic(ctx, elemHeight, 1)
		/*
			static size_t prog = 80;
			ui_header(ctx, media, "Progressbar");
			ui_widget(ctx, media, 35);
			nk_progress(ctx, &prog, 100, nk_true);
		*/
		{
			var cur nk.Size
			cur = 46
			nk.NkLabel(ctx, "Transfers", nk.TextAlignLeft)
			nk.NkProgress(ctx, &cur, 100, nk.True)

		}
		/*
			{
				if nk.NkContextualBegin(ctx, nk.WindowNoScrollbar, nk.NkVec2(150, 300), nk.NkWindowGetBounds(ctx)) > 0 {
					nk.NkLayoutRowDynamic(ctx, 32, 1)
					if nk.NkContextualItemText(ctx, "text", 1, nk.TextRight) > 0 {
						fmt.Println("ok")
					}
					if nk.NkContextualItemText(ctx, "text", 1, nk.TextRight) > 0 {
						fmt.Println("ok")
					}
					if nk.NkContextualItemText(ctx, "text", 1, nk.TextRight) > 0 {
						fmt.Println("ok")
					}
				}
			}
		*/
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

func gfxSettings(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(0, 0, 400, 500)
	update := nk.NkBegin(ctx, "Settings", bounds,
		nk.WindowBorder|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowDynamic(ctx, elemHeight, 2)
		{
			nk.NkLabel(ctx, "Access key", nk.TextAlignLeft)
			nk.NkEditBuffer(ctx, nk.EditField, &state.accessKey, nk.NkFilterDefault)
		}
		nk.NkLayoutRowDynamic(ctx, elemHeight, 1)
		{
			nk.NkLabel(ctx, "", nk.TextAlignLeft)
		}
		nk.NkLayoutRowDynamic(ctx, elemHeight, 4)
		{
			nk.NkSpacing(ctx, 1)
			nk.NkSpacing(ctx, 1)
			if nk.NkButtonLabel(ctx, "Save") > 0 {
				fmt.Println("[gfxSettings] Save button pressed!")
				fmt.Println(state.accessKey.GetGoString())
				settings.WriteAccessKey(state.accessKey.GetGoString())
				nk.NkTexteditInitDefault(&state.accessKey)
				ActiveNKWindow = "main"
			}
			if nk.NkButtonLabel(ctx, "Cancel") > 0 {
				fmt.Println("[gfxSettings] Cancel button pressed!")
				ActiveNKWindow = "main"
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

func gfxAbout(win *glfw.Window, ctx *nk.Context, state *State) {
	content := appName + ` - mount S3 buckets as network drives`
	version := "version " + appVersion + " (rclone version " + rcloneVersion + ")"
	homepageURL := "https://sites.google.com/mirball.com/mirball-drives"
	author := "© 2021 Innokentiy Sokolov"
	authorURL := "GitHub: https://github.com/keshon"
	authorEmail := "Email: keshon@zoho.com"

	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(0, 0, 400, 500)
	update := nk.NkBegin(ctx, "About software", bounds,
		nk.WindowBorder|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowDynamic(ctx, elemHeight, 1)
		{
			nk.NkLabel(ctx, content, nk.TextAlignLeft)
			nk.NkLabel(ctx, version, nk.TextAlignLeft)
			nk.NkLabel(ctx, "", nk.TextAlignLeft)
			nk.NkLabel(ctx, author, nk.TextAlignLeft)
			nk.NkLabel(ctx, authorURL, nk.TextAlignLeft)
			nk.NkLabel(ctx, authorEmail, nk.TextAlignLeft)
			nk.NkLabel(ctx, "", nk.TextAlignLeft)
		}
		nk.NkLayoutRowDynamic(ctx, elemHeight, 3)
		{
			nk.NkSpacing(ctx, 1)

			if nk.NkButtonLabel(ctx, "Visit website") > 0 {
				fmt.Println("[gfxAbout] Visit website button pressed!")
				utils.Openbrowser(homepageURL)
				ShowMainWindow = false
				ActiveNKWindow = "main"
			}
			if nk.NkButtonLabel(ctx, "Close") > 0 {
				fmt.Println("[gfxAbout] Close button pressed!")
				ActiveNKWindow = "main"
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}
