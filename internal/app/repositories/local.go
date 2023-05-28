package repositories

import (
	"errors"
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

func (r *LocalRepository) GetByShort(shortURL string) (string, error) {
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", errors.New("url not found")
}
