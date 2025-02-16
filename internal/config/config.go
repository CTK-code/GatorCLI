package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Get the full path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get path: %e", err)
	}
	fullPath := path.Join(homeDir, FileName)
	return fullPath, nil
}

func Read() (Config, error) {
	fullPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	// Read from file
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config: %e", err)
	}

	var config Config
	if err := json.Unmarshal(fileData, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal data: %e", err)
	}

	return config, nil
}

func SetUser(userName string) error {
	confData, err := Read()
	if err != nil {
		return err
	}
	confData.CurrentUserName = userName
	bytes, err := json.Marshal(confData)
	if err != nil {
		return err
	}
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}
