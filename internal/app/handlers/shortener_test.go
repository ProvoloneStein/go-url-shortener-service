package handlers

import (
	"context"
	"errors"
	mock_handlers "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers/mocks"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHandler_CreateShortURL(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_handlers.MockService, userID, fullURL string)

	type want struct {
		contentType string
		body        string
		statusCode  int
	}

	tests := []struct {
		name         string
		body         string
		userID       string
		mockBehavior mockBehavior
		contentType  string
		want         want
	}{
		{
			name:        "Wrong Content Type",
			contentType: "application/json",
			userID:      "123",
			body:        "https://ya.ru",
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
			contentType: "text/plain",
			userID:      "31",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", repositories.ErrUniqueViolation).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  409,
				body:        "123",
			},
		},
		{
			name:        "err noAuth",
			userID:      "",
			contentType: "text/plain",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", repositories.ErrUniqueViolation).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  401,
				body:        "Unauthorized\n",
			},
		},
		{
			name:        "custome service err",
			contentType: "text/plain",
			userID:      "321",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("123", errors.New("custome err")).MaxTimes(1)
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Неверный запрос\n",
			},
		},
		{
			name:        "Good test",
			contentType: "text/plain",
			userID:      "321",
			body:        "https://ya.ru",
			mockBehavior: func(r *mock_handlers.MockService, userID, fullURL string) {
				r.EXPECT().CreateShortURL(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					userID, fullURL).Return("1", nil).MaxTimes(1)
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
			tt.mockBehavior(mockServices, tt.userID, tt.body)
			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			// Create Request
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.userID != "" {
				request = request.WithContext(context.WithValue(request.Context(), userCtx, tt.userID))
			}
			request.Header.Set(contenntTypeHeader, tt.contentType)
			w := httptest.NewRecorder()
			handlers.CreateShortURL(w, request)
			result := w.Result()
			defer result.Body.Close()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get(contenntTypeHeader))
			assert.Equal(t, tt.want.body, string(respBody))
		})
	}
}
