package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Config consists project settings
type Config struct {
	Env            string `yaml:"env" env-default:"local" env:"ENV"`
	ServerAddr     string `yaml:"server_addr" env-default:":8080" env:"ADDRESS"`
	PollInterval   int    `yaml:"poll_interval" env-default:"2" env:"POLL_INTERVAL"`
	ReportInterval int    `yaml:"report_interval" env-default:"5" env:"REPORT_INTERVAL"`
	ClientTimeout  int    `yaml:"client_timeout" env-default:"5" env:"CLIENT_TIMEOUT"`
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

	err = env.Parse(cfg)
	if err != nil {
		fmt.Println("1111111111111111", err)
		return cfg, err
	}
	fmt.Println("2222222222222", cfg)

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
