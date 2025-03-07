package main

import (
	"embed"
	"log"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"pleimann.com/camel-do/services/oauth"
	"pleimann.com/camel-do/services/task"
)

var assets embed.FS

func main() {
	viper.SetConfigFile(".camel-do")
	viper.ReadInConfig()

	oauthService := oauth.NewOauthService(oauth.Config{
		// SessionSecret is the secret used to encrypt the session store.
		SessionSecret: "secret",
		Providers: []goth.Provider{
			google.New(
				os.Getenv("clientid"),
				os.Getenv("secret"),
				"http://localhost:9876/auth/google/callback",
				"email", "profile"),
		},
	})

	app := application.New(application.Options{
		Name:        "camel-do",
		Description: "Camel Task Manager",
		Services: []application.Service{
			application.NewService(task.NewTaskService(task.Config{})),
			application.NewService(oauthService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.OnEvent(oauth.Success, func(e *application.CustomEvent) {
		app.Logger.Info("Got oauth success response: " + e.ToJSON())

		/*
			 {
				"name":"wails:oauth:success",
				"data":[{
					"RawData":{
						"email":"mike@pleimann.com",
						"family_name":"Pleimann",
						"given_name":"Mike",
						"hd":"pleimann.com",
						"id":"100504129432409136786",
						"name":"Mike Pleimann",
						"picture":"https://lh3.googleusercontent.com/a/ACg8ocLZINNAbtr1Xv2vn2HuXk2c2AIGVE3NpeLxWFeLuUrOqUdIFd_1=s96-c",
						"verified_email":true
					},
					"Provider":"google",
					"Email":"mike@pleimann.com",
					"Name":"Mike Pleimann",
					"FirstName":"Mike",
					"LastName":"Pleimann",
					"NickName":"Mike Pleimann",
					"Description":"",
					"UserID":"100504129432409136786",
					"AvatarURL":"https://lh3.googleusercontent.com/a/ACg8ocLZINNAbtr1Xv2vn2HuXk2c2AIGVE3NpeLxWFeLuUrOqUdIFd_1=s96-c",
					"Location":"",
					"AccessToken":"ya29.a0AeXRPp4zfJWtFMSnFOTjn47r1ukeESWbOa-Gv6nzcgPLHjl0UiWpOMtMOTyFQ0JnmurGeUIQkqtsne8lALQmpyrcCT1ZpcJ4cubL2jPhztDVjbSL9_Daf2zUVqIjkc4qBm8_YYOnKmnrLdr6z1FgOCyatKSKqHHblE5HOT7XaCgYKAR4SARISFQHGX2Mi6N6UAYTC2ezptnY937WNjw0175",
					"AccessTokenSecret":"",
					"RefreshToken":"1//05rGbgXsJ9YZxCgYIARAAGAUSNwF-L9Irf23O4Q7B6LdqILnW58_Jr2Xgl1ho-zx7GtqI82gWkFqC8NmVkyYD0VSyfWQxgucOUVA",
					"ExpiresAt":"2025-03-06T15:57:35.345068-06:00",
					"IDToken":"eyJhbGciOiJSUzI1NiIsImtpZCI6IjI1ZjgyMTE3MTM3ODhiNjE0NTQ3NGI1MDI5YjAxNDFiZDViM2RlOWMiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI5MDI2OTEzMTk0NzQtdnBzYjcwN2V1bmZrN2YzN24xcnBmcGtzY2d0bnNyOHAuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI5MDI2OTEzMTk0NzQtdnBzYjcwN2V1bmZrN2YzN24xcnBmcGtzY2d0bnNyOHAuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDA1MDQxMjk0MzI0MDkxMzY3ODYiLCJoZCI6InBsZWltYW5uLmNvbSIsImVtYWlsIjoibWlrZUBwbGVpbWFubi5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IlFMeE1EOGJVM25yZ1R5emJGZFlHRHciLCJuYW1lIjoiTWlrZSBQbGVpbWFubiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMWklOTkFidHIxWHYydm4ySHVYazJjMkFJR1ZFM05wZUx4V0ZlTHVVck9xVWRJRmRfMT1zOTYtYyIsImdpdmVuX25hbWUiOiJNaWtlIiwiZmFtaWx5X25hbWUiOiJQbGVpbWFubiIsImlhdCI6MTc0MTI5NDY1NiwiZXhwIjoxNzQxMjk4MjU2fQ.WZRgRPT3_Rb0dFYnBh58s51ZUzv-OkHIXk1IOsndA1ZdQAyL9ZQm_N2q0-hTprX1n-4oUw3zbdhmPdt5nNUhPZd3akfB_9rPsWua2Khc3uUklaUCRxCnay2snxrZOfbwdQ9ukpPQpI4XZZ2YixCZurzSRRrwkLWTIQ5SkDgHdfToCkrAA_uOChcGK0tPutCkbPLZ-nbl4lR3VuPlIFq8TjwLB7GOXnj23PBDQ6mv-csDcCm7pvmJAtqc85AFwx37LLLGi8GoY17BPWxKMpyiZl7P8ffSmkll9788jDTYsTo8TG9FcCAwVBYx4OvFmQcGbhJSPEW9c82Q55CxVZRLGA"
				}],
				"sender":"",
				"Cancelled":false
			}
		*/
	})

	app.OnEvent(oauth.Error, func(e *application.CustomEvent) {
		app.Logger.Info("Got oauth error response: " + e.ToJSON())
	})

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

func createMainWindow(app *application.App) *application.WebviewWindow {
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
