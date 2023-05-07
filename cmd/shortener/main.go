package main

import (
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/interfaces"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"net/http"
)

func main() {
	config, _ := configs.InitConfig()
	store := map[string]string{}
	repos := repository.NewRepository(store)
	services := services.NewService(config, repos)
	handler := interfaces.NewHandler(services)
	err := http.ListenAndServe(config.Addr, handler.InitHandler())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
