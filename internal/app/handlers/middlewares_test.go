package handlers

import (
	"context"
	"errors"
	mock_handlers "github.com/ProvoloneStein/go-url-shortener-service/internal/app/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	// Init Test Table
	type mockGenerateTokenBehavior func(r *mock_handlers.MockService)
	type mockParseTokenBehavior func(r *mock_handlers.MockService)

	type want struct {
		statusCode int
	}

	tests := []struct {
		mockGenerateToken mockGenerateTokenBehavior
		mockParseToken    mockParseTokenBehavior
		name              string
		cookieVal         string
		want              want
	}{
		{
			name:      "ok test",
			cookieVal: "",
			mockGenerateToken: func(r *mock_handlers.MockService) {
				r.EXPECT().GenerateToken(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())).
					Return("123", nil).MaxTimes(1)
			},
			mockParseToken: func(r *mock_handlers.MockService) {
				r.EXPECT().ParseToken(gomock.AssignableToTypeOf("string")).Return("123", nil).MaxTimes(1)
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name:      "good test parse token",
			cookieVal: "any_token_value",
			mockGenerateToken: func(r *mock_handlers.MockService) {
				r.EXPECT().GenerateToken(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())).
					Return("123", nil).MaxTimes(1)
			},
			mockParseToken: func(r *mock_handlers.MockService) {
				r.EXPECT().ParseToken(gomock.AssignableToTypeOf("string")).Return("123", nil).MaxTimes(1)
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name:      "GenerateToken err",
			cookieVal: "",
			mockGenerateToken: func(r *mock_handlers.MockService) {
				r.EXPECT().GenerateToken(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())).
					Return("123", errors.New("any err")).MaxTimes(1)
			},
			mockParseToken: func(r *mock_handlers.MockService) {
				r.EXPECT().ParseToken(gomock.AssignableToTypeOf("string")).Return("123", nil).MaxTimes(1)
			},
			want: want{
				statusCode: 500,
			},
		},
		{
			name:      "parsetoken err",
			cookieVal: "any_token_value",
			mockGenerateToken: func(r *mock_handlers.MockService) {
				r.EXPECT().GenerateToken(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())).
					Return("123", nil).MaxTimes(1)
			},
			mockParseToken: func(r *mock_handlers.MockService) {
				r.EXPECT().ParseToken(gomock.AssignableToTypeOf("string")).Return("123", errors.New("any err")).MaxTimes(1)
			},
			want: want{
				statusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// create a handler to use as "next" which will verify the request
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			mockServices := mock_handlers.NewMockService(c)
			tt.mockGenerateToken(mockServices)
			tt.mockParseToken(mockServices)

			middleware := userIdentity(mockServices, zap.NewNop())
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			if tt.cookieVal != "" {
				request.AddCookie(&http.Cookie{Name: "authToken", Value: tt.cookieVal})
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)
			result := w.Result()
			defer func() {
				deferErr := result.Body.Close()
				assert.NoError(t, deferErr)
			}()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_gzipReadWriterHandler(t *testing.T) {

	type want struct {
		contentEncoding string
		statusCode      int
	}

	tests := []struct {
		name            string
		acceptEncoding  string
		contentEncoding string
		want            want
	}{
		{
			name:            "good test",
			acceptEncoding:  "gzip",
			contentEncoding: "gzip",
			want: want{
				statusCode:      200,
				contentEncoding: "gzip",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// create a handler to use as "next" which will verify the request
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			middleware := gzipReadWriterHandler(zap.NewNop())
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			request.Header.Add("Accept-Encoding", tt.acceptEncoding)
			request.Header.Add("Content-Encoding", tt.contentEncoding)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentEncoding, result.Header.Get("Content-Encoding"))
		})
	}
}
