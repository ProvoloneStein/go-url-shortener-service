package repository

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

func (r *ShortenerRepository) Create(fullUrl string) (string, error) {
	var shortUrl string
	for {
		shortUrl = RandomString()
		_, ok := r.store[shortUrl]
		if !ok {
			r.store[shortUrl] = fullUrl
			return shortUrl, nil
		}
	}
}

func (r *ShortenerRepository) GetByShort(shortUrl string) (string, error) {
	fullUrl, ok := r.store[shortUrl]
	if ok {
		return fullUrl, nil
	}
	return "", fmt.Errorf("Ошибочка")
}

func RandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
