package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ct := r.Header.Get("Content-Type")
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
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	res, err := h.services.CreateShortURL(ctx, userID, string(body))
	if err != nil && !errors.Is(err, repositories.ErrorUniqueViolation) {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if errors.Is(err, repositories.ErrorUniqueViolation) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err = w.Write([]byte(res)); err != nil {
		h.logger.Error("ошибка при записи ответа", zap.Error(err))
		return
	}
}

func (h *Handler) getByShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortURL := chi.URLParam(r, "id")
	_, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	res, err := h.services.GetFullByID(ctx, shortURL)
	if err != nil {
		if errors.Is(err, repositories.ErrURLNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, repositories.ErrDeleted) {
			http.Error(w, err.Error(), http.StatusGone)
		} else {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
		}
		return
	}
	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
