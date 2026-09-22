package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	AllowedOrigins       []string      `mapstructure:"ALLOWED_ORIGINS"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	for _, key := range []string{
		"ENVIRONMENT",
		"ALLOWED_ORIGINS",
		"DB_SOURCE",
		"MIGRATION_URL",
		"REDIS_ADDRESS",
		"HTTP_SERVER_ADDRESS",
		"GRPC_SERVER_ADDRESS",
		"TOKEN_SYMMETRIC_KEY",
		"ACCESS_TOKEN_DURATION",
		"REFRESH_TOKEN_DURATION",
		"EMAIL_SENDER_NAME",
		"EMAIL_SENDER_ADDRESS",
		"EMAIL_SENDER_PASSWORD",
	} {
		if err = viper.BindEnv(key); err != nil {
			return
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}