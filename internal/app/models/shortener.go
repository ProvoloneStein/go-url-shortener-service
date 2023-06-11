package models

type BatchCreateRequest struct {
	URL  string `json:"original_url" valid:"url" db:"url"`
	UUID string `json:"correlation_id"`
}

type BatchCreateResponse struct {
	URL  string `json:"short_url" valid:"url"`
	UUID string `json:"correlation_id"`
}
