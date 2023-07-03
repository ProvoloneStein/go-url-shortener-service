package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
)

type Repository interface {
	Create(ctx context.Context, userID, fullURL, shortURL string) (string, error)
	BatchCreate(ctx context.Context, data []models.BatchCreateData) ([]models.BatchCreateResponse, error)
	GetByShort(ctx context.Context, shortURL string) (string, error)
	GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error)
	DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error
	ValidateUniqueUser(ctx context.Context, userID string) error
	Ping() error
	Close() error
}

type Service struct {
	logger *zap.Logger
	repo   Repository
	cfg    configs.AppConfig
}

const defaultServiceError = "service:"

func defaultServiceErrWrapper(err error) error {
	return fmt.Errorf("%s %w", defaultServiceError, err)
}

func NewService(logger *zap.Logger, cfg configs.AppConfig, repo Repository) *Service {
	service := Service{logger: logger, cfg: cfg, repo: repo}
	return &service
}

func (s *Service) CreateShortURL(ctx context.Context, userID, fullURL string) (string, error) {
	for {
		shortID := repositories.RandomString()
		shortID, repoErr := s.repo.Create(ctx, userID, fullURL, shortID)
		if repoErr != nil {
			if errors.Is(repoErr, repositories.ErrShortURLExists) {
				continue
			}
			if !errors.Is(repoErr, repositories.ErrUniqueViolation) {
				return "", defaultServiceErrWrapper(repoErr)
			}
		}
		shortURL, err := url.JoinPath(s.cfg.BaseURL, shortID)
		if err != nil {
			return "", defaultServiceErrWrapper(err)
		}
		if repoErr != nil {
			return shortURL, defaultServiceErrWrapper(repoErr)
		}
		return shortURL, nil
	}
}

func (s *Service) BatchCreate(ctx context.Context, userID string,
	data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {

	queryData := make([]models.BatchCreateData, 0, len(data))
	for dataIndex := range data {
		shortID := repositories.RandomString()
		obj := models.BatchCreateData{ShortURL: shortID, URL: data[dataIndex].URL, UUID: data[dataIndex].UUID, UserID: userID}
		queryData = append(queryData, obj)
	}

	for {
		res, err := s.repo.BatchCreate(ctx, queryData)
		if err != nil {
			if errors.Is(err, repositories.ErrShortURLExists) {
			generator:
				for resIndex := range res {
					shortID := repositories.RandomString()
					// проверяем что не задублировали shortID
					for queryIndex := range queryData {
						if shortID == queryData[queryIndex].ShortURL {
							if len(res) > 1 {
								res = res[resIndex:]
							} else {
								res = nil
							}
							continue generator
						}
					}
					for queryIndex := range queryData {
						if res[resIndex].ShortURL == queryData[queryIndex].ShortURL {
							queryData[queryIndex].ShortURL = shortID
						}
					}
				}
				continue
			}
			return nil, defaultServiceErrWrapper(err)
		}
		return res, nil
	}
}

func (s *Service) GetFullByID(ctx context.Context, shortURL string) (string, error) {
	row, err := s.repo.GetByShort(ctx, shortURL)
	if err != nil {
		return row, defaultServiceErrWrapper(err)
	}
	return row, nil
}

func (s *Service) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	list, err := s.repo.GetListByUser(ctx, userID)
	if err != nil {
		return list, defaultServiceErrWrapper(err)
	}
	return list, nil
}

func (s *Service) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) {
	if err := s.repo.DeleteUserURLsBatch(ctx, userID, data); err != nil {
		s.logger.Error("ошибка при удалении", zap.Error(err))
	}
}

func (s *Service) Ping() error {
	if err := s.repo.Ping(); err != nil {
		return defaultServiceErrWrapper(err)
	}
	return nil
}
