package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Uri                  string
	TokenKey             string
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
	Port                 int
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
		configPath = "./local-config/auth.yaml"
	}

	viper.SetConfigFile(configPath)

	var cfg struct {
		DataSource           string        `mapstructure:"data_source"`
		PortPostgres         int           `mapstructure:"port_postgres"`
		AccessTokenLifeTime  time.Duration `mapstructure:"access_token_life"`
		RefreshTokenLifeTime time.Duration `mapstructure:"refresh_token_life"`
		Port                 int           `mapstructure:"port"`
		Mode                 string        `mapstructure:"mode"`
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}
	user := os.Getenv("POSTGRES_USER")
	pw := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	token := os.Getenv("TOKEN_KEY")
	if user == "" || pw == "" || db == "" || token == "" {
		return nil, format.Error(op, errors.New("missing environment variables"))
	}

	if cfg.Mode == "DEV" {
		log.Println(format.Struct(cfg), fmt.Sprintf("URI: postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, pw, cfg.DataSource, cfg.PortPostgres, db))
	}

	return &Config{
		Uri: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, pw, cfg.DataSource, cfg.PortPostgres, db),
		TokenKey:             token,
		AccessTokenLifeTime:  cfg.AccessTokenLifeTime,
		RefreshTokenLifeTime: cfg.RefreshTokenLifeTime,
		Port:                 cfg.Port,
	}, nil
}
