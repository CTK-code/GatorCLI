package config

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBURL string `json:"db_url"`
}

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	fullPath := path.Join(homeDir, FileName)
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(fileData, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
