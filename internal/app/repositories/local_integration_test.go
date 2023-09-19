package repositories

import (
	"context"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLocalRepository_Create(t *testing.T) {
	type want struct {
		err error
		res string
	}

	tests := []struct {
		name     string
		ctx      context.Context
		store    map[string][3]string
		userID   string
		shortURL string
		fullURL  string
		want     want
	}{
		{
			name:     "ShortURLExists test",
			ctx:      context.Background(),
			userID:   "341",
			shortURL: "faddsfa",
			fullURL:  "dasfaf",
			store: map[string][3]string{
				"faddsfa": {"dasfaf", "341", "f"},
			},
			want: want{
				res: "",
				err: ErrShortURLExists,
			},
		},
		{
			name:     "ErrUniqueViolation test",
			ctx:      context.Background(),
			userID:   "341",
			shortURL: "faddsfa",
			fullURL:  "dasfaf",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: "31fd",
				err: ErrUniqueViolation,
			},
		},
		{
			name:     "good test",
			ctx:      context.Background(),
			userID:   "341",
			shortURL: "faddsfa",
			fullURL:  "dasfaf",
			store:    map[string][3]string{},
			want: want{
				res: "faddsfa",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			res, err := repository.Create(tt.ctx, tt.userID, tt.fullURL, tt.shortURL)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_BatchCreate(t *testing.T) {
	type want struct {
		err error
		res []models.BatchCreateResponse
	}

	tests := []struct {
		name  string
		ctx   context.Context
		store map[string][3]string
		data  []models.BatchCreateData
		want  want
	}{
		{
			name: "ErrUniqueViolation test",
			ctx:  context.Background(),
			data: []models.BatchCreateData{
				{
					URL:      "dasfaf",
					UUID:     "usdsdu",
					ShortURL: "31fd",
					UserID:   "231",
				},
				{
					URL:      "fds",
					UUID:     "usd3sdu",
					ShortURL: "211343",
					UserID:   "231",
				},
			},
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: "31fd", UUID: "usdsdu"}},
				err: ErrShortURLExists,
			},
		},
		{
			name: "good test",
			ctx:  context.Background(),
			data: []models.BatchCreateData{
				{
					URL:      "dasfaf",
					UUID:     "usdsdu",
					ShortURL: "213",
					UserID:   "231",
				},
				{
					URL:      "fds",
					UUID:     "usd3sdu",
					ShortURL: "211343",
					UserID:   "231",
				},
			},
			store: map[string][3]string{},
			want: want{
				res: []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: "test/213", UUID: "usdsdu"}, models.BatchCreateResponse{ShortURL: "test/211343", UUID: "usd3sdu"}},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			res, err := repository.BatchCreate(tt.ctx, tt.data)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_GetByShort(t *testing.T) {
	type want struct {
		err error
		res string
	}

	tests := []struct {
		name     string
		ctx      context.Context
		store    map[string][3]string
		shortURL string
		want     want
	}{
		{
			name:     "ErrUniqueViolation test",
			ctx:      context.Background(),
			shortURL: "31fd",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: "dasfaf",
				err: nil,
			},
		},
		{
			name:     "ErrDeleted test",
			ctx:      context.Background(),
			shortURL: "31fd",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "t"},
			},
			want: want{
				res: "",
				err: ErrDeleted,
			},
		},
		{
			name:     "errWithVal test",
			ctx:      context.Background(),
			shortURL: "341",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: "",
				err: ErrURLNotFound,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			res, err := repository.GetByShort(tt.ctx, tt.shortURL)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_Ping(t *testing.T) {
	type want struct {
		err error
	}

	tests := []struct {
		want  want
		store map[string][3]string
		ctx   context.Context
		name  string
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			err := repository.Ping()
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_Close(t *testing.T) {
	type want struct {
		err error
	}

	tests := []struct {
		want  want
		store map[string][3]string
		ctx   context.Context
		name  string
	}{
		{
			name: "ok test",
			ctx:  context.Background(),
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			err := repository.Close()
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_GetListByUser(t *testing.T) {
	type want struct {
		err error
		res []models.GetURLResponse
	}

	tests := []struct {
		ctx    context.Context
		name   string
		store  map[string][3]string
		userID string
		want   want
	}{
		{
			name:   "ok test",
			ctx:    context.Background(),
			userID: "341",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: []models.GetURLResponse{
					{
						ShortURL: "31fd",
						URL:      "dasfaf",
					},
				},
				err: nil,
			},
		},
		{
			name:   "empty test",
			ctx:    context.Background(),
			userID: "433",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			want: want{
				res: nil,
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			res, err := repository.GetListByUser(tt.ctx, tt.userID)
			assert.Equal(t, tt.want.res, res)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestLocalRepository_DeleteUserURLsBatch(t *testing.T) {
	type want struct {
		err error
	}

	tests := []struct {
		want   want
		name   string
		ctx    context.Context
		store  map[string][3]string
		userID string
		data   []string
	}{
		{
			name:   "ok test",
			ctx:    context.Background(),
			userID: "341",
			store: map[string][3]string{
				"31fd": {"dasfaf", "341", "f"},
			},
			data: []string{
				"31fd",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			cfg := configs.AppConfig{BaseURL: "test", Addr: "test"}
			repository := LocalRepository{logger: zap.NewNop(), cfg: cfg, store: tt.store}

			err := repository.DeleteUserURLsBatch(tt.ctx, tt.userID, tt.data)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}
