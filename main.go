package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"pleimann.com/camel-do/services"
)

var assets embed.FS

func main() {
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

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:          "Camel Do",
		BackgroundType: application.BackgroundTypeTranslucent,
		URL:            "/",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 64,
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
		},
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
