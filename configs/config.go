package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	BaseURL        string `env:"BASE_URL" json:"base_url"`
	Addr           string `env:"SERVER_ADDRESS" json:"server_address"`
	FileStorage    string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN    string `env:"DATABASE_DSN" json:"database_dsn"`
	SigningKey     string `env:"SIGNING_KEY" json:"signing_key"`
	ConfigFileName string `env:"CONFIG" json:"-"`
	EnableHTTPS    bool   `env:"ENABLE_HTTPS" json:"enable_https"`
}

const signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"

func getFromFile(cfg *AppConfig) error {
	plan, err := os.ReadFile(cfg.ConfigFileName)
	if err != nil {
		return fmt.Errorf("file read err: %w", err)
	}
	data := AppConfig{}
	err = json.Unmarshal(plan, &data)
	if err != nil {
		return fmt.Errorf("file unmmarshal err: %w", err)
	}

	dataValueOf := reflect.ValueOf(&data)

	cfgValueOf := reflect.ValueOf(cfg)
	for i := 0; i < reflect.TypeOf(*cfg).NumField(); i++ {
		if cfgValueOf.Elem().Field(i).IsZero() {
			dataFieldValue := dataValueOf.Elem().FieldByName(reflect.TypeOf(*cfg).Field(i).Name)
			cfgValueOf.Elem().Field(i).Set(dataFieldValue)
		}
	}
	return nil
}

func InitConfig() (AppConfig, error) {
	var config AppConfig

	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.FileStorage, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection dsn")
	flag.StringVar(&config.SigningKey, "k", signingKey, "service signing key")
	flag.BoolVar(&config.EnableHTTPS, "s", false, "service signing key")
	flag.StringVar(&config.ConfigFileName, "c", "configs/config.json", "config file name")
	flag.StringVar(&config.ConfigFileName, "config", "configs/config.json", "config file name")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("ошибка при получении переменных окружения: %w", err)
	}
	err = getFromFile(&config)
	if err != nil {
		return config, fmt.Errorf("ошибка при чтении файла конфигурации: %w", err)
	}

	return config, nil
}
