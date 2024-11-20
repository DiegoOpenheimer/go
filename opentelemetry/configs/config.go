package configs

import (
	"errors"

	"github.com/spf13/viper"
)

type Conf struct {
	ServiceAPort             int    `mapstructure:"SERVICE_A_PORT"`
	ServiceBPort             int    `mapstructure:"SERVICE_B_PORT"`
	ServiceBUrl              string `mapstructure:"SERVICE_B_URL"`
	WeatherApiKey            string `mapstructure:"WEATHER_API_KEY"`
	ServiceName              string `mapstructure:"SERVICE_NAME"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

var cfg *Conf

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	_ = viper.BindEnv("SERVICE_A_PORT")
	_ = viper.BindEnv("SERVICE_B_PORT")
	_ = viper.BindEnv("SERVICE_B_URL")
	_ = viper.BindEnv("WEATHER_API_KEY")
	_ = viper.BindEnv("SERVICE_NAME")
	_ = viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	err := viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &configFileNotFoundError) {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	if cfg.WeatherApiKey == "" && cfg.ServiceBPort != 0 {
		return nil, errors.New("WEATHER_API_KEY is required")
	}
	return cfg, err
}

func GetConfig() Conf {
	return *cfg
}
