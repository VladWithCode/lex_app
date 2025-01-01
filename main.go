package main

import (
	"context"
	"embed"
	"log"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/controllers"
	_db "github.com/vladwithcode/lex_app/internal/db"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Connect to DB
	db, err := _db.Connect()
	if err != nil {
		log.Fatalf("Couldn't connect to DB: %v\n  ", err)
	}
	defer db.Close()

	// Create an instance of the app structure
	app := NewApp()
	caseCtl := controllers.NewCaseControler()
	accUpdtrCtl := controllers.NewAccordUpdaterCtl()

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "lex_app",
		MinWidth:  1024,
		MinHeight: 768,
		Width:     1600,
		Height:    900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx, db)
			caseCtl.Startup(ctx, db)
			accUpdtrCtl.Startup(ctx, db)
		},
		Bind: []interface{}{
			app,
			caseCtl,
			accUpdtrCtl,
		},
		EnumBind: []interface{}{
			internal.AllRegions,
			internal.AllCaseTypes,
		},
		Frameless: true,
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
		Windows: &windows.Options{
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
