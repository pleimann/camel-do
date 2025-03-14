package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	*User
	Token *oauth2.Token
}

func (ui *UserInfo) String() string {
	return fmt.Sprintf("Name: %s, Email: %s, Avatar URL: %s", ui.Name, ui.Email, ui.AvatarURL)
}

const googleUserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo"

func NewGoogleOauthCallbackHandler(oauthConfig *oauth2.Config) *GoogleOauthCallbackHandler {
	return &GoogleOauthCallbackHandler{
		data:        make(chan *UserInfo),
		oauthConfig: oauthConfig,
	}
}

type GoogleOauthCallbackHandler struct {
	data        chan (*UserInfo)
	oauthConfig *oauth2.Config
}

func (h *GoogleOauthCallbackHandler) GetUserAndToken() (*User, *oauth2.Token) {
	userInfo := <-h.data

	fmt.Printf("received user from channel %s\n", userInfo.String())

	return userInfo.User, userInfo.Token
}

func (h *GoogleOauthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	user, err := h.getUserDataFromGoogle(code)

	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("sending user to channel")
	h.data <- user
	defer close(h.data)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *GoogleOauthCallbackHandler) getUserDataFromGoogle(code string) (*UserInfo, error) {
	// Use code to get token and get user info from Google.

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	client := h.oauthConfig.Client(context.Background(), token)

	response, err := client.Get(googleUserInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	var userInfoMap map[string](interface{})
	if err := json.Unmarshal(body, &userInfoMap); err != nil {
		return nil, err
	}

	fmt.Printf("userInfoMap: %+v\n", userInfoMap)

	userInfo := &UserInfo{
		Token: token,
		User: &User{
			Name:      userInfoMap["name"].(string),
			Email:     userInfoMap["email"].(string),
			AvatarURL: userInfoMap["picture"].(string),
		},
	}

	return userInfo, nil
}
