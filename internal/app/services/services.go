package services

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"net/url"
)

type Repository interface {
	Create(fullURL string) (string, error)
	GetByShort(shortURL string) (string, error)
}

type Service struct {
	cfg  configs.AppConfig
	repo Repository
}

func NewService(cfg configs.AppConfig, repo Repository) *Service {
	return &Service{cfg: cfg, repo: repo}
}

func (s *Service) CreateShortURL(fullURL string) (string, error) {
	shortID, err := s.repo.Create(fullURL)
	if err != nil {
		return "", err
	}
	shortURL, err := url.JoinPath(s.cfg.BaseURL, shortID)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *Service) GetFullByID(shortURL string) (string, error) {
	return s.repo.GetByShort(shortURL)
}
