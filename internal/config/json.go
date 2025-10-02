package config

import (
	"encoding/json"
	"os"
)

func IntoJson(cfg Config) ([]byte, error) {
	jsonData, err := json.Marshal(cfg)
	if err != nil {return nil, err}
	return jsonData, nil
}

func FromJson(url string) (Config, error) {
	data, err := os.ReadFile(url)
	if err != nil {return Config{}, err}

	var newConfig Config
	if err := json.Unmarshal(data, &newConfig); err != nil {
		return Config{}, err
	}

	return newConfig, nil
}