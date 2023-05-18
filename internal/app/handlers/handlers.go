package handlers

import (
	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Service interface {
	CreateShortURL(fullURL string) (string, error)
	GetFullByID(shortURL string) (string, error)
}

type Handler struct {
	services Service
}

func NewHandler(services Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", h.createShortURL)
	router.Get("/{id}", h.getByShort)
	return router
}
