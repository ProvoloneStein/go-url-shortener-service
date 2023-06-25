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
	Create(ctx context.Context, userID, fullURL, shortURL string) (string, error)
	BatchCreate(ctx context.Context, data []models.BatchCreateData) ([]models.BatchCreateResponse, error)
	GetByShort(ctx context.Context, shortURL string) (string, error)
	GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error)
	ValidateUniqueUser(ctx context.Context, userID string) error
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

func (s *Service) CreateShortURL(ctx context.Context, userID, fullURL string) (string, error) {
	for {
		shortID := repositories.RandomString()
		shortID, repoErr := s.repo.Create(ctx, userID, fullURL, shortID)
		if repoErr != nil {
			if errors.Is(repoErr, repositories.ErrShortURLExists) {
				continue
			}
			if !errors.Is(repoErr, repositories.ErrorUniqueViolation) {
				return "", repoErr
			}
		}
		shortURL, err := url.JoinPath(s.cfg.BaseURL, shortID)
		if err != nil {
			return "", err
		}
		return shortURL, repoErr
	}
}

func (s *Service) BatchCreate(ctx context.Context, userID string, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var queryData []models.BatchCreateData
	var res []models.BatchCreateResponse

	for dataIndex, _ := range data {
		shortID := repositories.RandomString()
		queryData = append(queryData, models.BatchCreateData{ShortURL: shortID, URL: data[dataIndex].URL, UUID: data[dataIndex].UUID, UserID: userID})
	}

generator:
	for {
		for resIndex := range res {
			shortID := repositories.RandomString()
			// проверяем что не задублировали shortID
			for queryIndex := range queryData {
				if shortID == queryData[queryIndex].ShortURL {
					res = append(res[resIndex:])
					continue generator
				}
			}
			for queryIndex := range queryData {
				if res[resIndex].ShortURL == queryData[queryIndex].ShortURL {
					queryData[queryIndex].ShortURL = shortID
				}
			}
		}
		res, err := s.repo.BatchCreate(ctx, queryData)
		if err != nil {
			if errors.Is(err, repositories.ErrShortURLExists) {
				continue
			}
			return nil, err
		}
		return res, err
	}
}

func (s *Service) GetFullByID(ctx context.Context, userID, shortURL string) (string, error) {
	// todo оставлю userID до тестов. хз надо ли учитывать
	return s.repo.GetByShort(ctx, shortURL)
}

func (s *Service) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	return s.repo.GetListByUser(ctx, userID)
}

func (s *Service) Ping() error {
	return s.repo.Ping()
}
