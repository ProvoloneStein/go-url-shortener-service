package repositories

import (
	"fmt"
	"math/rand"
)

type ShortenerRepository struct {
	store map[string]string
}

func NewShortenerRepository(store map[string]string) *ShortenerRepository {
	return &ShortenerRepository{store: store}
}

func (r *ShortenerRepository) Create(fullURL string) (string, error) {
	var shortURL string
	for {
		shortURL = RandomString()
		_, ok := r.store[shortURL]
		if !ok {
			r.store[shortURL] = fullURL
			return shortURL, nil
		}
	}
}

func (r *ShortenerRepository) GetByShort(shortURL string) (string, error) {
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", fmt.Errorf("url not found")
}

func RandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
