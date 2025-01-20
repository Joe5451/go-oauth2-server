package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	GoogleOauth2ClientID     string `mapstructure:"GOOGLE_OAUTH2_CLIENT_ID"`
	GoogleOauth2ClientSecret string `mapstructure:"GOOGLE_OAUTH2_CLIENT_SECRET"`
}

var AppConfig Config

func InitializeAppConfig() error {
	// Automatically read environment variables
	viper.AutomaticEnv()

	// Optional: Read dotenv file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to load config file, %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}
