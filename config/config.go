package config

import (
	"fmt"
	"os"

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

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// Initialize Viper
		viper.SetConfigFile(filePath)
		viper.SetConfigType("yaml") // Setting the file type to yaml

		// Read the config file
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	} else if os.IsNotExist(err) {
		// If the file does not exist, log the error and fall back to environment variables
		fmt.Printf("Config file does not exist: %v. Falling back to environment variables.\n", err)
	} else {
		// If there's another error (like permission denied), return the error
		return nil, fmt.Errorf("error checking config file: %v", err)
	}

	// Automatically read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// Bind specific environment variables to struct fields
	viper.BindEnv("AuthService.HTTP.Port", "APP_HTTP_PORT")
	viper.BindEnv("AuthService.AWS.AccessKey", "APP_AWS_ACCESS_KEY")
	viper.BindEnv("AuthService.AWS.SecretAccessKey", "APP_AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("AuthService.Cognito.Region", "APP_COGNITO_REGION")
	viper.BindEnv("AuthService.Cognito.UserPoolId", "APP_COGNITO_USER_POOL_ID")
	viper.BindEnv("AuthService.Cognito.ClientId", "APP_COGNITO_CLIENT_ID")
	viper.BindEnv("AuthService.JWT.PublicKey", "APP_JWT_PUBLIC_KEY")

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}
