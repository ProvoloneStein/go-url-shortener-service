package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/go-chi/chi/v5"
)

// CreateShortURL - хэндлер создания короткого URL.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ct := r.Header.Get(contentTypeHeader)
	if !strings.HasPrefix(ct, "text/plain") && !strings.HasPrefix(ct, "application/x-gzip") {
		http.Error(w, "Неверный header запроса", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	res, err := h.services.CreateShortURL(ctx, userID, string(body))
	if err != nil && !errors.Is(err, repositories.ErrUniqueViolation) {
		h.logger.Error(defaultServiceError, zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set(contentTypeHeader, "text/plain; charset=utf-8")
	if errors.Is(err, repositories.ErrUniqueViolation) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err = w.Write([]byte(res)); err != nil {
		h.logger.Error("ошибка при записи ответа", zap.Error(err))
		return
	}
}

// GetByShort - хэндлер получения длинного URL по короткому URL.
func (h *Handler) GetByShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortURL := chi.URLParam(r, "id")
	fmt.Println(shortURL)
	_, err := getUserID(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	res, err := h.services.GetFullByID(ctx, shortURL)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrURLNotFound):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, repositories.ErrDeleted):
			http.Error(w, err.Error(), http.StatusGone)
		default:
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
		}
		return
	}
	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
