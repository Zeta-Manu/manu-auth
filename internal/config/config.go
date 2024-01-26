package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type AWSConfig struct {
	AWSRegion         string `json:"region"`
	CognitoUserPoolID string `json:"cognitoUserPoolID"`
	CognitoClientID   string `json:"cognitoClientID"`
}

type HTTPConfig struct {
	Port int `json:"port"`
}

type AppConfig struct {
	AWS  AWSConfig  `json:"aws"`
	HTTP HTTPConfig `json:"http"`
}

func InitConfig(filePath string) (*AppConfig, error) {
	// Read JSON configuration file
	configData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON config file: %v", err)
	}

	// Unmarshall JSON into AppConfig struct
	var config AppConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON config: %v", err)
	}

	return &config, nil
}
