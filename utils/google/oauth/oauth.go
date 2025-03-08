package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/keyring"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type User struct {
	Name      string
	Email     string
	AvatarURL string
}

type Config struct {
}

type TokenSourceProvider func(context.Context) (oauth2.TokenSource, error)

type OauthService struct {
	TokenSourceProvider

	config      *Config
	user        *User
	oauthConfig *oauth2.Config
	token       *oauth2.Token
	ring        keyring.Keyring
}

var clientId, clientSecret string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var exists bool

	clientId, exists = os.LookupEnv("GOOGLE_CLIENT_ID")
	if !exists {
		log.Fatal("GOOGLE_CLIENT_ID is not set")
	}

	clientSecret, exists = os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !exists {
		log.Fatal("GOOGLE_CLIENT_SECRET is not set")
	}
}

func NewOauthService(config *Config) *OauthService {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "camel-do",
	})

	if err != nil {
		log.Fatal(err)
	}

	service := &OauthService{
		config: config,
		oauthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:9876/auth/google/callback",
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/tasks",
				"https://www.googleapis.com/auth/calendar",
				"https://mail.google.com/",
			},
		},
		ring: ring,
	}

	return service
}

func (s *OauthService) TokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	if s.token == nil {
		tokenJson, err := s.ring.Get("googleToken")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(tokenJson.Data, &s.token)

		if err != nil {
			return nil, err
		}
	}

	return s.oauthConfig.TokenSource(ctx, s.token), nil
}

func (s *OauthService) User() *User {
	return s.user
}

func (s *OauthService) Authenticate(launcher func(url string)) (*User, error) {
	callbackHandler := NewGoogleOauthCallbackHandler(s.oauthConfig)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello World!")) })
	mux.HandleFunc("/auth/google/login", s.oauthGoogleLogin)
	mux.Handle("/auth/google/callback", callbackHandler)

	var r http.Handler = mux
	server := &http.Server{
		Addr:    ":9876",
		Handler: r,
	}
	defer server.Close()

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	go server.ListenAndServe()

	if err := waitForServerAvail("http://localhost:9876/"); err != nil {
		return nil, err
	}

	launcher("http://localhost:9876/auth/google/login")

	user, token := callbackHandler.GetUserAndToken()

	fmt.Println("got user", "user", user)
	fmt.Println("storing token in keyring", "token", token)
	s.storeToken(token)

	s.user = user

	return user, nil
}

func (s *OauthService) storeToken(token *oauth2.Token) {
	tokenJson, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}

	_ = s.ring.Set(keyring.Item{
		Key:   "googleToken",
		Data:  []byte(tokenJson),
		Label: "Google OAuth2 Token",
	})
}

func (s *OauthService) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: oauthState, Expires: expiration}
	http.SetCookie(w, &cookie)

	u := s.oauthConfig.AuthCodeURL(oauthState)

	fmt.Printf("redirecting to %s\n", u)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func waitForServerAvail(url string) error {
	// Keep trying to connect until we succeed
	var keepTrying = true
	var connected = false

	go func() {
		time.Sleep(3 * time.Second)
		keepTrying = false
	}()

	for keepTrying {
		_, err := http.Get(url)
		if err == nil {
			connected = true
			break
		}
	}

	if !connected {
		return fmt.Errorf("server failed to start")
	}

	return nil
}
