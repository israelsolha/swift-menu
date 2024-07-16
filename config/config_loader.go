package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(env string) (Config, error) {
	fmt.Println(env)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../resources")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if strings.HasPrefix(value, "env:") {
			envKey := strings.TrimPrefix(value, "env:")
			if envValue, exists := os.LookupEnv(envKey); exists {
				viper.Set(key, envValue)
			} else {
				return Config{}, fmt.Errorf("required environment variable not set: %s", envKey)
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
