package repositories

import (
	"context"
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
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

func (r *LocalRepository) GenerateShortURL(ctx context.Context) (string, error) {

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	for {
		shortURL := randomString()
		if _, ok := r.store[shortURL]; !ok {
			return shortURL, nil
		}
	}
}

func (r *LocalRepository) Create(ctx context.Context, fullURL, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
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
		shortID, err := r.GenerateShortURL(ctx)
		if err != nil {
			return response, err
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
	return "", NewValueError(shortURL, ErrURLNotFound)
}

func (r *LocalRepository) Ping() error {
	return nil
}

func (r *LocalRepository) Close() error {
	return nil
}
