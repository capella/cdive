package controllers

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SiteName string `mapstructure:"name"`
	Server   struct {
		Address string
		DSN     *string `mapstructure:"mysql"`
		Secret  string
	}
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("CDIVE")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
