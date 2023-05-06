package services

import "github.com/ProvoloneStein/go-url-shortener-service/internal/app/repository"

type ShortenerService struct {
	repo repository.Shortener
}

func NewShortenerService(repo repository.Shortener) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) CreateShortURL(fullUrl string) (string, error) {
	return s.repo.Create(fullUrl)
}

func (s *ShortenerService) GetFullByID(shortUrl string) (string, error) {
	return s.repo.GetByShort(shortUrl)
}
