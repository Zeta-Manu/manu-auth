package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AuthService struct {
		HTTP struct {
			Port int `mapstructure:"port"`
		} `mapstructure:"http"`
		AWS struct {
			AccessKey       string `mapstructure:"access_key"`
			SecretAccessKey string `mapstructure:"secret_access_key"`
		} `mapstructure:"aws"`
		Cognito struct {
			Region     string `mapstructure:"region"`
			UserPoolId string `mapstructure:"user_pool_id"`
			ClientId   string `mapstructure:"client_id"`
		} `mapstructure:"cognito"`
		JWT struct {
			PublicKey string `mapstructure:"public_key"`
		} `mapstructure:"jwt"`
	} `mapstructure: "authService"`
}

func LoadConfig(filePath string) (*Config, error) {
	var config Config

	// Initialize Viper
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml") // Setting the file type to yaml

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config file into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}
