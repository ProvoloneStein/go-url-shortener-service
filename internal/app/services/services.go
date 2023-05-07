package services

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"
)

type Shortener interface {
	CreateShortURL(fullURL string) (string, error)
	GetFullByID(shortURL string) (string, error)
}

type Service struct {
	Shortener
}

func NewService(cfg configs.AppConfig, repos *repository.Repository) *Service {
	return &Service{
		Shortener: NewShortenerService(cfg, repos.Shortener),
	}
}
