package config

import (
	"errors"
	"log"
	"os"

	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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
	if err := godotenv.Load(); err != nil {
		return nil, format.Error(op, err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, format.Error(op, errors.New("CONFIG_PATH is not set"))
	}

	viper.SetConfigFile(configPath)

	var cfg struct {
		RedisAddr string `mapstructure:"redis_addr"`
		RedisDb   int    `mapstructure:"redis_db"`
		Port      int    `mapstructure:"port"`
		Mode      string `mapstructure:"mode"`
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}

	pw := os.Getenv("REDISPASSWD")
	if pw == "" {
		return nil, format.Error(op, errors.New("missing environment variables"))
	}

	if cfg.Mode == "DEV" {
		log.Println(format.Struct(cfg))
	}

	return &Config{
		RedisAddr:   cfg.RedisAddr,
		RedisPasswd: pw,
		RedisDb:     cfg.RedisDb,
		Port:        cfg.Port,
	}, nil
}
