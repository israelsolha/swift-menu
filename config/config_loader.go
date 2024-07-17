package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"swift-menu-session/utils"

	"github.com/spf13/viper"
)

func LoadConfig(env string) Config {
	projectRoot := utils.FindProjectRoot()
	configPath := filepath.Join(projectRoot, "resources/config")

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if strings.HasPrefix(value, "env:") {
			envKey := strings.TrimPrefix(value, "env:")
			if envValue, exists := os.LookupEnv(envKey); exists {
				viper.Set(key, envValue)
			} else {
				panic(fmt.Errorf("required environment variable not set: %s", envKey))
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return config
}
