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
	Env                   string `yaml:"env" env-default:"local" env:"ENV"`
	ServerAddr            string `yaml:"server_addr" env-default:":8080" env:"ADDRESS"`
	PollInterval          int    `yaml:"poll_interval" env-default:"2" env:"POLL_INTERVAL"`
	ReportInterval        int    `yaml:"report_interval" env-default:"5" env:"REPORT_INTERVAL"`
	ClientTimeout         int    `yaml:"client_timeout" env:"CLIENT_TIMEOUT"`
	ServerReadTimeout     int    `yaml:"server_read_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerWriteTimeout    int    `yaml:"server_write_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerIdleTimeout     int    `yaml:"server_idle_timeout" env-default:"10" env:"SEVER_READ_TIMEOUT" envDefault:"10"`
	ServerStoreInterval   int    `yaml:"server_store_interval" env-default:"300" env:"STORE_INTERVAL" envDefault:"300"`
	ServerFileStoragePath string `yaml:"server_file_storage_path" env-default:"/tmp/metrics-db.json" env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	ServerRestore         bool   `yaml:"server_restore" env-default:"true" env:"RESTORE" envDefault:"true"`
	AgentRetryCount       int    `yaml:"agent_retry_count" env-default:"3" env:"AGENT_RETRY_COUNT" envDefault:"3"`
	AgentRetryWaitTime    int    `yaml:"agent_retry_wait_time" env-default:"30" env:"AGENT_RETRY_WAIT_TIME" envDefault:"30"`
	AgentRetryMaxWaitTime int    `yaml:"agent_retry_max_wait_time" env-default:"90" env:"AGENT_RETRY_MAX_WAIT_TIME" envDefault:"90"`
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Env: %s,"+
			" ServerAddr: %s, "+
			"PollInterval: %d,"+
			" ReportInterval: %d,"+
			" ClientTimeout: %d,"+
			"ServerReadTimeout: %d,"+
			"ServerWriteTimeout: %d,"+
			"ServerIdleTimeout: %d,",
		c.Env,
		c.ServerAddr,
		c.PollInterval,
		c.ReportInterval,
		c.ClientTimeout,
		c.ServerReadTimeout,
		c.ServerWriteTimeout,
		c.ServerIdleTimeout,
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
	flag.IntVar(&cfg.ReportInterval, "r", 10, "metrics report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&cfg.ClientTimeout, "t", 100, "agent request timeout")
	flag.IntVar(&cfg.ServerStoreInterval, "i", 300, "metrics store interval")
	flag.StringVar(&cfg.ServerFileStoragePath, "f", "/tmp/metrics-db.json", "metrics store path")
	flag.BoolVar(&cfg.ServerRestore, "re", true, "restore saved metrics")

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
		return &Config{}, ErrAbsentConfigFile
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &Config{}, ErrReadConfigFailed
	}
	return &cfg, nil
}
