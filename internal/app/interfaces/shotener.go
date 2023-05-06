package interfaces

import (
	"io"
	"net/http"
)

func (h *Handler) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte(res))
	} else if r.Method == http.MethodGet {
		shortURL := r.URL.Path[1:]
		res, err := h.services.GetFullByID(shortURL)
		if err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
	}
}
