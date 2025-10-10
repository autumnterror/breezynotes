package config

import (
	"flag"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	AddrAuth      string
	AddrBlockNote string
	AddrRedis     string
	Timeout       time.Duration
	Backoff       time.Duration
	RetriesCount  int
	Port          int
}

// MustSetup return config and panic if error
func MustSetup() *Config {
	cfg, err := setup()
	if err != nil {
		log.Panic(err)
	}
	return cfg
}

// setup create config structure
func setup() (*Config, error) {
	const op = "config.setup"
	configPath := flag.String("cfg", "./local-config/gateway.yaml", "path to config file")
	flag.Parse()
	viper.SetConfigFile(*configPath)

	var cfg struct {
		AddrAuth      string        `mapstructure:"addr_auth"`
		AddrBlockNote string        `mapstructure:"addr_blocknote"`
		AddrRedis     string        `mapstructure:"addr_redis"`
		Timeout       time.Duration `mapstructure:"timeout"`
		Backoff       time.Duration `mapstructure:"backoff"`
		RetriesCount  int           `mapstructure:"retries_count"`
		Port          int           `mapstructure:"port"`
		Mode          string        `mapstructure:"mode"`
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}

	if cfg.Mode == "DEV" {
		log.Println(format.Struct(cfg))
	}

	return &Config{
		AddrAuth:      cfg.AddrAuth,
		AddrBlockNote: cfg.AddrBlockNote,
		AddrRedis:     cfg.AddrRedis,
		Timeout:       cfg.Timeout,
		Backoff:       cfg.Backoff,
		RetriesCount:  cfg.RetriesCount,
		Port:          cfg.Port,
	}, nil
}
