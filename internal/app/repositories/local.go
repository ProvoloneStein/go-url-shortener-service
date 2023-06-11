package repositories

import (
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
)

type LocalRepository struct {
	store map[string]string
}

func NewLocalRepository() *LocalRepository {
	return &LocalRepository{
		store: make(map[string]string),
	}
}

func (r *LocalRepository) Create(fullURL string) (string, error) {
	var shortURL string
	for {
		shortURL = randomString()
		if _, ok := r.store[shortURL]; !ok {
			r.store[shortURL] = fullURL
			return shortURL, nil
		}
	}
}

func (r *LocalRepository) BatchCreate(data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	for _, val := range data {
		shortURL, err := r.Create(val.URL)
		if err == nil {
			response = append(response, models.BatchCreateResponse{URL: shortURL, UUID: val.UUID})
		}
	}
	return response, nil
}

func (r *LocalRepository) GetByShort(shortURL string) (string, error) {
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", errors.New("url not found")
}

func (r *LocalRepository) Ping() error {
	return nil
}
