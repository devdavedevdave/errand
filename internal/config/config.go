package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(path string) Config {
	config := Config{}
	configFile, err := os.ReadFile(path)

	if err != nil {

	}

	toml.Unmarshal(configFile, &config)
	return config
}

func Save(path string, config *Config) {
	file, err := toml.Marshal(config)

	if err != nil {

	}

	os.WriteFile(path, file, 0o644) // TODO: Find a better way to define this. Permission like this isn't ok.
}

func configPath() string {
path:
	string = "./config.toml"

	return path
}
