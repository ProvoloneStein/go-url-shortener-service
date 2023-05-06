package interfaces

import (
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"net/http"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.mainHandler)
	return mux
}
