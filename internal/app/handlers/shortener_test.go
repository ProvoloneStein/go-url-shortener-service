package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	mock_handlers "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers/mocks"
)

func TestHandler_createShortURL(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, fullURL string)

	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name         string
		body         string
		mockBehavior mockBehavior
		contentType  string
		want         want
	}{
		{
			name:        "Wrong Content Type",
			contentType: "type",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.AssignableToTypeOf("string"), fullURL).Return("123", nil).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный header запроса\n",
			},
		},
		{
			name:        "Good test",
			contentType: "text/plain",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.AssignableToTypeOf("string"), fullURL).Return("1", nil).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
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
			tt.mockBehavior(mockServices, tt.body)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			request = request.WithContext(context.WithValue(request.Context(), userCtx, "12345"))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			handlers.createShortURL(w, request)
			result := w.Result()
			defer result.Body.Close()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(respBody))
		})
	}
}

func TestHandler_getByShort(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, shortURL string)

	type want struct {
		statusCode  int
		contentType string
		location    string
		body        string
	}

	tests := []struct {
		name         string
		id           string
		url          string
		mockBehavior mockBehavior
		want         want
	}{
		{
			name: "Good test",
			id:   "gwrags",
			url:  "http://localhost:8080/gwrags",
			mockBehavior: func(r *mock_handlers.MockService, shortURL string) {
				r.EXPECT().GetFullByID(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), shortURL).Return("https://ya.ru", nil).MaxTimes(1)
			},
			want: want{
				statusCode: 307,
				location:   "https://ya.ru",
			},
		},
		{
			name: "Wrong id",
			id:   "adsga",
			url:  "http://localhost:8080/adsga",
			mockBehavior: func(r *mock_handlers.MockService, shortURL string) {
				r.EXPECT().GetFullByID(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), shortURL).Return("", fmt.Errorf("Ошибочка")).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный запрос\n",
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
			tt.mockBehavior(mockServices, tt.id)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			router := chi.NewRouteContext()

			router.URLParams.Add("id", tt.id)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, router))
			request = request.WithContext(context.WithValue(request.Context(), userCtx, "12345"))

			handlers.getByShort(w, request)
			result := w.Result()
			defer result.Body.Close()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(respBody))
		})
	}
}
