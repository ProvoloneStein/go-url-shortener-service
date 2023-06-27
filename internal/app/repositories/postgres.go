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
	cfg    configs.AppConfig
	db     *sqlx.DB
}

func initPG(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener " +
		"(id BIGSERIAL PRIMARY KEY, user_id VARCHAR(256), url VARCHAR(256) UNIQUE NOT NULL, shorten VARCHAR(256) UNIQUE NOT NULL, correlation_id VARCHAR(256), deleted BOOLEAN DEFAULT FALSE)")
	if err != nil {
		return fmt.Errorf("repository: ошибка при создании базы данных: %w", err)
	}
	return nil
}

func connectPG(dsnString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsnString)
	if err != nil {
		return nil, fmt.Errorf("repository: ошибка при подключении к базе данных: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("repository: ошибка при подключении к базе данных: %s", err)
	}
	return db, nil
}

func NewDBRepository(logger *zap.Logger, cfg configs.AppConfig) (*DBRepository, error) {
	db, err := connectPG(cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	if err := initPG(db); err != nil {
		return nil, fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
	}
	return &DBRepository{logger: logger, cfg: cfg, db: db}, nil
}

func (r *DBRepository) validateUniqueShortURL(ctx context.Context, tx *sqlx.Tx, shortURL string) error {
	var id int

	select {
	case <-ctx.Done():
		return fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	shortRow := tx.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
	if err := shortRow.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
	}
	return fmt.Errorf("%w: %s", ErrShortURLExists, shortURL)
}

func (r *DBRepository) Create(ctx context.Context, userID, fullURL, shortURL string) (string, error) {
	var shortRes string

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("repository: %w", ctx.Err())
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("repository: %w", err)
	}
	defer func() {
		errDefer := tx.Commit()
		if errDefer != nil {
			err = errDefer
		}
	}()

	if err := r.validateUniqueShortURL(ctx, tx, shortURL); err != nil {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				err = errDefer
			}
		}()
		return "", fmt.Errorf("repository: %w", err)
	}
	query := "INSERT INTO shortener (url, shorten, user_id) VALUES($1, $2, $3) ON CONFLICT(url) DO UPDATE SET shorten = shortener.shorten RETURNING shorten"
	res := tx.QueryRowContext(ctx, query, fullURL, shortURL, userID)
	if err := res.Scan(&shortRes); err != nil {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				err = errDefer
			}
		}()
		return "", fmt.Errorf("repository: %w", err)
	}

	if shortRes != shortURL {
		defer func() {
			errDefer := tx.Rollback()
			if errDefer != nil {
				r.logger.Error("repository:  ошибка при откате трансакции", zap.Error(errDefer))
			}
		}()
		return shortRes, ErrorUniqueViolation
	}

	return shortRes, tx.Commit()
}

func (r *DBRepository) BatchCreate(ctx context.Context, data []models.BatchCreateData) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("repository: %w", ctx.Err())
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}

	defer func() {
		errDefer := tx.Commit()
		if errDefer != nil {
			err = errDefer
		}
	}()
	// генерируем список сокращенных урлов
	query := "INSERT INTO shortener (url, shorten, correlation_id, user_id) VALUES(:url, :shorten, :correlation_id, :user_id) ON CONFLICT(url)  DO UPDATE SET shorten = shortener.shorten RETURNING shorten, correlation_id"
	// создаем
	rows, err := tx.NamedQuery(query, data)
	if err != nil {
		errDefer := tx.Rollback()
		if errDefer != nil {
			err = errDefer
		}
		r.logger.Error("repository:  ошибка при запросе к бд", zap.Error(err))
		return nil, fmt.Errorf("repository:  ошибка при запросе к бд: %w", err)
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
		if err := rows.StructScan(&row); err != nil {
			return nil, tx.Rollback()
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = errRollback
			}
			return nil, fmt.Errorf("repository: %w", err)
		}
		response = append(response, row)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = errRollback
		}
		return nil, fmt.Errorf("repository: %w", err)
	}

	return response, tx.Commit()
}

func (r *DBRepository) BatchDelete(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	row := r.db.QueryRowContext(ctx, "SELECT url FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %s", ErrURLNotFound, shortURL)
		}
		r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
		return "", fmt.Errorf("repository: ошибка при формировании ответа: %s", err)
	}
	return fullURL, nil
}

func (r *DBRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	var deleted bool
	row := r.db.QueryRowContext(ctx, "SELECT url, deleted FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL, &deleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %s", ErrURLNotFound, shortURL)
		}
		r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
		return "", fmt.Errorf("repository: ошибка при формировании ответа: %s", err)
	}
	if deleted {
		return "", ErrDeleted
	}
	return fullURL, nil
}

func (r *DBRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
	}
	return nil
}

func (r *DBRepository) Close() error {
	return r.db.Close()
}

func (r *DBRepository) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	_, err := r.db.ExecContext(ctx, "UPDATE shortener set deleted = True WHERE user_id LIKE $1 and shorten = any ($2::text[])", userID, data)
	if err != nil {
		r.logger.Error("repository: ошибка при запросе к бд", zap.Error(err))
		return fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
	}
	return nil
}

func (r *DBRepository) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	var response []models.GetURLResponse
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	query := "SELECT url, shorten FROM shortener WHERE user_id LIKE $1 and deleted = False"
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		r.logger.Error("repository: ошибка при запросе к бд", zap.Error(err))
		return nil, fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
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
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("repository: %w", err)
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
			return nil, fmt.Errorf("repository: %w", err)
		}
		response = append(response, row)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("repository: ошибка при формировании ответа", zap.Error(err))
		return nil, fmt.Errorf("repository: %w", err)
	}
	if response == nil {
		return nil, sql.ErrNoRows
	}
	return response, nil
}

func (r *DBRepository) ValidateUniqueUser(ctx context.Context, userID string) error {
	var id int

	select {
	case <-ctx.Done():
		return fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	shortRow := r.db.QueryRowContext(ctx, "SELECT id FROM shortener WHERE user_id = $1", userID)
	if err := shortRow.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("repository: ошибка при запросе к бд: %w", err)
	}
	return fmt.Errorf("%w: %s", ErrUserExists, userID)
}
