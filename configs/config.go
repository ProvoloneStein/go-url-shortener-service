package configs

import "flag"

type AppConfig struct {
	BaseURL string
	Addr    string
}

func InitConfig() (AppConfig, error) {
	var Config AppConfig
	flag.StringVar(&Config.BaseURL, "b", "http://localhost:8080", "short URL base address")
	flag.StringVar(&Config.Addr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	return Config, nil
}
