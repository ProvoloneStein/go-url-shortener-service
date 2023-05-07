package main

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"net/http"
)

func main() {
	config, _ := configs.InitConfig()
	store := map[string]string{}
	repos := repositories.NewRepository(store)
	services := services.NewService(config, repos)
	handler := handlers.NewHandler(services)
	err := http.ListenAndServe(config.Addr, handler.InitHandler())
	if err != nil {
		panic(err)
	}
}
