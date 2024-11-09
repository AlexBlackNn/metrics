package configserver

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/caarlos0/env/v6"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	MetricTypeCounter = "counter"
	MetricTypeGauge   = "gauge"
	configJSON        = "config.json"
)

// Config consists project settings.
type Config struct {
	Env                   string `yaml:"env" env-default:"local" env:"ENV"`
	ServerAddr            string `yaml:"server_addr" env-default:":8080" env:"ADDRESS"`
	ServerReadTimeout     int    `yaml:"server_read_timeout" env-default:"100" env:"SEVER_READ_TIMEOUT" envDefault:"100"`
	ServerWriteTimeout    int    `yaml:"server_write_timeout" env-default:"100" env:"SEVER_WRITE_TIMEOUT" envDefault:"1000000000"`
	ServerIdleTimeout     int    `yaml:"server_idle_timeout" env-default:"100" env:"SEVER_IDLE_TIMEOUT" envDefault:"100"`
	ServerRequestTimeout  int    `yaml:"server_request_timeout" env-default:"300" env:"SEVER_REQUEST_TIMEOUT" envDefault:"300"`
	ServerStoreInterval   int    `yaml:"server_store_interval" env:"STORE_INTERVAL"`
	ServerFileStoragePath string `yaml:"server_file_storage_path" env-default:"/tmp/metrics-db.json" env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	ServerRestore         bool   `yaml:"server_restore" env-default:"true" env:"RESTORE" envDefault:"true"`
	ServerRateLimit       int    `yaml:"server_rate_limit" env-default:"10000" env:"SERVER_RATE_LIMIT" envDefault:"10000"`
	ServerDataBaseDSN     string `yaml:"server_data_base_dsn" env:"DATABASE_DSN"`
	ServerMigrationTable  string `yaml:"server_migration_table_name" env-default:"migrations" env:"SERVER_MIGRATION_TABLE_NAME" envDefault:"migrations"`
	HashKey               string `yaml:"hash_key" env:"KEY"`
	CryptoKeyPath         string `yaml:"env" env-default:"secret-key" env:"ENV"`
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

// fetchConfigPath fetches config path from command line flag or env var.
// Priority: env -> yml -> flag -> default.

// New loads config
func New() (*Config, error) {
	cfg := &Config{}
	var err error
	var configPath string

	err = loadConfigFromFile(cfg)
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.Env, "e", "local", "project environment")
	flag.StringVar(&cfg.ServerAddr, "a", ":8080", "host address")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.ServerStoreInterval, "i", 1, "metrics store interval")
	flag.StringVar(&cfg.ServerFileStoragePath, "f", "/tmp/metrics-db.json", "metrics store path")
	flag.BoolVar(&cfg.ServerRestore, "r", true, "restore saved metrics")
	flag.StringVar(&cfg.ServerDataBaseDSN, "d", "", "database dsn")
	flag.StringVar(&cfg.CryptoKeyPath, "crypto-key", "", "path to crypto file")

	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath != "" {
		cfg, err = LoadByPath(configPath)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadByPath loads config by path
func LoadByPath(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, config.ErrAbsentConfigFile
		}
		return nil, fmt.Errorf("LoadByPath stat error: %w", err)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, config.ErrReadConfigFailed
	}
	return &cfg, nil
}

func (c *Config) SaveToJson() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(configJSON, data, 0644)
}

func loadConfigFromFile(cfg *Config) error {
	data, err := os.ReadFile(configJSON)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}
