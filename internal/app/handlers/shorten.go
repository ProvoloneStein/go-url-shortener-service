package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
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

	ctx := r.Context()
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "application/x-gzip") {
		http.Error(w, "неверный header запоса", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "неверое тело запроса", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(requestBody); err != nil {
		http.Error(w, "ошибка валидации тела запроса", http.StatusBadRequest)
		return
	}
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	res, err := h.services.CreateShortURL(ctx, userID, requestBody.URL)
	if err != nil && !errors.Is(err, repositories.ErrorUniqueViolation) {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		http.Error(w, "неверный запрос", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, repositories.ErrorUniqueViolation) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	b, err := json.Marshal(&responseData{Result: res})
	if err != nil {
		h.logger.Error("ошибка при сериализации url", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(b); err != nil {
		h.logger.Error("ошибка при записи ответа", zap.Error(err))
		return
	}
}

func (h *Handler) batchCreateURLByJSON(w http.ResponseWriter, r *http.Request) {
	var requestBody []models.BatchCreateRequest

	ctx := r.Context()
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
		http.Error(w, "Неверное тело запрос", http.StatusBadRequest)
		return
	}
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	res, err := h.services.BatchCreate(ctx, userID, requestBody)
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

func (h *Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	res, err := h.services.GetListByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNoContent)
		} else if errors.Is(err, repositories.ErrDeleted) {
			http.Error(w, err.Error(), http.StatusGone)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		h.logger.Error("ошибка при сериализации url", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(b); err != nil {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		return
	}
}

func (h *Handler) deleteUserURLsBatch(w http.ResponseWriter, r *http.Request) {
	var reqBody []string

	ctx := r.Context()
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "Неверное тело запрос", http.StatusBadRequest)
		return
	}
	err = h.services.DeleteUserURLsBatch(ctx, userID, reqBody)
	if err != nil {
		h.logger.Error(fmt.Sprintf("error while DeleteURLBatch: %s", err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

}
