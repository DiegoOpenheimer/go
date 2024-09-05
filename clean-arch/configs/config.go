package configs

import (
	"errors"
	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	_ = viper.BindEnv("DB_DRIVER")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("WEB_SERVER_PORT")
	_ = viper.BindEnv("GRPC_SERVER_PORT")
	_ = viper.BindEnv("GRAPHQL_SERVER_PORT")
	_ = viper.BindEnv("RABBITMQ_URL")
	err := viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &configFileNotFoundError) {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
