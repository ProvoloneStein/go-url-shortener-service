package repositories

import (
	"context"
	"database/sql"
	"errors"
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
		mockBehavior mockBehavior
		name         string
		shortURL     string
		args         args
		wantErr      bool
		ctxEnd       bool
		id           int
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
			} else {
				defer cancel()
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
		mockBehavior mockBehavior
		name         string
		args         args
		wantErr      bool
		ctxEnd       bool
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
			} else {
				defer cancel()
			}
			testCase.mockBehavior(testCase.args)
			_, err := repo.BatchCreate(ctx, testCase.args.data)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

			}

		})
	}
}

func TestDBRepository_GetByShort(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	type mockBehavior func(fullURL, shortURL string)

	tests := []struct {
		mockBehavior mockBehavior
		name         string
		fullURL      string
		shortURL     string
		wantErr      bool
		ctxEnd       bool
	}{
		{
			name:     "ok",
			ctxEnd:   false,
			fullURL:  "url",
			shortURL: "dsds",
			wantErr:  false,
			mockBehavior: func(fullURL, shortURL string) {
				selectMockRows := sqlmock.NewRows([]string{"shorten", "deleted"}).
					AddRow(fullURL, false)
				mock.ExpectQuery("SELECT url, deleted FROM shortener").
					WithArgs(shortURL).WillReturnRows(selectMockRows)
			},
		},
		{
			name:    "err URLNotFound",
			ctxEnd:  false,
			wantErr: true,
			mockBehavior: func(fullURL, shortURL string) {
				mock.ExpectQuery("SELECT url, deleted FROM shortener").
					WithArgs(shortURL).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:    "err select",
			ctxEnd:  false,
			wantErr: true,
			mockBehavior: func(fullURL, shortURL string) {
				selectMockRows := sqlmock.NewRows([]string{"shorten", "deleted"}).
					AddRow(fullURL, true)
				mock.ExpectQuery("SELECT url, deleted FROM shortener").
					WithArgs(shortURL).WillReturnRows(selectMockRows)
			},
		},
		{
			name:    "err delete",
			ctxEnd:  false,
			wantErr: true,
			mockBehavior: func(fullURL, shortURL string) {
				mock.ExpectQuery("SELECT url, deleted FROM shortener").
					WithArgs(shortURL).WillReturnError(errors.New("sdsd"))
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			if testCase.ctxEnd {
				cancel()
			} else {
				defer cancel()
			}
			testCase.mockBehavior(testCase.fullURL, testCase.shortURL)
			got, err := repo.GetByShort(ctx, testCase.shortURL)
			if testCase.wantErr {
				assert.Equal(t, "", got)
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.fullURL, got)
				assert.NoError(t, err)

			}

		})
	}
}

func TestDBRepository_Ping(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	type mockBehavior func()

	tests := []struct {
		mockBehavior mockBehavior
		name         string
		wantErr      bool
	}{
		{
			name:    "ok",
			wantErr: false,
			mockBehavior: func() {
				mock.ExpectPing().WillReturnError(nil)
			},
		},
		{
			name:    "err",
			wantErr: true,
			mockBehavior: func() {
				mock.ExpectPing().WillReturnError(errors.New("any err"))
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()
			err := repo.Ping()
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}

func TestDBRepository_Close(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "Error - %v", err)
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	repo := DBRepository{
		logger: zap.NewNop(),
		cfg:    configs.AppConfig{},
		db:     db,
	}

	mock.ExpectClose().WillReturnError(errors.New("dsds"))
	assert.Error(t, repo.Close())

	mock.ExpectClose().WillReturnError(nil)
	assert.NoError(t, repo.Close())
}

// TODO: не получается сделать мок - пробелма с аргументом data

// type CustomConverter struct{}
//
// func (s CustomConverter) ConvertValue(v interface{}) (driver.Value, error) {
//	switch v.(type) {
//	case string:
//		return v.(string), nil
//	case []string:
//		res := pgtype.Array[string]{
//			Elements: v.([]string),
//			Dims:     []pgtype.ArrayDimension{{Length: int32(len(v.([]string))), LowerBound: 1}},
//			Valid:    true,
//		}
//		return res, nil
//	case int:
//		return v.(int), nil
//	default:
//		return nil, errors.New(fmt.Sprintf("cannot convert %T with value %v", v, v))
//	}
//}
//
// func TestDBRepository_DeleteUserURLsBatch(t *testing.T) {
//	mockDB, mock, err := sqlmock.New(sqlmock.ValueConverterOption(CustomConverter{}))
//	//mockDB, mock, err := sqlmock.New()
//	require.NoError(t, err, "Error - %v", err)
//	db := sqlx.NewDb(mockDB, "sqlmock")
//	defer db.Close()
//
//	repo := DBRepository{
//		logger: zap.NewNop(),
//		cfg:    configs.AppConfig{},
//		db:     db,
//	}
//
//	type args struct {
//		userID string
//		data   []string
//	}
//
//	type mockBehavior func(args)
//
//	tests := []struct {
//		name         string
//		wantErr      bool
//		args         args
//		mockBehavior mockBehavior
//	}{
//		{
//			name:    "ok",
//			wantErr: false,
//			args: args{
//				userID: "2",
//				data:   []string{"1", "2"},
//			},
//			mockBehavior: func(args args) {
//				mock.ExpectExec("UPDATE shortener set").WithArgs(args.userID, args.data).WillReturnResult(nil)
//			},
//		},
//		{
//			name:    "err",
//			wantErr: true,
//			args: args{
//				userID: "1",
//				data:   []string{"1", "2"},
//			},
//			mockBehavior: func(args args) {
//				mock.ExpectExec("UPDATE shortener set deleted").WithArgs(args.userID, args.data).WillReturnError(errors.New("any"))
//			},
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mockBehavior(testCase.args)
//			err := repo.DeleteUserURLsBatch(context.Background(), testCase.args.userID, testCase.args.data)
//			if testCase.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//
//			}
//		})
//	}
//}

func TestDBRepository_GetListByUser(t *testing.T) {
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
		userID string
	}

	type mockBehavior func(args)

	tests := []struct {
		args         args
		mockBehavior mockBehavior
		name         string
		ctxEnd       bool
		wantErr      bool
	}{
		{
			name:    "ok",
			ctxEnd:  false,
			wantErr: false,
			args: args{
				userID: "2",
			},
			mockBehavior: func(args args) {
				selectMockRows := sqlmock.NewRows([]string{"url", "shorten"}).
					AddRow("12", "23")
				mock.ExpectQuery("SELECT url, shorten FROM").WithArgs(args.userID).WillReturnRows(selectMockRows)
			},
		},
		{
			name:    "err",
			ctxEnd:  false,
			wantErr: true,
			args: args{
				userID: "2",
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT url, shorten FROM").WithArgs(args.userID).WillReturnError(errors.New("sd"))
			},
		},
		{
			name:    "ctx err",
			ctxEnd:  true,
			wantErr: true,
			args: args{
				userID: "2",
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT url, shorten FROM").WithArgs(args.userID).WillReturnError(errors.New("sd"))
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)
			ctx, cancel := context.WithCancel(context.Background())

			if testCase.ctxEnd {
				cancel()
			} else {
				defer cancel()
			}
			_, err := repo.GetListByUser(ctx, testCase.args.userID)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}
