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

// filePerms - стандартные права файлового репозитория.
const filePerms = 0600

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	if buildVersion != "" {
		fmt.Printf("Build version: %s\n", buildVersion)
	} else {
		fmt.Println("Build version: N/A")
	}
	if buildDate != "" {
		fmt.Printf("Build date: %s\n", buildDate)
	} else {
		fmt.Println("Build date: N/A")
	}
	if buildCommit != "" {
		fmt.Printf("Build commit: %s\n", buildCommit)
	} else {
		fmt.Println("Build commit: N/A")
	}
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
			if deferErr := repos.Close(); deferErr != nil {
				logger.Error("ошибка при закрытии подключения базы данных.", zap.Error(deferErr))
			}
		}()
	case config.FileStorage == "":
		repos = repositories.NewLocalRepository(logger, config)
	default:
		file, repoInitErr := os.OpenFile(config.FileStorage, os.O_CREATE|os.O_RDWR, filePerms)
		if repoInitErr != nil {
			logger.Fatal("ошибка при попытке открытия файла", zap.Error(repoInitErr))
		}
		defer func() {
			if deferErr := file.Close(); deferErr != nil {
				logger.Error("ошибка при закрытии файлового репозитория.", zap.Error(deferErr))
			}
		}()
		repos, repoInitErr = repositories.NewFileRepository(config, logger, file)
		if repoInitErr != nil {
			logger.Fatal("ошибка при иницилизации файлового репозитория.", zap.Error(repoInitErr))
		}
	}
	services := services.NewService(logger, config, repos)
	handler := handlers.NewHandler(logger, services)
	logger.Info(fmt.Sprintf("запускается сервер по адресу %s", config.Addr))
	if err = server.Run(config.Addr, handler.InitHandler()); err != nil {
		logger.Fatal("ошибка при запуске сервера", zap.Error(err))
	}
}
