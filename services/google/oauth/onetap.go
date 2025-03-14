package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// OneTapHandler handles Google One Tap authentication
type OneTapHandler struct {
	oauthConfig *oauth2.Config
	data        chan (*UserInfo)
}

// NewOneTapHandler creates a new Google One Tap handler
func NewOneTapHandler(oauthConfig *oauth2.Config) *OneTapHandler {
	return &OneTapHandler{
		data:        make(chan *UserInfo),
		oauthConfig: oauthConfig,
	}
}

// GetUserAndToken waits for and returns the authenticated user and token
func (h *OneTapHandler) GetUserAndToken() (*User, *oauth2.Token) {
	userInfo := <-h.data
	fmt.Printf("received user from one tap: %s\n", userInfo.String())
	return userInfo.User, userInfo.Token
}

// ServeHTTP handles the One Tap callback
func (h *OneTapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the credential from the request
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	credential := r.FormValue("credential")
	if credential == "" {
		http.Error(w, "No credential provided", http.StatusBadRequest)
		return
	}

	// Verify the ID token with Google
	payload, err := VerifyGoogleIDToken(credential)
	if err != nil {
		log.Printf("Error verifying ID token: %v", err)
		http.Error(w, "Failed to verify ID token", http.StatusUnauthorized)
		return
	}

	// Verify that the audience matches our client ID
	if payload.Aud != h.oauthConfig.ClientID {
		log.Printf("Token audience mismatch: %s != %s", payload.Aud, h.oauthConfig.ClientID)
		http.Error(w, "Token audience mismatch", http.StatusUnauthorized)
		return
	}

	// Create an OAuth2 token from the ID token
	token, err := CreateOAuth2TokenFromIDToken(payload)
	if err != nil {
		log.Printf("Error creating OAuth2 token: %v", err)
		http.Error(w, "Failed to create OAuth2 token", http.StatusInternalServerError)
		return
	}

	// Extract user information from the token payload
	user := GetUserFromTokenPayload(payload)

	// Create user info
	userInfo := &UserInfo{
		Token: token,
		User:  user,
	}

	// Send the user info to the channel
	h.data <- userInfo
	defer close(h.data)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// verifyGoogleCredential verifies the ID token with Google and returns user info
func (h *OneTapHandler) verifyGoogleCredential(credential string) (*UserInfo, error) {
	// Create a token from the credential
	token := &oauth2.Token{
		AccessToken: credential,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour), // Temporary expiry
	}

	// Get user info from Google
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

	var userInfoMap map[string]interface{}
	if err := json.Unmarshal(body, &userInfoMap); err != nil {
		return nil, err
	}

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
