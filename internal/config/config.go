package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisSecret   string `mapstructure:"REDIS_SECRET"`

	GoogleOauth2ClientID     string `mapstructure:"GOOGLE_OAUTH2_CLIENT_ID"`
	GoogleOauth2ClientSecret string `mapstructure:"GOOGLE_OAUTH2_CLIENT_SECRET"`

	FacebookOauth2ClientID     string `mapstructure:"FACEBOOK_OAUTH2_CLIENT_ID"`
	FacebookOauth2ClientSecret string `mapstructure:"FACEBOOK_OAUTH2_CLIENT_SECRET"`

	TwitchOauth2ClientID     string `mapstructure:"TWITCH_OAUTH2_CLIENT_ID"`
	TwitchOauth2ClientSecret string `mapstructure:"TWITCH_OAUTH2_CLIENT_SECRET"`

	JwtSecret string `mapstructure:"JWT_SECRET_KEY"`

	CSRFSecret string `mapstructure:"CSRF_SECRET_KEY"`
	CSRFSecure bool   `mapstructure:"CSRF_SECURE"`

	UploadBaseUrl string `mapstructure:"UPLOAD_BASE_URL"`
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
