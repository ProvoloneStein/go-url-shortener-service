package services

import (
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"net/url"
)

type Repository interface {
	Create(fullURL string) (string, error)
	BatchCreate(data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	GetByShort(shortURL string) (string, error)
	Ping() error
}

type Service struct {
	cfg  configs.AppConfig
	repo Repository
}

func NewService(cfg configs.AppConfig, repo Repository) *Service {
	return &Service{cfg: cfg, repo: repo}
}

func (s *Service) CreateShortURL(fullURL string) (string, error) {
	shortID, repoErr := s.repo.Create(fullURL)
	if repoErr != nil && !errors.Is(repoErr, repositories.ErrorUniqueViolation) {
		return "", repoErr
	}
	shortURL, err := url.JoinPath(s.cfg.BaseURL, shortID)
	if err != nil {
		return "", err
	}
	return shortURL, repoErr
}

func (s *Service) BatchCreate(data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	return s.repo.BatchCreate(data)
}

func (s *Service) GetFullByID(shortURL string) (string, error) {
	return s.repo.GetByShort(shortURL)
}

func (s *Service) Ping() error {
	return s.repo.Ping()
}
