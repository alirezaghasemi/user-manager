package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   Server
	Database Database
	Redis    Redis
	Swagger  Swagger
}

type Server struct {
	Host         string        `envconfig:"SERVER_HOST" default:"127.0.0.1"`
	Port         int           `envconfig:"SERVER_PORT" default:"6661"`
	WriteTimeout time.Duration `envconfig:"SERVER_WITH_TIMEOUT" default:"10s"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	Debug        bool          `envconfig:"SERVER_DEBUG" default:"true"`
}

type Database struct {
	Host              string `envconfig:"DATABASE_HOST"`
	Port              int    `envconfig:"DATABASE_PORT"`
	Username          string `envconfig:"DATABASE_USERNAME"`
	Password          string `envconfig:"DATABASE_PASSWORD"`
	Name              string `envconfig:"DATABASE_NAME"`
	SSLMode           string `envconfig:"DATABASE_SSLMODE" default:"disable"`
	MaxOpenConnection int    `envconfig:"DATABASE_MAX_OPEN_CONNECTION" default:"100"`
	Driver            string `envconfig:"DATABASE_DRIVER" default:"postgres"`
}

type Redis struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     string `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"`
}

type Swagger struct {
	Host    string   `envconfig:"SWAGGER_HOST_PORT" default:"127.0.0.1:6661"`
	Schemas []string `envconfig:"SWAGGER_SCHEMA" default:"http"`
}

func Load(envPath string) (*Config, error) {
	if envPath != "" {
		err := godotenv.Load(envPath)
		if err != nil {
			return nil, fmt.Errorf("could not load .env file: %w", err)
		}
	} else {
		_ = godotenv.Load()
	}

	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load env variable into config struct: %w", err)
	}

	return &cfg, nil
}
