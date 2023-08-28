package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
)

const contentTypeHeader = "Content-Type"

// Service -  интерфейс с необходимым набором методов сервиса.
//
//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Service interface {
	CreateShortURL(ctx context.Context, userID, fullURL string) (string, error)
	BatchCreate(ctx context.Context, userID string, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	// GetFullByID - получить полный URL по короткому URL.
	GetFullByID(ctx context.Context, shortURL string) (string, error)
	// GetListByUser - получить список связей короткий/длинный URL по пользователю.
	GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error)
	// DeleteUserURLsBatch - удалить связть короткий/длинный URL.
	DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error
	// GenerateToken - получить токен авторизации.
	GenerateToken(ctx context.Context) (string, error)
	// ParseToken - расшифровать токен авторизации.
	ParseToken(accessToken string) (string, error)
	// Ping - проверить доступность сервиса.
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
	router.Use(userIdentity(h.services, h.logger))
	router.Mount("/debug", middleware.Profiler())
	router.Get("/ping", h.pingDB)

	router.Post("/", h.CreateShortURL)
	router.Get("/{id}", h.GetByShort)

	router.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", h.CreateShortURLByJSON)
			r.Post("/batch", h.BatchCreateURLByJSON)
		})
		r.Route("/user", func(r chi.Router) {
			r.Get("/urls", h.GetUserURLs)
			r.Delete("/urls", h.DeleteUserURLsBatch)
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
