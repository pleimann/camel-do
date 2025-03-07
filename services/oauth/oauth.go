package oauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ---------------- OauthService Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

const (
	Success   = "wails:oauth:success"
	Error     = "wails:oauth:error"
	LoggedOut = "wails:oauth:loggedout"
)

type OauthService struct {
	config Config
	server *http.Server
	router *pat.Router
}

type Config struct {

	// Address to bind the temporary webserver to
	// Defaults to localhost:9876
	Address string

	// SessionSecret is the secret used to encrypt the session store.
	SessionSecret string

	// MaxAge is the maximum age of the session in seconds.
	MaxAge int

	// Providers is a list of goth providers to use.
	Providers []goth.Provider

	// WindowConfig is the configuration for the window that will be opened
	// to perform the OAuth login.
	WindowConfig *application.WebviewWindowOptions
}

func NewOauthService(config Config) *OauthService {
	result := &OauthService{
		config: config,
	}

	if result.config.MaxAge == 0 {
		result.config.MaxAge = 86400 * 30 // 30 days
	}

	if result.config.Address == "" {
		result.config.Address = "localhost:9876"
	}

	if result.config.WindowConfig == nil {
		result.config.WindowConfig = &application.WebviewWindowOptions{
			Title:  "OAuth Login",
			Width:  600,
			Height: 850,
			Hidden: true,
		}
	}
	return result
}

func (p *OauthService) Shutdown() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *OauthService) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	store := sessions.NewCookieStore([]byte(p.config.SessionSecret))
	store.MaxAge(p.config.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store
	goth.UseProviders(p.config.Providers...)

	return nil
}

func (p *OauthService) CallableByJS() []string {
	return []string{
		"Github",
		"LogoutGithub",
		"AzureAD",
		"LogoutAzureAD",
		"Google",
		"LogoutGoogle",
		"MicrosoftOnline",
		"LogoutMicrosoftOnline",
		"Okta",
		"LogoutOkta",
		"OpenIDConnect",
		"LogoutOpenIDConnect",
		"Slack",
		"LogoutSlack",
		"Zoom",
		"LogoutZoom",
	}
}

func (p *OauthService) InjectJS() string {
	return ""
}

func (p *OauthService) start(provider string) error {
	if p.server != nil {
		return fmt.Errorf("server already processing request. Please wait for the current login to complete")
	}

	router := pat.New()
	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			application.Get().EmitEvent(Error, err.Error())
		} else {
			application.Get().EmitEvent(Success, user)
		}

		_ = p.server.Close()
		p.server = nil
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: router,
	}

	go p.server.ListenAndServe()

	// Keep trying to connect until we succeed
	var keepTrying = true
	var connected = false

	go func() {
		time.Sleep(3 * time.Second)
		keepTrying = false
	}()

	for keepTrying {
		_, err := http.Get("http://" + p.config.Address)
		if err == nil {
			connected = true
			break
		}
	}

	if !connected {
		return fmt.Errorf("server failed to start")
	}

	// create a window
	p.config.WindowConfig.URL = "http://" + p.config.Address + "/auth/" + provider
	window := application.Get().NewWebviewWindowWithOptions(*p.config.WindowConfig)
	window.Show()

	application.Get().OnEvent(Success, func(event *application.CustomEvent) {
		window.Close()
	})
	application.Get().OnEvent(Error, func(event *application.CustomEvent) {
		window.Close()
	})

	return nil
}

func (p *OauthService) logout(provider string) error {
	if p.server != nil {
		return fmt.Errorf("server already processing request. Please wait for the current operation to complete")
	}

	router := pat.New()
	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		err := gothic.Logout(res, req)
		if err != nil {
			application.Get().EmitEvent(Error, err.Error())
		} else {
			application.Get().EmitEvent(LoggedOut)
		}
		_ = p.server.Close()
		p.server = nil
	})

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: router,
	}

	go p.server.ListenAndServe()

	// Keep trying to connect until we succeed
	var keepTrying = true
	var connected = false

	go func() {
		time.Sleep(3 * time.Second)
		keepTrying = false
	}()

	for keepTrying {
		_, err := http.Get("http://" + p.config.Address)
		if err == nil {
			connected = true
			break
		}
	}

	if !connected {
		return fmt.Errorf("server failed to start")
	}

	// create a window
	p.config.WindowConfig.URL = "http://" + p.config.Address + "/logout/" + provider
	window := application.Get().NewWebviewWindowWithOptions(*p.config.WindowConfig)
	window.Show()

	application.Get().OnEvent(LoggedOut, func(event *application.CustomEvent) {
		window.Close()
	})
	application.Get().OnEvent(Error, func(event *application.CustomEvent) {
		window.Close()
	})

	return nil
}

// ---------------- OauthService Methods ----------------

func (p *OauthService) Apple() error {
	return p.start("apple")
}

func (p *OauthService) Github() error {
	return p.start("github")
}

func (p *OauthService) Google() error {
	return p.start("google")
}

func (p *OauthService) MicrosoftOnline() error {
	return p.start("microsoftonline")
}

func (p *OauthService) Okta() error {
	return p.start("okta")
}

func (p *OauthService) Onedrive() error {
	return p.start("onedrive")
}

func (p *OauthService) OpenIDConnect() error {
	return p.start("openid-connect")
}

func (p *OauthService) Slack() error {
	return p.start("slack")
}

func (p *OauthService) Zoom() error {
	return p.start("zoom")
}

func (p *OauthService) LogoutApple() error {
	return p.logout("apple")
}

func (p *OauthService) LogoutGithub() error {
	return p.logout("github")
}

func (p *OauthService) LogoutGoogle() error {
	return p.logout("google")
}

func (p *OauthService) LogoutMicrosoftOnline() error {
	return p.logout("microsoftonline")
}

func (p *OauthService) LogoutOkta() error {
	return p.logout("okta")
}

func (p *OauthService) LogoutOnedrive() error {
	return p.logout("onedrive")
}

func (p *OauthService) LogoutOpenIDConnect() error {
	return p.logout("openid-connect")
}

func (p *OauthService) LogoutSlack() error {
	return p.logout("slack")
}

func (p *OauthService) LogoutZoom() error {
	return p.logout("zoom")
}
