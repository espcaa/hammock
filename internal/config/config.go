package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Scale float32
}

func NewConfig() *Config {
	return &Config{
		Scale: 1.0,
	}
}

func (c *Config) Save() error {

	log.Println("Saving config:", c)

	// serialize as json & save to app dir
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configPath, err := configPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	configPath, err := configPath()
	if err != nil {
		return nil, err
	}

	jsonData, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return NewConfig(), nil // first run, no config file yet
	}
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, "hammock")
	if err := os.MkdirAll(appConfigDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}
