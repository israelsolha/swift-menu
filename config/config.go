package config

type Config struct {
	Api         Api         `mapstructure:"api"`
	Database    Database    `mapstructure:"database"`
	Oauth2      Oauth2      `mapstructure:"oauth2"`
	CookieStore CookieStore `mapstructure:"cookie-store"`
	Docker      Docker      `mapstructure:"docker"`
}

type Api struct {
	Port int `mapstructure:"port"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
	Port     int    `mapstructure:"port"`
}

type Oauth2 struct {
	ClientID      string `mapstructure:"client-id"`
	ClientSecret  string `mapstructure:"client-secret"`
	CallbackURL   string `mapstructure:"callback-url"`
	AuthURL       string `mapstructure:"auth-url"`
	TokenURL      string `mapstructure:"token-url"`
	DeviceAuthURL string `mapstructure:"device-auth-url"`
	AuthStyle     int    `mapstructure:"auth-style"`
}

type CookieStore struct {
	Secret string `mapstructure:"secret"`
}

type Docker struct {
	TearDownTags []string `mapstructure:"tear-down-tags"`
}
