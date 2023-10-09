package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlersgrpc "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers_grpc"
	handlersrest "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers_rest"
	"github.com/ProvoloneStein/go-url-shortener-service/pkg/api/shorten"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/server"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
)

// filePerms - стандартные права файлового репозитория.
const filePerms = 0600
const shutdownTimeout = 30

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
	config, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logger.Initialize("info")
	if err != nil {
		log.Fatal(err)
	}
	var repos services.Repository
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
		file, err := os.OpenFile(config.FileStorage, os.O_CREATE|os.O_RDWR, filePerms)
		if err != nil {
			logger.Fatal("ошибка при попытке открытия файла", zap.Error(err))
		}
		defer func() {
			if deferErr := file.Close(); deferErr != nil {
				logger.Error("ошибка при закрытии файлового репозитория.", zap.Error(deferErr))
			}
		}()
		repos, err = repositories.NewFileRepository(config, logger, file)
		if err != nil {
			logger.Fatal("ошибка при иницилизации файлового репозитория.", zap.Error(err))
		}
	}
	services := services.NewService(logger, config, repos)
	handler := handlersrest.NewHandler(logger, services, config)
	srv := server.InitServer(config.Addr, handler.InitHandler())
	logger.Info(fmt.Sprintf("запускается сервер по адресу %s", config.Addr))

	serverErr := make(chan error, 1)

	go func() {
		if err = server.Run(logger, config.EnableHTTPS, srv); err != nil {
			serverErr <- err
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen:", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	shorten.RegisterShortenServer(grpcServer, handlersgrpc.NewServer(logger, services))

	logger.Info(fmt.Sprintf("запускается handlers_grpc сервер по порту %d", config.GRPCPort))

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			serverErr <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("start server error", zap.Error(err))
		}
	case sig := <-c:
		logger.Info(fmt.Sprintf("got signal %s", sig))
		logger.Info("Shutting everything down gracefully")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Graceful shutdown failed", zap.Error(err))
		}
	}
	logger.Info("Server shutdown successfully")
}
