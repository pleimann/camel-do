package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
)

type GoogleAuth struct {
	config *oauth2.Config
}

func NewGoogleAuth(credentials string) *GoogleAuth {
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(credentials), tasks.TasksReadonlyScope, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return &GoogleAuth{
		config: config,
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func (a *GoogleAuth) GetClient() *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Unable to get user cache directory: %v", err)
	}

	tokFile := path.Join(userCacheDir, "camel-do", "token.json")

	tok, err := a.tokenFromFile(tokFile)
	if err != nil {
		slog.Info("No existing token found, getting token from web...")
		tok = a.getTokenFromWeb(a.config)
		a.saveTokenToFile(tokFile, tok)
	} else {
		slog.Info("Got token", "file", tokFile)

		// Check if token is expired and refresh token is not available
		if !tok.Valid() && tok.RefreshToken == "" {
			slog.Info("Token expired and no refresh token available, getting new token from web...")
			tok = a.getTokenFromWeb(a.config)
			a.saveTokenToFile(tokFile, tok)
		}
	}

	return a.config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func (a *GoogleAuth) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// TODO implement PKCE
	authURL := config.AuthCodeURL("state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce, // Force approval prompt to ensure refresh token
	)

	var authCodeChannel = make(chan string)
	defer close(authCodeChannel)

	srv := a.startServer(authCodeChannel)

	a.launchBrowser(authURL)

	authCode := <-authCodeChannel

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	// Log whether we got a refresh token
	if tok.RefreshToken != "" {
		slog.Info("Successfully obtained refresh token")
	} else {
		slog.Warn("No refresh token received - user may need to revoke and re-authorize")
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "error", err)
	}

	return tok
}

func (a *GoogleAuth) startServer(authCodeChannel chan string) *http.Server {
	srv := &http.Server{
		Addr: ":9876",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Query().Has("code") {
				authCodeChannel <- req.URL.Query().Get("code")
				w.Header().Set("Location", "http://localhost:4000")
				w.WriteHeader(http.StatusTemporaryRedirect)
				fmt.Fprintf(w, "<html><body>Authorization successful. You can close this window now.</body></html>")
			}
		}),
	}

	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	slog.Debug("Listening on http://localhost:9876...")

	// returning reference so caller can call Shutdown()
	return srv
}

func (a *GoogleAuth) launchBrowser(url string) {
	var err error

	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()

	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()

	case "darwin":
		err = exec.Command("open", url).Start()

	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Retrieves a token from a local file.
func (a *GoogleAuth) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (a *GoogleAuth) saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Helper method to save token with directory creation and logging
func (a *GoogleAuth) saveTokenToFile(tokFile string, tok *oauth2.Token) {
	slog.Info("Saving token...", "file", tokFile)

	if err := os.MkdirAll(path.Dir(tokFile), 0700); err != nil {
		log.Fatalf("Unable to create directory for token file: %v", err)
	}

	a.saveToken(tokFile, tok)
	slog.Info("Saved token", "file", tokFile)
}
