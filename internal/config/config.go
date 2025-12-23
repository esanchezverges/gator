package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const gatorConfigFileName string = ".gatorconfig.json"

func (c *Config) SetUser(username string) error {
	c.Currentusername = username
	if err := write(c); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading file in Read: %v", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("Error unmarshaling data in Read: %v", err)
	}
	return config, nil
}

func write(cfg *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("Error marshalingindent in write: %v", err)
	}
	if err := os.WriteFile(path, data, os.ModeDevice); err != nil {
		return fmt.Errorf("Error writing data in write: %v", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error in getConfigFilePath: %v", err)
	}
	path += "/" + gatorConfigFileName
	return path, nil
}

type Config struct {
	Dburl           string `json:"db_url"`
	Currentusername string `json:"current_user_name"`
}
