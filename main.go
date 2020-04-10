package main

import (
	"fmt"
	"runtime"

	assets "thdwb/assets"
	bun "thdwb/bun"
	gg "thdwb/gg"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	structs "thdwb/structs"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var perf *profiler.Profiler

func main() {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	mustard.SetGLFWHints()

	perf = profiler.CreateProfiler()

	browser := &structs.WebBrowser{
		Document: loadDocumentFromAsset(assets.HomePage()),
	}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton, urlInput := createMainBar(window, browser)
	urlInput.SetReturnCallback(func() {
		fmt.Println("enter")
		goButton.Click()
	})

	debugFrame := createDebugFrame(window, browser)
	rootFrame.AttachWidget(appBar)

	viewPort := mustard.CreateContextWidget(func(ctx *gg.Context) {
		perf.Start("parse")
		parsedDoc := ketchup.ParseDocument(browser.Document.RawDocument)
		perf.Stop("parse")

		perf.Start("render")
		bun.RenderDocument(ctx, parsedDoc)
		perf.Stop("render")

		statusLabel.SetContent("Loaded; " +
			"Render: " + perf.GetProfile("render").GetElapsedTime().String() + "; " +
			"Parsing: " + perf.GetProfile("parse").GetElapsedTime().String() + "; ")
	})

	//viewPort.EnableScrolling()
	window.RegisterButton(menuButton, func() {
		if debugFrame.GetHeight() != 300 {
			debugFrame.SetHeight(300)
		} else {
			debugFrame.SetHeight(0)
		}
	})

	window.RegisterButton(goButton, func() {
		if urlInput.GetValue() != browser.Document.URL {
			statusLabel.SetContent("Loading: " + urlInput.GetValue())
			go func() {
				loadDocument(browser, urlInput.GetValue(), func() {})
				ctx := viewPort.GetContext()
				ctx.SetRGB(1, 1, 1)
				ctx.Clear()

				perf.Start("parse")
				parsedDoc := ketchup.ParseDocument(browser.Document.RawDocument)
				perf.Stop("parse")

				perf.Start("render")
				bun.RenderDocument(ctx, parsedDoc)
				perf.Stop("render")

				statusLabel.SetContent(
					"Loaded; " +
						"Render: " + perf.GetProfile("render").GetElapsedTime().String() + "; " +
						"Parsing: " + perf.GetProfile("parse").GetElapsedTime().String() + "; ",
				)

				viewPort.RequestRepaint()
				statusLabel.RequestRepaint()
			}()
		}
	})

	rootFrame.AttachWidget(viewPort)
	rootFrame.AttachWidget(debugFrame)

	window.SetRootFrame(rootFrame)
	window.Show()

	app.AddWindow(window)
	app.Run(func() {})
}
