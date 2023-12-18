package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey  string        `mapstructure:"SYMMETRIC_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetDefault("DB_USERNAME", "")
	viper.SetDefault("DB_USERPASSWORD", "")
	viper.SetDefault("DB_DRIVER", "")
	viper.SetDefault("DB_SOURCE", "")
	viper.SetDefault("SERVER_ADDRESS", "")
	viper.SetDefault("TOKEN_DURATION", "")
	viper.SetDefault("SYMMETRIC_KEY", "")

	if _, err := os.Stat(path + "/app.env"); err == nil {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		errRc := viper.ReadInConfig()
		if errRc != nil {
			return config, errRc
		}
	}
	if err != nil {
		return
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
