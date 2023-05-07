package services

import (
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"
)

type ShortenerService struct {
	cfg  configs.AppConfig
	repo repository.Shortener
}

func NewShortenerService(cfg configs.AppConfig, repo repository.Shortener) *ShortenerService {
	return &ShortenerService{cfg: cfg, repo: repo}
}

func (s *ShortenerService) CreateShortURL(fullURL string) (string, error) {
	shortID, err := s.repo.Create(fullURL)
	if err != nil {
		return "", err
	}
	shortURL := s.cfg.BaseURL + "/" + shortID
	return shortURL, nil
}

func (s *ShortenerService) GetFullByID(shortURL string) (string, error) {
	return s.repo.GetByShort(shortURL)
}
