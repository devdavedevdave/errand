package config

import (
	"errors"
	"fmt"
	"io/fs"
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
		return nil, fmt.Errorf("load: %w", err)
	}

	return &config, nil
}

func Save(path string, config *Config) error {
	permission := 0o644
	file, err := toml.Marshal(config)

	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	errWriteFile := os.WriteFile(path, file, os.FileMode(permission))

	if errWriteFile != nil {
		return fmt.Errorf("save: %w", errWriteFile)
	}

	return nil
}

func SaveInitial(path string) error {
	config := Config{Version: 1}

	// TODO: Add syncing with existing configuration. There might be something existing,
	// therefore this need to be added

	err := Save(path, &config)

	if err != nil {
		return fmt.Errorf("save initial: %s, %q, %w", path, config, err)
	}

	return nil
}

func initConfigIfNeeded() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("initConfigIfNeeded: %w", err)
	}

	configExists, err := configExists(homeDir)

	if err != nil {
		return fmt.Errorf("init config if needed: %w", err)
	}

	if configExists {
		return nil
	}

	errNewConfig := createNewConfig(homeDir)
	if errNewConfig != nil {
		return fmt.Errorf("init config if needed: %w", errNewConfig)
	}

	return nil
}

func configExists(dir string) (bool, error) {
	path := filepath.Join(dir, configPath())
	_, err := os.Stat(path)

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("config exists: %s, %w", dir, err)
	}

	return true, nil
}

func createNewConfig(dir string) error {
	path := filepath.Join(dir, configPath())

	permission := os.FileMode(0o755)
	errDir := os.MkdirAll(filepath.Dir(path), permission)

	if errDir != nil {
		return fmt.Errorf("create new config: %w, %s", errDir, path)
	}

	errFile := SaveInitial(path)

	if errFile != nil {
		return fmt.Errorf("create new config: %w, %s", errFile, path)
	}

	return nil
}

func configPath() string {
	return ".config/errand/config.toml"
}
