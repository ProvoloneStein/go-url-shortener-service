package repository

type Shortener interface {
	Create(fullUrl string) (string, error)
	GetByShort(shortUrl string) (string, error)
}

type Repository struct {
	Shortener
}

func NewRepository(store map[string]string) *Repository {
	return &Repository{
		Shortener: NewShortenerRepository(store),
	}
}
