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
