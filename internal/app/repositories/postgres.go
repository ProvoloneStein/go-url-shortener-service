package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
)

type DBRepository struct {
	logger *zap.Logger
	db     *sqlx.DB
	cfg    configs.AppConfig
}

const (
	defaultMaxConns      = 200
	defaultMaxIdleConns  = 10
	defaultConnsLifetime = 0
)

func initPG(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener " +
		"(id BIGSERIAL PRIMARY KEY, user_id VARCHAR(256), url VARCHAR(256) UNIQUE NOT NULL, " +
		"shorten VARCHAR(256) UNIQUE NOT NULL, correlation_id VARCHAR(256), deleted BOOLEAN DEFAULT FALSE)")
	if err != nil {
		return fmt.Errorf("ошибка при создании базы данных: %w", err)
	}
	return nil
}

func connectPG(dsnString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsnString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}
	db.SetMaxOpenConns(defaultMaxConns)
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetConnMaxLifetime(defaultConnsLifetime)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к базе данных: %w", err)
	}
	return db, nil
}

func NewDBRepository(logger *zap.Logger, cfg configs.AppConfig) (*DBRepository, error) {
	db, err := connectPG(cfg.DatabaseDSN)
	if err != nil {
		return nil, defaultRepoErrWrapper(err)
	}
	if err := initPG(db); err != nil {
		return nil, defaultRepoErrWrapper(err)
	}
	return &DBRepository{logger: logger, cfg: cfg, db: db}, nil
}

func (r *DBRepository) validateUniqueShortURL(ctx context.Context, tx *sqlx.Tx, shortURL string) error {
	var id int

	select {
	case <-ctx.Done():
		return defaultRepoErrWrapper(ctx.Err())
	default:
	}
	shortRow := tx.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
	if err := shortRow.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return defaultRepoErrWrapper(err)
	}
	return errWithVal(ErrShortURLExists, shortURL)
}

func (r *DBRepository) Create(ctx context.Context, userID, fullURL, shortURL string) (string, error) {
	var shortRes string

	select {
	case <-ctx.Done():
		return "", defaultRepoErrWrapper(ctx.Err())
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", defaultRepoErrWrapper(err)
	}

	if err := r.validateUniqueShortURL(ctx, tx, shortURL); err != nil {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
		}()
		return "", defaultRepoErrWrapper(err)
	}
	query := "INSERT INTO shortener (url, shorten, user_id) VALUES($1, $2, $3) " +
		"ON CONFLICT(url) DO UPDATE SET shorten = shortener.shorten RETURNING shorten"
	res := tx.QueryRowxContext(ctx, query, fullURL, shortURL, userID)
	if err := res.Scan(&shortRes); err != nil {
		defer func() {
			if errDefer := tx.Rollback(); errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
		}()
		return "", defaultRepoErrWrapper(err)
	}

	if shortRes != shortURL {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
		}()
		return shortRes, ErrUniqueViolation
	}
	if err := tx.Commit(); err != nil {
		r.logger.Error(txRollbackError, zap.Error(err))
		return "", defaultRepoErrWrapper(err)
	}
	return shortRes, nil
}

