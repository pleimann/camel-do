package oauth

import "github.com/wailsapp/wails/v3/pkg/application"

type BrowserLauncher struct {
	window *application.WebviewWindow
}

func (l *BrowserLauncher) Launch(url string) error {
	// // create a window
	// window := application.Get().NewWebviewWindowWithOptions(
	// 	application.WebviewWindowOptions{
	// 		Title:  "OAuth Login",
	// 		Width:  600,
	// 		Height: 850,
	// 		URL:    url,
	// 	})
	// window.Show()

	// l.window = window

	

	return nil
}

func (l *BrowserLauncher) Shutdown() {
	if l.window != nil {
		l.window.Close()
	}
}
