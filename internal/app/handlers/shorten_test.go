package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mock_handlers "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers/mocks"
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
				contentType: contentTypeJSON,
				statusCode:  409,
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
			request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			request.Header.Set(contenntTypeHeader, tt.contentType)
			w := httptest.NewRecorder()
			handlers.CreateShortURLByJSON(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get(contenntTypeHeader))
		})
	}
}
