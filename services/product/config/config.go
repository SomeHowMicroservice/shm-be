package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		GRPCPort int `mapstructure:"grpc_port"`
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
