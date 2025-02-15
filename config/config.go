package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApiServerConfig ApiServer
	DatabaseConfig  Database
}

// ApiServer представляет конфигурацию сервера API
type ApiServer struct {
	Host          string        `env:"API_SERVER_HOST" env-default:"localhost"`
	Port          string        `env:"API_SERVER_PORT" env-default:"8080"`
	AuthSecretKey string        `env:"API_SERVER_AUTH_SECRET_KEY" env-required:"true"`
	Timeout       time.Duration `env:"API_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout   time.Duration `env:"API_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

// Database представляет конфигурацию подключения к базе данных
type Database struct {
	Driver   string `env:"DB_DRIVER" env-default:"postgres"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-default:"db"`
	DBPort   string `env:"DB_PORT" env-default:"5432"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

// MustLoad загружает конфигурацию
func MustLoad() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err = cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	log.Println("Config loaded successfully")
	return cfg, nil
}

func (c *Config) validate() error {
	if c.ApiServerConfig.AuthSecretKey == "" {
		return fmt.Errorf("API_SERVER_AUTH_SECRET_KEY is required")
	}
	if c.ApiServerConfig.Host == "" || c.ApiServerConfig.Port == "" {
		return fmt.Errorf("API_SERVER_HOST and API_SERVER_PORT are required")
	}
	if c.DatabaseConfig.Driver == "" || c.DatabaseConfig.Name == "" || c.DatabaseConfig.User == "" || c.DatabaseConfig.Password == "" {
		return fmt.Errorf("database configuration is incomplete")
	}
	return nil
}
