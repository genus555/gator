package config

import (
	//"fmt"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL				string `json:"db_url"`
	CurrentUserName		string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	/*//for testing
	home_dir, err := os.UserHomeDir()
	if err != nil {return "", err}
	full_path := home_dir + configFileName*/

	//actual
	work_dir, err := os.Getwd()
	if err != nil {return "", err}
	full_path := work_dir + configFileName

	return full_path, nil
}

func writeToConfig(cfg Config) error {
	full_path, err := getConfigFilePath()
	if err != nil {return err}

	jsonData, err := IntoJson(cfg)
	if err != nil {return err}

	err = os.WriteFile(full_path, jsonData, 0644)
	if err != nil {return err}

	return nil
}

func ReadJson() (Config, error) {
	full_path, err := getConfigFilePath()
	if err != nil {return Config{}, err}

	newConfig, err := FromJson(full_path)
	if err != nil {return Config{}, err}

	return newConfig, nil
}

func SetUser(cfg Config, username string) error {
	cfg.CurrentUserName = username
	err := writeToConfig(cfg)
	if err != nil {return err}
	return nil
}
