package util

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstruture:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetDefault("DB_USERNAME", "")
	viper.SetDefault("DB_USERPASSWORD", "")
	viper.SetDefault("DB_DRIVER", "")
	viper.SetDefault("DB_SOURCE", "")
	viper.SetDefault("SERVER_ADDRESS", "")

	if _, err := os.Stat(path + "/app.env"); err == nil {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
