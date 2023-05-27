package main

import (
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/server"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
	"go.uber.org/zap"
)

import (
	"log"
)

func main() {
	config, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize("info"); err != nil {
		log.Fatal(err)
	}
	repos := repositories.NewRepository()
	services := services.NewService(config, repos)
	handler := handlers.NewHandler(services)
	logger.Log.Info(fmt.Sprintf("запускается сервер по адресу %s", config.Addr))
	if err = server.Run(config.Addr, handler.InitHandler()); err != nil {
		logger.Log.Fatal("ошибка при запуске сервера", zap.Error(err))
	}
}
