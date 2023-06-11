package configs

import (
	"flag"
	"fmt"
)

import (
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	BaseURL     string `env:"BASE_URL"`
	Addr        string `env:"SERVER_ADDRESS"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func InitConfig() (AppConfig, error) {
	var config AppConfig

	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.FileStorage, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection dsn")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("ошибка при получении переменных окружения: %s", err)
	}
	return config, err
}
