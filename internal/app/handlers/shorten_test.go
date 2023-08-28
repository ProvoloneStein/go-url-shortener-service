package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	mock_handlers "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers/mocks"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_CreateShortURLByJSON(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, userID, fullURL string)

	type want struct {
		contentType string
		body        string
		statusCode  int
	}

	tests := []struct {
		name         string
		body         map[string]string
		userID       string
		mockBehavior mockBehavior
		contentType  string
		want         want
	}{
		{
			name:        "Wrong Content Type",
			contentType: "text/plain",
			userID:      "123",
			body: map[string]string{
				"url": "https://ya.ru",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", nil).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный header запроса\n",
			},
		},
		{
			name:        "err ErrUniqueViolation",
			contentType: contentTypeJSON,
			userID:      "31",
			body: map[string]string{
				"url": "https://ya.ru",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", repositories.ErrUniqueViolation).MaxTimes(1)
			},
			want: want{
				contentType: contentTypeJSON,
				statusCode:  409,
				body:        "123",
			},
		},
		{
			name:        "err Validation",
			contentType: contentTypeJSON,
			userID:      "31",
			body: map[string]string{
				"url": "asdasd",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", nil).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "123",
			},
		},
		{
			name:        "err noAuth",
			userID:      "",
			contentType: contentTypeJSON,
			body: map[string]string{
				"url": "https://ya.ru",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", repositories.ErrUniqueViolation).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  401,
				body:        "123",
			},
		},
		{
			name:        "custome service err",
			contentType: contentTypeJSON,
			userID:      "321",
			body: map[string]string{
				"url": "https://ya.ru",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", errors.New("custome err")).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "неверный запрос\n",
			},
		},
		{
			name:        "Good test",
			userID:      "321",
			contentType: contentTypeJSON,
			body: map[string]string{
				"url": "https://ya.ru",
			},
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("1", nil).MaxTimes(1)
			},
			want: want{
				contentType: contentTypeJSON,
				statusCode:  201,
				body:        "1",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			mockServices := mock_handlers.NewMockService(c)
			tt.mockBehavior(mockServices, tt.userID, tt.body["url"])
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			// Create Request
			data, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			if tt.userID != "" {
				request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			}
			request.Header.Set(contentTypeHeader, tt.contentType)
			w := httptest.NewRecorder()
			handlers.CreateShortURLByJSON(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get(contentTypeHeader))
		})
	}
}

func TestHandler_BatchCreateURLByJSON(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, userID string)

	type want struct {
		statusCode int
	}

	tests := []struct {
		name         string
		body         []map[string]string
		userID       string
		mockBehavior mockBehavior
		contentType  string
		want         want
	}{
		{
			name:        "Wrong Content Type",
			contentType: "text/plain",
			userID:      "123",
			body: []map[string]string{
				{
					"url":            "https://ya.ru",
					"correlation_id": "vfwt4312",
				},
				{
					"url":            "https://yand.ru",
					"correlation_id": "fwef13",
				},
			},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().BatchCreate(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return([]models.BatchCreateResponse{}, nil).MaxTimes(1)
			},
			want: want{
				statusCode: 400,
			},
		},
		{
			name:        "good test",
			contentType: contentTypeJSON,
			userID:      "123",
			body: []map[string]string{
				{
					"original_url":   "https://ya.ru",
					"correlation_id": "vfwt4312",
				},
				{
					"original_url":   "https://yand.ru",
					"correlation_id": "fwef13",
				},
			},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().BatchCreate(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return([]models.BatchCreateResponse{}, nil).MaxTimes(1)
			},
			want: want{
				statusCode: 201,
			},
		},
		{
			name:        "err noAuth",
			userID:      "",
			contentType: contentTypeJSON,
			body: []map[string]string{
				{
					"original_url":   "https://ya.ru",
					"correlation_id": "vfwt4312",
				},
				{
					"original_url":   "https://yand.ru",
					"correlation_id": "fwef13",
				},
			},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().BatchCreate(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return([]models.BatchCreateResponse{}, errors.New("any err")).MaxTimes(1)
			},
			want: want{
				statusCode: 401,
			},
		},
		{
			name:        "good test",
			contentType: contentTypeJSON,
			userID:      "123",
			body: []map[string]string{
				{
					"original_url":   "https://ya.ru",
					"correlation_id": "vfwt4312",
				},
				{
					"original_url":   "https://yand.ru",
					"correlation_id": "fwef13",
				},
			},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().BatchCreate(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return([]models.BatchCreateResponse{}, errors.New("any err")).MaxTimes(1)
			},
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			mockServices := mock_handlers.NewMockService(c)
			tt.mockBehavior(mockServices, tt.userID)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			data, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(data))
			if tt.userID != "" {
				request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			}
			request.Header.Set(contentTypeHeader, tt.contentType)
			w := httptest.NewRecorder()
			handlers.BatchCreateURLByJSON(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_GetUserURLs(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, userID string)

	type want struct {
		statusCode int
	}

	tests := []struct {
		name         string
		userID       string
		mockBehavior mockBehavior
		want         want
	}{
		{
			name:   "no auth",
			userID: "",
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().GetListByUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID).Return([]models.GetURLResponse{}, nil).MaxTimes(1)
			},
			want: want{
				statusCode: 401,
			},
		},
		{
			name:   "good test",
			userID: "123",
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().GetListByUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID).Return([]models.GetURLResponse{}, nil).MaxTimes(1)
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name:   "err ErrNoRows",
			userID: "123",
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().GetListByUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID).Return([]models.GetURLResponse{}, sql.ErrNoRows).MaxTimes(1)
			},
			want: want{
				statusCode: 204,
			},
		},
		{
			name:   "err ErrDeleted",
			userID: "123",
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().GetListByUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID).Return([]models.GetURLResponse{}, repositories.ErrDeleted).MaxTimes(1)
			},
			want: want{
				statusCode: 410,
			},
		},
		{
			name:   "err custome",
			userID: "123",
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().GetListByUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID).Return([]models.GetURLResponse{}, errors.New("new")).MaxTimes(1)
			},
			want: want{
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			mockServices := mock_handlers.NewMockService(c)
			tt.mockBehavior(mockServices, tt.userID)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/user/urls", nil)
			if tt.userID != "" {
				request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			}
			w := httptest.NewRecorder()
			handlers.GetUserURLs(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_DeleteUserURLsBatch(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, userID string)

	type want struct {
		statusCode int
	}

	tests := []struct {
		name         string
		userID       string
		body         []string
		mockBehavior mockBehavior
		want         want
	}{
		{
			name:   "no auth",
			userID: "",
			body:   []string{"fafs", "fafadfd"},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().DeleteUserURLsBatch(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return(nil).MaxTimes(1)
			},
			want: want{
				statusCode: 401,
			},
		},
		{
			name:   "good test",
			userID: "123",
			body:   []string{"fafs", "fafadfd"},
			mockBehavior: func(r *mock_handlers.MockService, userID string) {
				r.EXPECT().DeleteUserURLsBatch(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, gomock.Any()).Return(nil).MaxTimes(1)
			},
			want: want{
				statusCode: 202,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			mockServices := mock_handlers.NewMockService(c)
			tt.mockBehavior(mockServices, tt.userID)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			data, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/api/shorten/user/urls", bytes.NewReader(data))
			if tt.userID != "" {
				request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			}
			w := httptest.NewRecorder()
			handlers.DeleteUserURLsBatch(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
