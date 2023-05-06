package main

import (
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/interfaces"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"net/http"
)

func main() {
	store := map[string]string{}
	repos := repository.NewRepository(store)
	services := services.NewService(repos)
	handler := interfaces.NewHandler(services)
	err := http.ListenAndServe(`:8080`, handler.InitHandler())
	if err != nil {
		panic(err)
	}
}
