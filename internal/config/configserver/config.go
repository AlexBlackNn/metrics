package configserver

import (
	"flag"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/caarlos0/env/v6"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Config consists project settings
type Config struct {
	Env                   string `yaml:"env" env-default:"local" env:"ENV"`
	ServerAddr            string `yaml:"server_addr" env-default:":8080" env:"ADDRESS"`
	ServerReadTimeout     int    `yaml:"server_read_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerWriteTimeout    int    `yaml:"server_write_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerIdleTimeout     int    `yaml:"server_idle_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerStoreInterval   int    `yaml:"server_store_interval" env:"STORE_INTERVAL"`
	ServerFileStoragePath string `yaml:"server_file_storage_path" env-default:"/tmp/metrics-db.json" env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	ServerRestore         bool   `yaml:"server_restore" env-default:"true" env:"RESTORE" envDefault:"true"`
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Env: %s,"+
			" ServerAddr: %s, "+
			"ServerReadTimeout: %d,"+
			"ServerWriteTimeout: %d,"+
			"ServerIdleTimeout: %d,"+
			"ServerStoreInterval: %d,",
		c.Env,
		c.ServerAddr,
		c.ServerReadTimeout,
		c.ServerWriteTimeout,
		c.ServerIdleTimeout,
		c.ServerStoreInterval,
	)
}

// fetchConfigPath fetches config path from command line flag or env var
// Priority: env -> yml -> flag -> default

// New loads config
func New() (*Config, error) {
	cfg := &Config{}
	var err error
	var configPath string

	flag.StringVar(&cfg.Env, "e", "local", "project environment")
	flag.StringVar(&cfg.ServerAddr, "a", ":8080", "host address")
	flag.IntVar(&cfg.ServerStoreInterval, "i", 1, "metrics store interval")
	flag.StringVar(&cfg.ServerFileStoragePath, "f", "/tmp/metrics-db.json", "metrics store path")
	flag.BoolVar(&cfg.ServerRestore, "r", true, "restore saved metrics")

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
		return cfg, err
	}
	return cfg, nil
}

// LoadByPath loads config by path
func LoadByPath(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, config.ErrAbsentConfigFile
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &Config{}, config.ErrReadConfigFailed
	}
	return &cfg, nil
}
