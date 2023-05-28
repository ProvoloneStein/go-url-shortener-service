package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/logger"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type getData struct {
	url string `json: "url" valid: "url"`
}

func (h *Handler) createShortURLByJSON(w http.ResponseWriter, r *http.Request) {
	var requestData getData
	var buf bytes.Buffer

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "Неверный header", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &requestData); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(requestData); err != nil {
		http.Error(w, "Неверое тело запроса", http.StatusBadRequest)
		return
	}
	res, err := h.services.CreateShortURL(requestData.url)
	if err != nil {
		logger.Log.Error("ошибка при создании url", zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write([]byte(res)); err != nil {
		logger.Log.Error("ошибка при создании url", zap.Error(err))
	}
}
