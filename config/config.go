package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server `envPrefix:"Server_"`
		YaApp  `envPrefix:"YaApp_"`
		PG     `envPrefix:"PG_"`
		TG     `envPrefix:"TG_"`
	}

	Server struct {
		Port string `env:"Port"`
	}
	YaApp struct {
		Host        string `env:"host"`
		ClienID     string `env:"ClientID"`
		ClienSecret string `env:"ClientSecret"`
	}

	PG struct {
		URL string `env:"URL"`
	}

	TG struct {
		Endpoint string `env:"endpoint"`
		Token    string `env:"Token"`
		Webhook  string `env:"Webhook"`
	}
)

func NewConfig() (*Config, error) {
	loadEnv()
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}
	return &cfg, nil
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("%+v\n", err)
	}
}
