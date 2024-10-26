package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func OAuth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  confString("REDIRECT_URL", ""),
		ClientID:     confString("CLIENT_ID", ""),
		ClientSecret: confString("CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// GetGoogleOAuthConfig returns the registered Google OAuth2 configuration
func GetGoogleOAuthConfig() *oauth2.Config {
	if googleOauthConfig == nil {
		panic("OAuth configuration not initialized")
	}
	return googleOauthConfig
}
