package handlers

import (
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", h.createShortURL)
	router.Get("/{id}", h.getByShort)
	return router
}
