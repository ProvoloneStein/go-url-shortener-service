package repository

type Shortener interface {
	Create(fullURL string) (string, error)
	GetByShort(shortURL string) (string, error)
}

type Repository struct {
	Shortener
}

func NewRepository(store map[string]string) *Repository {
	return &Repository{
		Shortener: NewShortenerRepository(store),
	}
}
