package main

import (
	"embed"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"pleimann.com/camel-do/services/task"
	"pleimann.com/camel-do/utils/google/oauth"
)

var assets embed.FS

func main() {
	godotenv.Load()

	viper.SetConfigFile(".camel-do")
	viper.ReadInConfig()

	oauthService := oauth.NewOauthService(&oauth.Config{})
	tasksService := task.NewTaskService(&task.Config{
		TokenSourceProvider: &oauthService.TokenSourceProvider,
	})

	app := application.New(application.Options{
		Name:        "camel-do",
		Description: "Camel Task Manager",
		Services: []application.Service{
			application.NewService(tasksService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	var appWindow application.Window
	user, err := oauthService.Authenticate(func(url string) {
		app.Logger.Info("Opening...", "url", url)
		window := app.NewWebviewWindowWithOptions(
			application.WebviewWindowOptions{
				Title:  "OAuth Login",
				Width:  600,
				Height: 850,
				URL:    url,
				Hidden: false,
			})

		appWindow = window.Show()
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	appWindow.Close()

	app.Logger.Info("Authenticated", "User: ", user)

	window := createMainWindow(app)

	app.OnApplicationEvent(events.Windows.SystemThemeChanged, func(event *application.ApplicationEvent) {
		app.Logger.Info("System theme changed!")
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")
		} else {
			app.Logger.Info("System is now using light mode!")
		}
	})

	window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window = createMainWindow(app)
		window.Show()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func createMainWindow(app *application.App) *application.WebviewWindow {
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "Camel Do",
		BackgroundType:  application.BackgroundTypeTranslucent,
		URL:             "/",
		Hidden:          false,
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
