package services

import (
	"context"
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"
	mock_services "github.com/ProvoloneStein/go-url-shortener-service/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestService_CreateShortURL(t *testing.T) {

	type mockRepo func(r *mock_services.MockRepository, ctx context.Context, userID, fullURL string)

	type want struct {
		res string
		err error
	}

	tests := []struct {
		name     string
		ctx      context.Context
		userID   string
		fullURL  string
		repoFunc mockRepo
		want     want
	}{
		{
			name:    "ok",
			ctx:     context.Background(),
			userID:  "123",
			fullURL: "sdsds",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID, fullURL string) {
				r.EXPECT().Create(ctx, userID, fullURL, gomock.Any()).Return("test", nil).MaxTimes(1)
			},
			want: want{
				res: "test",
				err: nil,
			},
		},
		{
			name:    "ErrUniqueViolation",
			ctx:     context.Background(),
			userID:  "123",
			fullURL: "sdsds",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID, fullURL string) {
				r.EXPECT().Create(ctx, userID, fullURL, gomock.Any()).Return("", repositories.ErrUniqueViolation).MaxTimes(1)
			},
			want: want{
				res: "",
				err: repositories.ErrUniqueViolation,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			repo := mock_services.NewMockRepository(c)
			tt.repoFunc(repo, tt.ctx, tt.userID, tt.fullURL)
			defer c.Finish()
			service := Service{zap.NewNop(), repo, configs.AppConfig{}}
			res, err := service.CreateShortURL(tt.ctx, tt.userID, tt.fullURL)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestService_BatchCreate(t *testing.T) {

	type mockRepo func(r *mock_services.MockRepository, ctx context.Context, shortURL string)

	type want struct {
		res []models.BatchCreateResponse
		err error
	}

	tests := []struct {
		name     string
		ctx      context.Context
		userID   string
		data     []models.BatchCreateRequest
		repoFunc mockRepo
		want     want
	}{
		{
			name:   "ok",
			ctx:    context.Background(),
			userID: "123",
			data:   []models.BatchCreateRequest{models.BatchCreateRequest{URL: "dsadas", UUID: "dsadas"}},
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, shortURL string) {
				r.EXPECT().BatchCreate(ctx, gomock.Any()).Return([]models.BatchCreateResponse{}, nil).MaxTimes(1)
			},
			want: want{
				res: []models.BatchCreateResponse{},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			repo := mock_services.NewMockRepository(c)
			tt.repoFunc(repo, tt.ctx, tt.userID)
			defer c.Finish()
			service := Service{zap.NewNop(), repo, configs.AppConfig{}}
			res, err := service.BatchCreate(tt.ctx, tt.userID, tt.data)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, tt.want.err, err)
		})
	}
}

