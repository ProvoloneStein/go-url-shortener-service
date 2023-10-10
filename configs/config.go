package configs

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"os"
	"reflect"

	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	BaseURL             string `env:"BASE_URL" json:"base_url"`
	Addr                string `env:"SERVER_ADDRESS" json:"server_address"`
	FileStorage         string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN         string `env:"DATABASE_DSN" json:"database_dsn"`
	SigningKey          string `env:"SIGNING_KEY" json:"signing_key"`
	ConfigFileName      string `env:"CONFIG" json:"-"`
	TrustedSubnetString string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	TrustedSubnet       *net.IPNet
	GRPCPort            int  `env:"GRPC_PORT" json:"grpc_port"`
	EnableHTTPS         bool `env:"ENABLE_HTTPS" json:"enable_https"`
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

const defaultGRPCPort = 9090

func New() (AppConfig, error) {
	var config AppConfig

	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.FileStorage, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection dsn")
	flag.StringVar(&config.SigningKey, "k", signingKey, "service signing key")
	flag.BoolVar(&config.EnableHTTPS, "s", false, "service signing key")
	flag.StringVar(&config.ConfigFileName, "c", "configs/config.json", "config file name")
	flag.StringVar(&config.TrustedSubnetString, "t", "", "trusted subnet")
	flag.StringVar(&config.ConfigFileName, "config", "configs/config.json", "config file name")
	flag.IntVar(&config.GRPCPort, "p", defaultGRPCPort, "handlers_grpc port")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("ошибка при получении переменных окружения: %w", err)
	}
	err = getFromFile(&config)
	if err != nil {
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			return config, fmt.Errorf("ошибка при чтении файла конфигурации: %w", err)
		}
	}
	_, config.TrustedSubnet, err = net.ParseCIDR(config.TrustedSubnetString)
	if err != nil {
		return config, fmt.Errorf("ошибка парсинга CIDR %w", err)
	}

	return config, nil
}
