package handlers

import (
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
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
	router.Use(logger.RequestLogger)
	router.Post("/", h.createShortURL)
	router.Get("/{id}", h.getByShort)

	router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.createShortURLByJSON)
	})
	return router
}