func TestService_GetFullByID(t *testing.T) {

	type mockRepo func(r *mock_services.MockRepository, ctx context.Context, shortURL string)

	type want struct {
		res   string
		isErr bool
	}

	tests := []struct {
		name     string
		ctx      context.Context
		shortURL string
		repoFunc mockRepo
		want     want
	}{
		{
			name:     "ok",
			ctx:      context.Background(),
			shortURL: "123",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, shortURL string) {
				r.EXPECT().GetByShort(ctx, shortURL).Return("321", nil).MaxTimes(1)
			},
			want: want{
				res:   "321",
				isErr: false,
			},
		},
		{
			name:     "err",
			ctx:      context.Background(),
			shortURL: "123",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, shortURL string) {
				r.EXPECT().GetByShort(ctx, shortURL).Return("", errors.New("any")).MaxTimes(1)
			},
			want: want{
				res:   "",
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			repo := mock_services.NewMockRepository(c)
			tt.repoFunc(repo, tt.ctx, tt.shortURL)
			defer c.Finish()
			service := Service{zap.NewNop(), repo, configs.AppConfig{}}
			res, err := service.GetFullByID(tt.ctx, tt.shortURL)
			assert.Equal(t, tt.want.res, res)
			if tt.want.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetListByUser(t *testing.T) {

	type mockRepo func(r *mock_services.MockRepository, ctx context.Context, userID string)

	type want struct {
		res   []models.GetURLResponse
		isErr bool
	}

	tests := []struct {
		name     string
		ctx      context.Context
		userID   string
		repoFunc mockRepo
		want     want
	}{
		{
			name:   "ok",
			ctx:    context.Background(),
			userID: "123",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID string) {
				r.EXPECT().GetListByUser(ctx, userID).Return([]models.GetURLResponse{models.GetURLResponse{
					ShortURL: "sdsd", URL: "sdds"}}, nil).MaxTimes(1)
			},
			want: want{
				res: []models.GetURLResponse{models.GetURLResponse{
					ShortURL: "sdsd", URL: "sdds"}},
				isErr: false,
			},
		},
		{
			name:   "err",
			ctx:    context.Background(),
			userID: "123",
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID string) {
				r.EXPECT().GetListByUser(ctx, userID).Return([]models.GetURLResponse{models.GetURLResponse{
					ShortURL: "sdsd", URL: "sdds"}}, errors.New("err")).MaxTimes(1)
			},
			want: want{
				res: []models.GetURLResponse{models.GetURLResponse{
					ShortURL: "sdsd", URL: "sdds"}},
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			repo := mock_services.NewMockRepository(c)
			tt.repoFunc(repo, tt.ctx, tt.userID)
			defer c.Finish()
			service := Service{zap.NewNop(), repo, configs.AppConfig{}}
			res, err := service.GetListByUser(tt.ctx, tt.userID)
			assert.Equal(t, tt.want.res, res)
			if tt.want.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteUserURLsBatch(t *testing.T) {

	type mockRepo func(r *mock_services.MockRepository, ctx context.Context, userID string, data []string)

	type want struct {
		isErr bool
	}

	tests := []struct {
		name     string
		ctx      context.Context
		userID   string
		data     []string
		repoFunc mockRepo
		want     want
	}{
		{
			name:   "ok",
			ctx:    context.Background(),
			userID: "123",
			data:   []string{"dfdf", "fdaf"},
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID string, data []string) {
				r.EXPECT().DeleteUserURLsBatch(ctx, userID, data).Return(nil).MaxTimes(1)
			},
			want: want{
				isErr: false,
			},
		},
		{
			name:   "err",
			ctx:    context.Background(),
			userID: "123",
			data:   []string{"dfdf", "fdaf"},
			repoFunc: func(r *mock_services.MockRepository, ctx context.Context, userID string, data []string) {
				r.EXPECT().DeleteUserURLsBatch(ctx, userID, data).Return(errors.New("err")).MaxTimes(1)
			},
			want: want{
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			repo := mock_services.NewMockRepository(c)
			tt.repoFunc(repo, tt.ctx, tt.userID, tt.data)
			defer c.Finish()
			service := Service{zap.NewNop(), repo, configs.AppConfig{}}
			err := service.DeleteUserURLsBatch(tt.ctx, tt.userID, tt.data)
			if tt.want.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Ping(t *testing.T) {

	type mockPing func(r *mock_services.MockRepository)

	type want struct {
		isErr bool
	}

	tests := []struct {
		name string
		ping mockPing
		want want
	}{
		{
			name: "ok",
			ping: func(r *mock_services.MockRepository) {
				r.EXPECT().Ping().Return(nil).MaxTimes(1)
			},
			want: want{
				isErr: false,
			},
		},
		{
			name: "err",
			ping: func(r *mock_services.MockRepository) {
				r.EXPECT().Ping().Return(errors.New("sd")).MaxTimes(1)
			},
			want: want{
				isErr: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			mockRepo := mock_services.NewMockRepository(c)
			tt.ping(mockRepo)
			defer c.Finish()
			service := Service{zap.NewNop(), mockRepo, configs.AppConfig{}}
			err := service.Ping()
			if tt.want.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
