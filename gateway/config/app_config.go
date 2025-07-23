package config

import "github.com/spf13/viper"

type AppConfig struct {
	App struct {
		ServerHost string `mapstructure:"server_host"`
		GRPCPort   int    `mapstructure:"grpc_port"`
	} `mapstructure:"app"`

	Jwt struct {
		SecretKey   string `mapstructure:"secret_key"`
		AccessName  string `mapstructure:"access_name"`
		RefreshName string `mapstructure:"refresh_name"`
	} `mapstructure:"jwt"`

	Services struct {
		AuthPort int `mapstructure:"auth_port"`
		UserPort int `mapstructure:"user_port"`
	} `mapstructure:"services"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
