package repositories

import (
	"context"
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type LocalRepository struct {
	logger *zap.Logger
	cfg    configs.AppConfig
	store  map[string]string
}

func NewLocalRepository(logger *zap.Logger, cfg configs.AppConfig) *LocalRepository {
	return &LocalRepository{
		logger: logger,
		cfg:    cfg,
		store:  make(map[string]string),
	}
}

func (r *LocalRepository) Create(ctx context.Context, fullURL string) (string, error) {
	var shortURL string
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		for key, val := range r.store {
			if val == fullURL {
				return key, ErrorUniqueViolation
			}
		}
		for {
			shortURL = randomString()
			if _, ok := r.store[shortURL]; !ok {
				r.store[shortURL] = fullURL
				return shortURL, nil
			}
		}
	}
}

func (r *LocalRepository) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, val := range data {
			time.Sleep(time.Second)
			shortURL, err := r.Create(ctx, val.URL)
			if err == nil {
				resShortURL, err := url.JoinPath(r.cfg.BaseURL, shortURL)
				if err == nil {
					response = append(response, models.BatchCreateResponse{URL: resShortURL, UUID: val.UUID})
				}
			}
		}
		return response, nil
	}
}

func (r *LocalRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		fullURL, ok := r.store[shortURL]
		if ok {
			return fullURL, nil
		}
		return "", errors.New("url not found")
	}
}

func (r *LocalRepository) Ping() error {
	return nil
}
