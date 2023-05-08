package handlers

import (
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	mock_services "github.com/ProvoloneStein/go-url-shortener-service/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_createShortURL(t *testing.T) {

	// Init Test Table
	type mockBehavior func(r *mock_services.MockShortener, fullURL string)

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
			mockBehavior: func(r *mock_services.MockShortener, fullURL string) {
				r.EXPECT().CreateShortURL(fullURL).Return("123", nil).AnyTimes()
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный header\n",
			},
		},
		{
			name:        "Good test",
			contentType: "text/plain",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_services.MockShortener, fullURL string) {
				r.EXPECT().CreateShortURL(fullURL).Return("1", nil)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				body:        "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_services.NewMockShortener(c)
			tt.mockBehavior(repo, tt.body)
			services := &services.Service{Shortener: repo}
			handlers := Handler{services}

			// Create Request
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			handlers.createShortURL(w, request)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(respBody))
		})
	}
}

func TestHandler_getByShort(t *testing.T) {

	// Init Test Table
	type mockBehavior func(r *mock_services.MockShortener, fullURL string)

	type want struct {
		statusCode  int
		contentType string
		location    string
		body        string
	}

	tests := []struct {
		name         string
		id           string
		mockBehavior mockBehavior
		want         want
	}{
		{
			name: "Good test",
			id:   "gwrags",
			mockBehavior: func(r *mock_services.MockShortener, shortURL string) {
				r.EXPECT().GetFullByID(shortURL).Return("https://ya.ru", nil)
			},
			want: want{
				statusCode: 307,
				location:   "https://ya.ru",
			},
		},
		{
			name: "Wrong id",
			id:   "adsga",
			mockBehavior: func(r *mock_services.MockShortener, shortURL string) {
				r.EXPECT().GetFullByID(shortURL).Return("", fmt.Errorf("Ошибочка")).AnyTimes()
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный запрос\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_services.NewMockShortener(c)
			tt.mockBehavior(repo, tt.id)
			services := &services.Service{Shortener: repo}
			handlers := Handler{services}

			// Create Request
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			w := httptest.NewRecorder()
			handlers.getByShort(w, request)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, string(respBody))
		})
	}
}