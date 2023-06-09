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
	"os"
)

import (
	"log"
)

func main() {
	config, err := configs.InitConfig()
	var repos services.Repository
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logger.Initialize("info")
	if err != nil {
		log.Fatal(err)
	}
	if config.FileStorage == "" {
		repos = repositories.NewLocalRepository()
	} else {
		file, err := os.OpenFile(config.FileStorage, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Fatal("ошибка при попытке открытия файла", zap.Error(err))
		}
		defer file.Close()
		repos, err = repositories.NewFileRepository(logger, file)
		if err != nil {
			logger.Fatal("ошибка при иницилизации репозитория.", zap.Error(err))
		}
	}
	services := services.NewService(config, repos)
	handler := handlers.NewHandler(logger, services)
	logger.Info(fmt.Sprintf("запускается сервер по адресу %s", config.Addr))
	if err = server.Run(config.Addr, handler.InitHandler()); err != nil {
		logger.Fatal("ошибка при запуске сервера", zap.Error(err))
	}
}
