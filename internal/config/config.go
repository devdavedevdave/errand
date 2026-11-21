package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func Load(path string) (*Config, error) {
	config := Config{}
	configFile, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	err = toml.Unmarshal(configFile, &config)

	if err != nil {
		return nil, fmt.Errorf("load: %q", err)
	}

	return &config, nil
}

func Save(path string, config *Config) error {
	permission := 0o644
	file, err := toml.Marshal(config)

	if err != nil {
		return fmt.Errorf("save: ", err)
	}

	os.WriteFile(path, file, os.FileMode(permission))
	return nil
}

func initConfigIfNeeded() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("initConfigIfNeeded: %q", err)
	}

	configExists, err := configExists(homeDir)

	if err != nil {
		fmt.Errorf("init config if needed: %q", err)
	}

	if configExists {
		return nil
	}

	errNewConfig := createNewConfig(homeDir)
	if errNewConfig != nil {
		return fmt.Errorf("init config if needed: %q", errNewConfig)
	}

	return nil
}

func configExists(dir string) (bool, error) {
	path := filepath.Join(dir, configPath())
	fileInfo, err := os.Stat(path)

	if errors.Is() {
		return false, fmt.Errorf("config exists: %q", err)
	}

	return fileInfo.Size() > 0, nil
}

func createNewConfig(dir string) error {
	path := filepath.Join(dir, configPath())

	_, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("create new config: %q, %path", err, path)
	}

	return nil
}

func configPath() string {
	return ".config/errand/config.toml"
}
