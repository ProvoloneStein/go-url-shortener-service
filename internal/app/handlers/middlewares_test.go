package handlers

import (
	"context"
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
	type mockPingBehavior func(r *mock_handlers.MockService)

	type want struct {
		statusCode int
	}

	tests := []struct {
		name              string
		ctx               context.Context
		mockGenerateToken mockGenerateTokenBehavior
		mockParseToken    mockParseTokenBehavior
		mockPing          mockPingBehavior
		want              want
	}{
		{
			name: "good test",
			ctx:  context.Background(),
			mockGenerateToken: func(r *mock_handlers.MockService) {
				r.EXPECT().GenerateToken(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())).
					Return("123", nil).MaxTimes(1)
			},
			mockParseToken: func(r *mock_handlers.MockService) {
				r.EXPECT().ParseToken(gomock.AssignableToTypeOf("string")).Return("123", nil).MaxTimes(1)
			},
			mockPing: func(r *mock_handlers.MockService) {
				r.EXPECT().Ping().Return(nil).MaxTimes(1)
			},
			want: want{
				statusCode: 200,
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
			tt.mockGenerateToken(mockServices)
			tt.mockParseToken(mockServices)
			tt.mockPing(mockServices)

			handlers := Handler{logger: zap.NewNop(), services: mockServices}

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			request = request.WithContext(tt.ctx)
			w := httptest.NewRecorder()
			handlers.pingDB(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
