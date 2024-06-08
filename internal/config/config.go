package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Config consists project settings
type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	ServerAddr     string `yaml:"server_addr" env-default:":8080"`
	PollInterval   int    `yaml:"poll_interval" env-default:"2"`
	ReportInterval int    `yaml:"report_interval" env-default:"5"`
	ClientTimeout  int    `yaml:"client_timeout" env-default:"5"`
}

// fetchConfigPath fetches config path from command line flag or env var
// Priority: flag -> env -> default
// Default value is empty string

// Load loads config
func Load() (*Config, error) {
	cfg := &Config{}
	var err error
	var configPath string

	flag.StringVar(&cfg.Env, "e", "local", "project environment")
	flag.StringVar(&cfg.ServerAddr, "a", ":8080", "host address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "metrics report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&cfg.ClientTimeout, "t", 3, "agent request timeout")
	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath != "" {
		cfg, err = LoadByPath(configPath)
		if err != nil {
			return &Config{}, err
		}
		return cfg, nil
	}
	return cfg, nil
}

// LoadByPath loads config by path
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
