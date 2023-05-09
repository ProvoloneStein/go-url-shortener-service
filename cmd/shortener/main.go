package main

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
)

func main() {
	config, err := configs.InitConfig()
	if err != nil {
		panic(err)
	}
	store := map[string]string{}
	repos := repositories.NewRepository(store)
	services := services.NewService(config, repos)
	handler := handlers.NewHandler(services)
	srv := new(app.Server)
	err = srv.Run(config.Addr, handler.InitHandler())
	if err != nil {
		panic(err)
	}
}
