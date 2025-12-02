package config

import (
	"fmt"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Uri  string
	Port int
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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./local-config/blocknote.yaml"
	}

	viper.SetConfigFile(configPath)

	var cfg struct {
		Db         string
		Pw         string
		User       string
		DataSource string `mapstructure:"data_source"`
		PortMongo  int    `mapstructure:"port_mongo"`
		Uri        string
		Port       int
		Mode       string
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}

	if cfg.Mode == "DEV" {
		log.Println(format.Struct(cfg), fmt.Sprintf("URI: mongodb://%s:%s@%s:%d/%s?authSource=admin",
			cfg.User, cfg.Pw, cfg.DataSource, cfg.PortMongo, cfg.Db))
	}

	return &Config{
		Uri: fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			cfg.User, cfg.Pw, cfg.DataSource, cfg.PortMongo, cfg.Db),
		Port: cfg.Port,
	}, nil
}
