package handlers

import (
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Service interface {
	CreateShortURL(fullURL string) (string, error)
	GetFullByID(shortURL string) (string, error)
}

type Handler struct {
	logger   *zap.Logger
	services Service
}

func NewHandler(logger *zap.Logger, services Service) *Handler {
	return &Handler{logger: logger, services: services}
}

func (h *Handler) InitHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(logger.RequestLogger(h.logger))
	router.Use(gzipReadWriterHandler(h.logger))
	router.Post("/", h.createShortURL)
	router.Get("/{id}", h.getByShort)

	router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.createShortURLByJSON)
	})
	return router
}
