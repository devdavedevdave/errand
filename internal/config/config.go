package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	filePerm os.FileMode = 0o644 // config file: owner read/write, everyone else read
	dirPerm  os.FileMode = 0o755 // config dir: owner full access, everyone else read/enter
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
	file, err := toml.Marshal(config)

	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	err = os.WriteFile(path, file, filePerm)

	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func SaveInitial(path string) error {
	config := Config{Version: 1}

	err := Save(path, &config)

	if err != nil {
		return fmt.Errorf("save initial: %w", err)
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

	err = createNewConfig(homeDir)
	if err != nil {
		return fmt.Errorf("init config if needed: %w", err)
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

	err := os.MkdirAll(filepath.Dir(path), dirPerm)

	if err != nil {
		return fmt.Errorf("create new config: %w, %s", err, path)
	}

	err = SaveInitial(path)

	if err != nil {
		return fmt.Errorf("create new config: %w, %s", err, path)
	}

	return nil
}

func configPath() string {
	return ".config/errand/config.toml"
}
