package main

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/server"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
)

import (
	"log"
)

func main() {
	config, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	repos := repositories.NewRepository()
	services := services.NewService(config, repos)
	handler := handlers.NewHandler(services)
	if err = server.Run(config.Addr, handler.InitHandler()); err != nil {
		log.Fatal(err)
	}
}
