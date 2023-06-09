package handlers

import (
	"context"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Service interface {
	CreateShortURL(ctx context.Context, fullURL string) (string, error)
	BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	GetFullByID(ctx context.Context, shortURL string) (string, error)
	Ping() error
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
	router.Get("/ping", h.pingDB)
	router.Post("/", h.createShortURL)
	router.Get("/{id}", h.getByShort)

	router.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", h.createShortURLByJSON)
			r.Post("/batch", h.batchCreateURLByJSON)
		})
	})

	return router
}

func (h *Handler) pingDB(w http.ResponseWriter, r *http.Request) {
	err := h.services.Ping()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
