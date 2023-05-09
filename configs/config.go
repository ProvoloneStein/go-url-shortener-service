package configs

import (
	"flag"
	"github.com/caarlos0/env/v8"
	"log"
)

type AppConfig struct {
	BaseURL string `env:"BASE_URL"`
	Addr    string `env:"SERVER_ADDRESS"`
}

func InitConfig() (AppConfig, error) {
	var Config AppConfig

	flag.StringVar(&Config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&Config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}
	return Config, nil
}