func (r *DBRepository) BatchCreate(ctx context.Context,
	data []models.BatchCreateData) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, defaultRepoErrWrapper(ctx.Err())
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, defaultRepoErrWrapper(err)
	}

	for _, val := range data {
		if validateErr := r.validateUniqueShortURL(ctx, tx, val.ShortURL); validateErr != nil {
			if errDefer := tx.Rollback(); errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
			return []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: val.ShortURL, UUID: val.UUID}}, validateErr
		}
	}

	// генерируем список сокращенных урлов
	query := "INSERT INTO shortener (url, shorten, correlation_id, user_id) " +
		"VALUES(:url, :shorten, :correlation_id, :user_id) " +
		"ON CONFLICT(url)  DO UPDATE SET shorten = shortener.shorten RETURNING shorten, correlation_id"
	// создаем
	rows, err := tx.NamedQuery(query, data)
	if err != nil {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
		}()
		r.logger.Error(queryErrorMessage, zap.Error(err))
		return nil, defaultRepoErrWrapper(err)
	}
	defer func() {
		errDefer := rows.Close()
		if errDefer != nil {
			err = errDefer
		}
	}()
	// обрабатываем ответ
	for rows.Next() {
		var row = models.BatchCreateResponse{}
		if scanErr := rows.StructScan(&row); scanErr != nil {
			defer func() {
				errDefer := tx.Rollback()
				if errDefer != nil {
					r.logger.Error(txRollbackError, zap.Error(errDefer))
				}
			}()
			return nil, defaultRepoErrWrapper(scanErr)
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			r.logger.Error(defaultRepoError, zap.Error(err))
			defer func() {
				errDefer := tx.Rollback()
				if errDefer != nil {
					r.logger.Error(txRollbackError, zap.Error(errDefer))
				}
			}()
			return nil, defaultRepoErrWrapper(err)
		}
		response = append(response, row)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(defaultRepoError, zap.Error(err))
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				r.logger.Error(txRollbackError, zap.Error(errDefer))
			}
		}()
		return nil, defaultRepoErrWrapper(err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error(txRollbackError, zap.Error(err))
		return nil, defaultRepoErrWrapper(err)
	}
	return response, nil
}

//func (r *DBRepository) BatchDelete(ctx context.Context, shortURL string) (string, error) {
//	var fullURL string
//	row := r.db.QueryRowContext(ctx, "SELECT url FROM shortener WHERE  shorten = $1", shortURL)
//	if err := row.Scan(&fullURL); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return "", errWithVal(ErrURLNotFound, shortURL)
//		}
//		r.logger.Error(defaultRepoError, zap.Error(err))
//		return "", defaultRepoErrWrapper(err)
//	}
//	return fullURL, nil
//}

func (r *DBRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	var deleted bool
	row := r.db.QueryRowContext(ctx, "SELECT url, deleted FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL, &deleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errWithVal(ErrURLNotFound, shortURL)
		}
		r.logger.Error(defaultRepoError, zap.Error(err))
		return "", defaultRepoErrWrapper(err)
	}
	if deleted {
		return "", ErrDeleted
	}
	return fullURL, nil
}

func (r *DBRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return defaultRepoErrWrapper(err)
	}
	return nil
}

func (r *DBRepository) Close() error {
	if err := r.db.Close(); err != nil {
		return defaultRepoErrWrapper(err)
	}
	return nil
}

func (r *DBRepository) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error {
	query := "UPDATE shortener set deleted = True WHERE user_id = $1 and shorten = any ($2::text[])"
	_, err := r.db.ExecContext(ctx, query, userID, data)
	if err != nil {
		r.logger.Error(queryErrorMessage, zap.Error(err))
		return defaultRepoErrWrapper(err)
	}
	return nil
}

func (r *DBRepository) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	var response []models.GetURLResponse
	select {
	case <-ctx.Done():
		return nil, defaultRepoErrWrapper(ctx.Err())
	default:
	}
	query := "SELECT url, shorten FROM shortener WHERE user_id = $1 and deleted = False"
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		r.logger.Error(queryErrorMessage, zap.Error(err))
		return nil, defaultRepoErrWrapper(err)
	}
	defer func() {
		errDefer := rows.Close()
		if errDefer != nil {
			err = errDefer
		}
	}()
	// обрабатываем ответ
	for rows.Next() {
		var row = models.GetURLResponse{}
		if scanErr := rows.StructScan(&row); scanErr != nil {
			return nil, defaultRepoErrWrapper(scanErr)
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
			return nil, defaultRepoErrWrapper(err)
		}
		response = append(response, row)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		r.logger.Error("repository: ошибка при формировании ответа", zap.Error(rowsErr))
		return nil, defaultRepoErrWrapper(rowsErr)
	}
	if response == nil {
		return nil, sql.ErrNoRows
	}
	return response, nil
}
