package services

import (
	"context"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	mock_services "github.com/ProvoloneStein/go-url-shortener-service/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestHandler_GenerateToken(t *testing.T) {

	type want struct {
		res string
		err error
	}

	tests := []struct {
		name string
		ctx  context.Context
		want want
	}{
		{
			name: "ok",
			ctx:  context.Background(),
			want: want{
				res: "asdf",
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			mockRepo := mock_services.NewMockRepository(c)
			defer c.Finish()
			service := Service{zap.NewNop(), mockRepo, configs.AppConfig{}}
			res, err := service.GenerateToken(tt.ctx)
			assert.ErrorIs(t, tt.want.err, err)
			if tt.want.res == "" {
				assert.Equal(t, res, tt.want.res)
			}
		})
	}
}

func TestHandler_ParseToken(t *testing.T) {

	type want struct {
		res   string
		isErr bool
	}

	tests := []struct {
		name  string
		token string
		want  want
	}{
		{
			name:  "ok",
			token: "",
			want: want{
				res:   "asdf",
				isErr: false,
			},
		},
		{
			name:  "invalid token",
			token: "token",
			want: want{
				res:   "asdf",
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var err error
			// Init Dependencies
			c := gomock.NewController(t)
			mockRepo := mock_services.NewMockRepository(c)
			defer c.Finish()
			service := Service{zap.NewNop(), mockRepo, configs.AppConfig{}}
			if tt.token == "" {
				tt.token, err = service.GenerateToken(context.Background())
				assert.NoError(t, err)
			}
			res, err := service.ParseToken(tt.token)
			if tt.want.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want.res == "" {
				assert.Equal(t, res, tt.want.res)
			}
		})
	}
}
