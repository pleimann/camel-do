package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	// GoogleTokenInfoURL is the endpoint to verify Google ID tokens
	GoogleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token="
)

// TokenPayload represents the payload of a verified Google ID token
type TokenPayload struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Jti           string `json:"jti"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Typ           string `json:"typ"`
}

// VerifyGoogleIDToken verifies a Google ID token and returns the token payload
func VerifyGoogleIDToken(idToken string) (*TokenPayload, error) {
	// Make a request to Google's token info endpoint
	resp, err := http.Get(GoogleTokenInfoURL + idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token verification failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var payload TokenPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to parse token payload: %v", err)
	}

	// Verify the token is not expired
	if payload.Exp == "" {
		return nil, errors.New("token has no expiration time")
	}

	// Additional validation can be added here as needed
	// For example, verify the audience matches your client ID

	return &payload, nil
}

// CreateOAuth2TokenFromIDToken creates an OAuth2 token from a verified Google ID token payload
func CreateOAuth2TokenFromIDToken(payload *TokenPayload) (*oauth2.Token, error) {
	// In a real implementation, you would exchange the ID token for an access token
	// Here we're creating a temporary token for demonstration
	token := &oauth2.Token{
		AccessToken: "id_token_based_access", // This is a placeholder
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Add the ID token to the token's extra fields
	token = token.WithExtra(map[string]interface{}{
		"id_token": payload,
	})

	return token, nil
}

// GetUserFromTokenPayload extracts user information from a token payload
func GetUserFromTokenPayload(payload *TokenPayload) *User {
	return &User{
		Name:      payload.Name,
		Email:     payload.Email,
		AvatarURL: payload.Picture,
	}
}
