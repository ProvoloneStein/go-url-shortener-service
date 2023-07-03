package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/asaskevich/govalidator"
)

type requestData struct {
	URL string `json:"url" valid:"url"`
}

type responseData struct {
	Result string `json:"result" valid:"url"`
}

const (
	contentTypeJSON = "application/json"
)

func (h *Handler) createShortURLByJSON(w http.ResponseWriter, r *http.Request) {
	var requestBody requestData

	ctx := r.Context()
	ct := r.Header.Get(contenntTypeHeader)
	if !strings.HasPrefix(ct, contentTypeJSON) && !strings.HasPrefix(ct, "application/x-gzip") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
	if err != nil && !errors.Is(err, repositories.ErrUniqueViolation) {
		h.logger.Error("ошибка при создании url", zap.Error(err))
		http.Error(w, "неверный запрос", http.StatusBadRequest)
		return
	}

	w.Header().Set(contenntTypeHeader, contentTypeJSON)

	if errors.Is(err, repositories.ErrUniqueViolation) {
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
	ct := r.Header.Get(contenntTypeHeader)
	if !strings.HasPrefix(ct, contentTypeJSON) && !strings.HasPrefix(ct, "application/x-gzip") {
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
	w.Header().Set(contenntTypeHeader, contentTypeJSON)
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
	w.Header().Set(contenntTypeHeader, contentTypeJSON)
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
	go h.services.DeleteUserURLsBatch(context.Background(), userID, reqBody)
	w.Header().Set(contenntTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusAccepted)
}
