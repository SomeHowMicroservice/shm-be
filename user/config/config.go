package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		GRPCPort int `mapstructure:"grpc_port"`
	} `mapstructure:"app"`

	Database struct {
		DBHost           string `mapstructure:"pg_host"`
		DBName           string `mapstructure:"pg_database"`
		DBUser           string `mapstructure:"pg_user"`
		DBPassword       string `mapstructure:"pg_password"`
		DBSSLMode        string `mapstructure:"pg_ssl_mode"`
		DBChannelBinding string `mapstructure:"pg_channel_binding"`
	} `mapstructure:"database"`
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
