package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestDBRepository_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	type args struct {
		userID   string
		fullURL  string
		shortURL string
	}

	type mockBehavior func(args args, shortURL string)

	tests := []struct {
		wantErr      bool
		ctxEnd       bool
		args         args
		id           int
		name         string
		shortURL     string
		mockBehavior mockBehavior
	}{
		{
			name:   "ok",
			ctxEnd: false,
			args: args{
				userID:   "1",
				fullURL:  "12",
				shortURL: "21",
			},
			shortURL: "21",
			wantErr:  false,
			mockBehavior: func(args args, shortURL string) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").
					WithArgs(args.shortURL).WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten"}).
					AddRow(shortURL)

				mock.ExpectQuery("INSERT INTO shortener").
					WithArgs(args.fullURL, args.shortURL, args.userID).
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
		{
			name:   "err commitTx",
			ctxEnd: false,
			args: args{
				userID:   "1",
				fullURL:  "12",
				shortURL: "21",
			},
			shortURL: "21",
			wantErr:  true,
			mockBehavior: func(args args, shortURL string) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").
					WithArgs(args.shortURL).WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten"}).
					AddRow(shortURL)

				mock.ExpectQuery("INSERT INTO shortener").
					WithArgs(args.fullURL, args.shortURL, args.userID).
					WillReturnRows(insertMockRows)
				mock.ExpectCommit().WillReturnError(errors.New("sdsdsd"))
			},
		},
		{
			name:   "err commitTx",
			ctxEnd: false,
			args: args{
				userID:   "1",
				fullURL:  "12",
				shortURL: "21",
			},
			shortURL: "21",
			wantErr:  true,
			mockBehavior: func(args args, shortURL string) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").
					WithArgs(args.shortURL).WillReturnError(errors.New("sdds"))

				mock.ExpectRollback()

				insertMockRows := sqlmock.NewRows([]string{"shorten"}).
					AddRow(shortURL)

				mock.ExpectQuery("INSERT INTO shortener").
					WithArgs(args.fullURL, args.shortURL, args.userID).
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
		{
			name:   "err beginTx",
			ctxEnd: false,
			args: args{
				userID:   "1",
				fullURL:  "12",
				shortURL: "21",
			},
			shortURL: "",
			wantErr:  true,
			mockBehavior: func(args args, shortURL string) {
				mock.ExpectBegin().WillReturnError(errors.New("sd"))

				mock.ExpectQuery("SELECT id FROM shortener").
					WithArgs(args.shortURL).WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten"}).
					AddRow(shortURL)

				mock.ExpectQuery("INSERT INTO shortener").
					WithArgs(args.fullURL, args.shortURL, args.userID).
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
		{
			name:   "context err",
			ctxEnd: true,
			args: args{
				userID:   "1",
				fullURL:  "12",
				shortURL: "21",
			},
			shortURL: "",
			wantErr:  true,
			mockBehavior: func(args args, shortURL string) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").
					WithArgs(args.shortURL).WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten"}).
					AddRow(shortURL)

				mock.ExpectQuery("INSERT INTO shortener").
					WithArgs(args.fullURL, args.shortURL, args.userID).
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			if testCase.ctxEnd {
				cancel()
			}
			testCase.mockBehavior(testCase.args, testCase.shortURL)
			got, err := repo.Create(ctx, testCase.args.userID, testCase.args.fullURL, testCase.args.shortURL)
			if testCase.wantErr {
				assert.Equal(t, "", got)
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.shortURL, got)
				assert.NoError(t, err)

			}

		})
	}
}

func TestDBRepository_BatchCreate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	type args struct {
		data []models.BatchCreateData
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		wantErr      bool
		ctxEnd       bool
		args         args
		mockBehavior mockBehavior
	}{
		{
			name:   "ok",
			ctxEnd: false,
			args: args{
				data: []models.BatchCreateData{models.BatchCreateData{UUID: "sdsd", URL: "dsds", ShortURL: "sdsd", UserID: "ds"}},
			},
			wantErr: false,
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten", "correlation_id"}).
					AddRow("dsds", 1)

				mock.ExpectQuery("INSERT INTO shortener").
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			if testCase.ctxEnd {
				cancel()
			}
			testCase.mockBehavior(testCase.args)
			got, err := repo.BatchCreate(ctx, testCase.args.data)
			fmt.Println(got)
			if testCase.wantErr {
				//assert.Equal(t, "", got)
				assert.Error(t, err)
			} else {
				//assert.Equal(t, testCase.shortURL, got)
				assert.NoError(t, err)

			}

		})
	}
}

func TestDBRepository_BatchDelete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	type args struct {
		data []models.BatchCreateData
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		wantErr      bool
		ctxEnd       bool
		args         args
		mockBehavior mockBehavior
	}{
		{
			name:   "ok",
			ctxEnd: false,
			args: args{
				data: []models.BatchCreateData{models.BatchCreateData{UUID: "sdsd", URL: "dsds", ShortURL: "sdsd", UserID: "ds"}},
			},
			wantErr: false,
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT id FROM shortener").WillReturnError(sql.ErrNoRows)

				insertMockRows := sqlmock.NewRows([]string{"shorten", "correlation_id"}).
					AddRow("dsds", 1)

				mock.ExpectQuery("INSERT INTO shortener").
					WillReturnRows(insertMockRows)
				mock.ExpectCommit()
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			if testCase.ctxEnd {
				cancel()
			}
			testCase.mockBehavior(testCase.args)
			got, err := repo.BatchCreate(ctx, testCase.args.data)
			fmt.Println(got)
			if testCase.wantErr {
				//assert.Equal(t, "", got)
				assert.Error(t, err)
			} else {
				//assert.Equal(t, testCase.shortURL, got)
				assert.NoError(t, err)

			}

		})
	}
}
