package configs

import (
	"flag"
	"fmt"
)

import (
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	BaseURL string `env:"BASE_URL"`
	Addr    string `env:"SERVER_ADDRESS"`
}

func InitConfig() (AppConfig, error) {
	var config AppConfig

	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("ошибка при получении переменных окружения: %s", err)
	}
	return config, err
}
