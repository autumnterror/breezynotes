package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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

	if err := godotenv.Load(); err != nil {
		return nil, format.Error(op, err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./local-config/blocknote.yaml"
	}

	viper.SetConfigFile(configPath)

	var cfg struct {
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

	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	pw := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	db := os.Getenv("MONGO_INITDB_DATABASE")
	if user == "" || pw == "" || db == "" {
		return nil, format.Error(op, errors.New("missing environment variables"))
	}

	if cfg.Mode == "DEV" {
		log.Println(format.Struct(cfg), fmt.Sprintf("URI: mongodb://%s:%s@%s:%d/%s?authSource=admin",
			user, pw, cfg.DataSource, cfg.PortMongo, db))
	}

	return &Config{
		Uri: fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			user, pw, cfg.DataSource, cfg.PortMongo, db),
		Port: cfg.Port,
	}, nil
}
