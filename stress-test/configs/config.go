package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Url         string `mapstructure:"url"`
	Concurrency int    `mapstructure:"concurrency"`
	Requests    int    `mapstructure:"requests"`
}

var cfg Config

func LoadConfig() {
	viper.AutomaticEnv()

	_ = viper.BindEnv("CONCURRENCY")
	_ = viper.BindEnv("REQUESTS")
	_ = viper.BindEnv("URL")

	_ = viper.Unmarshal(&cfg)
}

func GetConfig() *Config {
	return &cfg
}
