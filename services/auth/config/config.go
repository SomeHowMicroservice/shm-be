package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		ServerHost string `mapstructure:"server_host"`
		GRPCPort int64 `mapstructure:"grpc_port"`
	} `mapstructure:"app"`

	Services struct {
		UserPort int64 `mapstructure:"user_port"`
	} `mapstructure:"services"`

	Cache struct {
		CHost     string `mapstructure:"rd_host"`
		CPort     string `mapstructure:"rd_port"`
		CPassword string `mapstructure:"rd_password"`
	} `mapstructure:"cache"`
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
