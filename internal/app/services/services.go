package services

import (
	"context"
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"go.uber.org/zap"
	"net/url"
)

type Repository interface {
	Create(ctx context.Context, fullURL, shortURL string) (string, error)
	GenerateShortURL(ctx context.Context) (string, error)
	BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	GetByShort(ctx context.Context, shortURL string) (string, error)
	Ping() error
	Close() error
}

type Service struct {
	logger *zap.Logger
	cfg    configs.AppConfig
	repo   Repository
}

func NewService(logger *zap.Logger, cfg configs.AppConfig, repo Repository) *Service {
	return &Service{logger: logger, cfg: cfg, repo: repo}
}

func (s *Service) CreateShortURL(ctx context.Context, fullURL string) (string, error) {
	shortID, err := s.repo.GenerateShortURL(ctx)
	if err != nil {
		return "", err
	}
	shortID, repoErr := s.repo.Create(ctx, fullURL, shortID)
	if repoErr != nil && !errors.Is(repoErr, repositories.ErrorUniqueViolation) {
		return "", repoErr
	}
	shortURL, err := url.JoinPath(s.cfg.BaseURL, shortID)
	if err != nil {
		return "", err
	}
	return shortURL, repoErr
}

func (s *Service) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	return s.repo.BatchCreate(ctx, data)
}

func (s *Service) GetFullByID(ctx context.Context, shortURL string) (string, error) {
	return s.repo.GetByShort(ctx, shortURL)
}

func (s *Service) Ping() error {
	return s.repo.Ping()
}
