package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		GRPCPort int `mapstructure:"grpc_port"`
	} `mapstructure:"app"`
	Database struct {
		DBUser string `mapstructure:"mongo_user"`
		DBPassword string `mapstructure:"mongo_password"`
		DBHost string `mapstructure:"mongo_host"`
		DBName string `mapstructure:"mongo_db_name"`
		DBAppName string `mapstructure:"mongo_app_name"`
		DBRetryWrites bool `mapstructure:"mongo_retry_writes"`
		DBW string `mapstructure:"mongo_w"`
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
