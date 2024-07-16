package config

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Schema   string `mapstructure:"schema"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"database"`
	Oauth2 struct {
		ClientId     string `mapstructure:"client-id"`
		ClientSecret string `mapstructure:"client-secret"`
		CallbackUrl  string `mapstructure:"callback-url"`
	} `mapstructure:"oauth2"`
	CookieStore struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"cookie-store"`
}
