package repositories

import (
	"errors"
	"math/rand"
)

type Repository struct {
	store map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		store: make(map[string]string),
	}
}

func (r *Repository) Create(fullURL string) (string, error) {
	var shortURL string
	for {
		shortURL = randomString()
		if _, ok := r.store[shortURL]; !ok {
			r.store[shortURL] = fullURL
			return shortURL, nil
		}
	}
}

func (r *Repository) GetByShort(shortURL string) (string, error) {
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", errors.New("url not found")
}

func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
