package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/server"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
)

const filePerms = 0600

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
	switch {
	case config.DatabaseDSN != "":
		repos, err = repositories.NewDBRepository(logger, config)
		if err != nil {
			logger.Fatal("ошибка при иницилизации базы данных.", zap.Error(err))
		}
		defer func() {
			if err := repos.Close(); err != nil {
				logger.Error("ошибка при закрытии подключения базы данных.", zap.Error(err))
			}
		}()
	case config.FileStorage == "":
		repos = repositories.NewLocalRepository(logger, config)
	default:
		file, err := os.OpenFile(config.FileStorage, os.O_CREATE|os.O_RDWR, filePerms)
		if err != nil {
			logger.Fatal("ошибка при попытке открытия файла", zap.Error(err))
		}
		defer func() {
			if err := file.Close(); err != nil {
				logger.Error("ошибка при закрытии файлового репозитория.", zap.Error(err))
			}
		}()
		repos, err = repositories.NewFileRepository(config, logger, file)
		if err != nil {
			logger.Fatal("ошибка при иницилизации файлового репозитория.", zap.Error(err))
		}
	}
	services := services.NewService(logger, config, repos)
	handler := handlers.NewHandler(logger, services)
	logger.Info(fmt.Sprintf("запускается сервер по адресу %s", config.Addr))
	if err = server.Run(config.Addr, handler.InitHandler()); err != nil {
		logger.Fatal("ошибка при запуске сервера", zap.Error(err))
	}
}
