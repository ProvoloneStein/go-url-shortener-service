package repositories

import (
	"context"
	"errors"
	"net/url"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
)

type LocalRepository struct {
	logger *zap.Logger
	store  map[string][3]string
	cfg    configs.AppConfig
}

func NewLocalRepository(logger *zap.Logger, cfg configs.AppConfig) *LocalRepository {
	return &LocalRepository{
		logger: logger,
		cfg:    cfg,
		store:  make(map[string][3]string),
	}
}

func (r *LocalRepository) validateUniqueShortURL(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return defaultRepoErrWrapper(ctx.Err())
	default:
	}
	if _, ok := r.store[shortURL]; !ok {
		return nil
	}
	return ErrShortURLExists
}

func (r *LocalRepository) Create(ctx context.Context, userID, fullURL, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", defaultRepoErrWrapper(ctx.Err())
	default:
	}
	if err := r.validateUniqueShortURL(ctx, shortURL); err != nil {
		return "", err
	}
	for key, val := range r.store {
		if val[0] == fullURL {
			return key, ErrUniqueViolation
		}
	}
	r.store[shortURL] = [3]string{fullURL, userID, "f"}
	return shortURL, nil
}

func (r *LocalRepository) BatchCreate(ctx context.Context,
	data []models.BatchCreateData) ([]models.BatchCreateResponse, error) {
	select {
	case <-ctx.Done():
		return nil, defaultRepoErrWrapper(ctx.Err())
	default:
	}
	for _, val := range data {
		if err := r.validateUniqueShortURL(ctx, val.ShortURL); err != nil {
			return []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: val.ShortURL, UUID: val.UUID}}, err
		}
	}
	response := make([]models.BatchCreateResponse, 0, len(data))
	for _, val := range data {
		shortURL, err := url.JoinPath(r.cfg.BaseURL, val.ShortURL)
		if err != nil {
			r.logger.Error("ошибка при формировании url", zap.Error(err))
			return response, defaultRepoErrWrapper(err)
		}
		_, err = r.Create(ctx, val.UserID, val.URL, val.ShortURL)
		if err != nil && !errors.Is(err, ErrUniqueViolation) {
			r.logger.Error("ошибка при записи url", zap.Error(err))
			return response, defaultRepoErrWrapper(err)
		}
		response = append(response, models.BatchCreateResponse{ShortURL: shortURL, UUID: val.UUID})
	}
	return response, nil
}

func (r *LocalRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", defaultRepoErrWrapper(ctx.Err())
	default:
	}
	data, ok := r.store[shortURL]
	if ok {
		if data[2] == "t" {
			return "", ErrDeleted
		}
		return data[0], nil
	}
	return "", errWithVal(ErrURLNotFound, shortURL)
}

func (r *LocalRepository) Ping() error {
	return nil
}

func (r *LocalRepository) Close() error {
	return nil
}

func (r *LocalRepository) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	var result []models.GetURLResponse
	select {
	case <-ctx.Done():
		return nil, defaultRepoErrWrapper(ctx.Err())
	default:
	}
	for key, val := range r.store {
		if userID == val[1] && val[2] == "f" {
			result = append(result, models.GetURLResponse{ShortURL: key, URL: val[0]})
		}
	}
	return result, nil
}

func (r *LocalRepository) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error {
	select {
	case <-ctx.Done():
		return defaultRepoErrWrapper(ctx.Err())
	default:
	}
	for _, short := range data {
		row, ok := r.store[short]
		if ok {
			if row[1] == userID {
				row[2] = "t"
			}
		}
	}
	return nil
}
