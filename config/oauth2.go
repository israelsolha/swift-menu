package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOauth2Config(config Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  config.Oauth2.CallbackUrl,
		ClientID:     config.Oauth2.ClientId,
		ClientSecret: config.Oauth2.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}
