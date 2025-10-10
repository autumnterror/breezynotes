package config

import (
	"flag"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	RedisAddr   string
	RedisPasswd string
	RedisDb     int
	Port        int
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
		RedisAddr   string `mapstructure:"redis_addr"`
		RedisPasswd string `mapstructure:"redis_passwd"`
		RedisDb     int    `mapstructure:"redis_db"`
		Port        int    `mapstructure:"port"`
		Mode        string `mapstructure:"mode"`
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
		RedisAddr:   cfg.RedisAddr,
		RedisPasswd: cfg.RedisPasswd,
		RedisDb:     cfg.RedisDb,
		Port:        cfg.Port,
	}, nil
}
