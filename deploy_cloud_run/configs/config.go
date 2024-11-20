package configs

import (
	"errors"

	"github.com/spf13/viper"
)

type Conf struct {
	Port          int    `mapstructure:"PORT"`
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
}

var cfg *Conf

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("WEATHER_API_KEY")
	err := viper.ReadInConfig()
	viper.SetDefault("PORT", 8080)
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &configFileNotFoundError) {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	if cfg.WeatherApiKey == "" {
		return nil, errors.New("WEATHER_API_KEY is required")
	}
	return cfg, err
}

func GetConfig() Conf {
	return *cfg
}
