package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"go.uber.org/zap"
	"net/url"
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

func (r *LocalRepository) validateUniqueShortURL(ctx context.Context, shortURL string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if _, ok := r.store[shortURL]; !ok {
		return nil
	}
	return ErrShortURLExists
}

func (r *LocalRepository) Create(ctx context.Context, fullURL, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	if err := r.validateUniqueShortURL(ctx, shortURL); err != nil {
		return "", err
	}
	for key, val := range r.store {
		if val == fullURL {
			return key, ErrorUniqueViolation
		}
	}
	r.store[shortURL] = fullURL
	return shortURL, nil
}

func (r *LocalRepository) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for _, val := range data {
		for {
			shortID := services.RandomString()
			if _, ok := r.store[shortID]; ok {
				continue
			}
			shortURL, err := url.JoinPath(r.cfg.BaseURL, shortID)
			if err != nil {
				r.logger.Error("ошибка при формировании url", zap.Error(err))
				return response, err
			}
			_, err = r.Create(ctx, val.URL, shortID)
			if err != nil && !errors.Is(err, ErrorUniqueViolation) {
				r.logger.Error("ошибка при записи url", zap.Error(err))
				return response, err
			}
			response = append(response, models.BatchCreateResponse{ShortURL: shortURL, UUID: val.UUID})
			break
		}
	}
	return response, nil
}

func (r *LocalRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", fmt.Errorf("%w: %s", ErrURLNotFound, shortURL)
}

func (r *LocalRepository) Ping() error {
	return nil
}

func (r *LocalRepository) Close() error {
	return nil
}
