package handlers

import (
	"encoding/json"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type requestData struct {
	URL string `json:"url" valid:"url"`
}

type responseData struct {
	Result string `json:"result" valid:"url"`
}

func (h *Handler) createShortURLByJSON(w http.ResponseWriter, r *http.Request) {
	var requestBody requestData

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "application/x-gzip") {
		http.Error(w, "Неверный header", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(requestBody); err != nil {
		http.Error(w, "Неверое тело запроса", http.StatusBadRequest)
		return
	}
	res, err := h.services.CreateShortURL(requestBody.URL)
	if err != nil {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(&responseData{Result: res})
	if err != nil {
		h.logger.Error("ошибка при сериализации url", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(b); err != nil {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		return
	}
}

func (h *Handler) batchCreateURLByJSON(w http.ResponseWriter, r *http.Request) {
	var requestBody []models.BatchCreateRequest

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "application/x-gzip") {
		http.Error(w, "Неверный header", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(requestBody); err != nil {
		http.Error(w, "Неверое тело запроса", http.StatusBadRequest)
		return
	}
	res, err := h.services.BatchCreate(requestBody)
	if err != nil {
		h.logger.Error("ошибка при создании urls", zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		h.logger.Error("ошибка при сериализации url", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(b); err != nil {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		return
	}
}
