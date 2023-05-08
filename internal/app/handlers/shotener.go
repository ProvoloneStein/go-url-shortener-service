package handlers

import (
	"io"
	"net/http"
	"strings"
)

func (h *Handler) createShortURL(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		http.Error(w, "Неверный header", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	res, err := h.services.CreateShortURL(string(body))
	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (h *Handler) getByShort(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	res, err := h.services.GetFullByID(shortURL)
	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
