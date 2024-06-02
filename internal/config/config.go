package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
}

// fetchConfigPath fetches config path from command line flag or env var
// Priority: flag -> env -> default
// Default value is empty string
func fetchConfigPath() string {
	var res string
	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	if res == "" {
		res = "./cmd/server/config/demo.yaml"
	}

	return res
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return &Config{}, ErrEmptyConfigPath
	}
	cfg, err := LoadByPath(configPath)
	if err != nil {
		return &Config{}, err
	}
	return cfg, nil
}

func LoadByPath(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, ErrAbsentConfigFile
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &Config{}, ErrReadConfigFailed
	}
	return &cfg, nil
}
