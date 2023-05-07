package services

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
)

//go:generate mockgen -source=services.go -destination=mocks/mock.go

type Shortener interface {
	CreateShortURL(fullURL string) (string, error)
	GetFullByID(shortURL string) (string, error)
}

type Service struct {
	Shortener
}

func NewService(cfg configs.AppConfig, repos *repositories.Repository) *Service {
	return &Service{
		Shortener: NewShortenerService(cfg, repos.Shortener),
	}
}
