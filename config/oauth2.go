package config

import (
	"golang.org/x/oauth2"
)

func NewOauth2Config(oauth2Config Oauth2) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  oauth2Config.CallbackURL,
		ClientID:     oauth2Config.ClientID,
		ClientSecret: oauth2Config.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Config.AuthURL,
			DeviceAuthURL: oauth2Config.DeviceAuthURL,
			TokenURL:      oauth2Config.TokenURL,
			AuthStyle:     oauth2.AuthStyle(oauth2Config.AuthStyle),
		},
	}
}
