package models

type BatchCreateRequest struct {
	URL  string `json:"original_url" db:"url"`
	UUID string `json:"correlation_id"`
}

type BatchCreateData struct {
	URL      string `json:"-" db:"url"`
	ShortURL string `json:"short_url" db:"shorten"`
	UUID     string `json:"correlation_id" db:"correlation_id"`
	UserID   string `json:"user_id" db:"user_id"`
}

type BatchCreateResponse struct {
	ShortURL string `json:"short_url" db:"shorten"`
	UUID     string `json:"correlation_id" db:"correlation_id"`
}

type GetURLResponse struct {
	ShortURL string `json:"short_url" db:"shorten"`
	URL      string `json:"original_url" db:"url"`
}
