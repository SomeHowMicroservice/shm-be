package config

import "github.com/spf13/viper"

type AppConfig struct {
	App struct {
		ServerHost  string `mapstructure:"server_host"`
		GRPCPort int64  `mapstructure:"grpc_port"`
	} `mapstructure:"app"`

	Services struct {
		AuthPort int64 `mapstructure:"auth_port"`
		UserPort int64 `mapstructure:"user_port"`
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
