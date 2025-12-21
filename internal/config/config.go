package config

import (
	"encoding/json"
	"os"
)

const gatorConfigFileName string = ".gatorconfig.json"

func SetUser(username string) error {
	conf, err := Read()
	if err != nil {
		return err
	}
	conf.CurrentUserName = username
	if err := write(conf); err != nil {
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
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, os.ModeDevice); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path += gatorConfigFileName
	return path, nil
}

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}
