package configs

type AppConfig struct {
	BaseURL string
	Addr    string
}

func InitConfig() (AppConfig, error) {
	var Config AppConfig
	Config.BaseURL = "http://localhost:8080"
	Config.Addr = "localhost:8080"
	return Config, nil
}
