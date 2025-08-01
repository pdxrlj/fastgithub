package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	host := NewHost()

	err := wails.Run(&options.App{
		Title: "githubfast",
		Windows: &windows.Options{
			WebviewIsTransparent: true,
		},
		Width:  1024,
		Height: 768,

		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        host.StartPoll,
		Bind: []interface{}{
			app,
			host,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
