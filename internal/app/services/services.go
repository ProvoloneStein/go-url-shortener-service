package services

import "github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"

type Shortener interface {
	CreateShortURL(fullUrl string) (string, error)
	GetFullByID(shortUrl string) (string, error)
}

type Service struct {
	Shortener
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Shortener: NewShortenerService(repos.Shortener),
	}
}
