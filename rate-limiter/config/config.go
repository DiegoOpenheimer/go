package config

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	RedisHost      string        `mapstructure:"REDIS_HOST"`
	RedisPassword  string        `mapstructure:"REDIS_PASSWORD"`
	RateLimitIP    int           `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken int           `mapstructure:"RATE_LIMIT_TOKEN"`
	BlockedTime    time.Duration `mapstructure:"BLOCKED_TIME"`
}

var cfg Config

func LoadConfig(path string) *Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &configFileNotFoundError) {
		panic(err)
	}
	viper.AutomaticEnv()

	_ = viper.BindEnv("REDIS_HOST")
	_ = viper.BindEnv("REDIS_PASSWORD")
	_ = viper.BindEnv("RATE_LIMIT_IP")
	_ = viper.BindEnv("RATE_LIMIT_TOKEN")
	_ = viper.BindEnv("BLOCKED_TIME")

	_ = viper.Unmarshal(&cfg)
	return &cfg
}

func GetConfig() *Config {
	return &cfg
}
