package main

import (
	"embed"
	"log"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"pleimann.com/camel-do/services"
)

var assets embed.FS

func main() {
	viper.SetConfigFile(".camel-do")
	viper.ReadInConfig()

	app := application.New(application.Options{
		Name:        "camel-do",
		Description: "Camel Task Manager",
		Services: []application.Service{
			application.NewService(&services.GreetService{}),
			application.NewService(&services.TaskService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := createWindow(app)

	app.OnApplicationEvent(events.Windows.SystemThemeChanged, func(event *application.ApplicationEvent) {
		app.Logger.Info("System theme changed!")
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")
		} else {
			app.Logger.Info("System is now using light mode!")
		}
	})

	window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window = createWindow(app)
		window.Show()
	})

	app.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")

		} else {
			app.Logger.Info("System is now using light mode!")

		}
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func createWindow(app *application.App) *application.WebviewWindow {
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "Camel Do",
		BackgroundType:  application.BackgroundTypeTranslucent,
		URL:             "/",
		InitialPosition: application.WindowXY,
		X:               viper.GetInt("WINDOW_X"),
		Y:               viper.GetInt("WINDOW_Y"),
		Width:           768,
		Height:          1024,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 64,
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
	})

	return window
}
