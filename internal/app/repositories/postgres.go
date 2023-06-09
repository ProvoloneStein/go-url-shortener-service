package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/url"
)

type DBRepository struct {
	logger *zap.Logger
	cfg    configs.AppConfig
	db     *sqlx.DB
}

func initPG(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener " +
		"(id BIGSERIAL PRIMARY KEY, url VARCHAR(256) UNIQUE NOT NULL, shorten VARCHAR(256) UNIQUE NOT NULL, correlation_id VARCHAR(256))")
	if err != nil {
		return fmt.Errorf("ошибка при создании базы данных: %s", err)
	}
	return nil
}

func connectPG(dsnString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsnString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %s", err)
	}
	return db, nil
}

func NewDBRepository(logger *zap.Logger, cfg configs.AppConfig) (*DBRepository, error) {
	db, err := connectPG(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	if err := initPG(db); err != nil {
		return nil, fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return &DBRepository{logger: logger, cfg: cfg, db: db}, nil
}

func (r *DBRepository) validateUniqueShortURL(ctx context.Context, tx *sqlx.Tx, shortURL string) error {
	var id int

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	shortRow := tx.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
	if err := shortRow.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return fmt.Errorf("%w: %s", ErrShortURLExists, shortURL)
}

func (r *DBRepository) Create(ctx context.Context, fullURL, shortURL string) (string, error) {
	var shortRes string

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Commit()

	if err := r.validateUniqueShortURL(ctx, tx, shortURL); err != nil {
		tx.Rollback()
		return "", err
	}
	res := tx.QueryRowContext(ctx, "INSERT INTO shortener (url, shorten) VALUES($1, $2) ON CONFLICT(url) DO UPDATE SET shorten = shortener.shorten RETURNING shorten", fullURL, shortURL)
	if err := res.Scan(&shortRes); err != nil {
		tx.Rollback()
		return "", err
	}

	if shortRes != shortURL {
		tx.Rollback()
		return shortRes, ErrorUniqueViolation
	}

	return shortRes, tx.Commit()
}

func (r *DBRepository) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var queryData []models.BatchCreateData
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Commit()
	// генерируем список сокращенных урлов
	for _, val := range data {
	generator:
		for {
			shortURL := RandomString()
			err := r.validateUniqueShortURL(ctx, tx, shortURL)
			if err != nil {
				if errors.Is(err, ErrShortURLExists) {
					continue
				}
				tx.Rollback()
				return nil, err
			}
			// проверяем, что не задублировали ссылку
			for _, row := range queryData {
				if row.ShortURL == shortURL {
					continue generator
				}
			}
			queryData = append(queryData, models.BatchCreateData{URL: val.URL, UUID: val.UUID, ShortURL: shortURL})
			break
		}
	}
	query := "INSERT INTO shortener (url, shorten, correlation_id) VALUES(:url, :shorten, :correlation_id) ON CONFLICT(url)  DO UPDATE SET shorten = shortener.shorten RETURNING shorten, correlation_id"
	// создаем
	rows, err := tx.NamedQuery(query, queryData)
	if err != nil {
		tx.Rollback()
		r.logger.Error("ошибка при запросе к бд", zap.Error(err))
		return nil, fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	defer rows.Close()
	// обрабатываем ответ
	for rows.Next() {
		var row = models.BatchCreateResponse{}
		if err := rows.StructScan(&row); err != nil {
			return nil, tx.Rollback()
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			tx.Rollback()
			r.logger.Error("ошибка при формировании ответа", zap.Error(err))
			return nil, err
		}
		response = append(response, row)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		r.logger.Error("ошибка при формировании ответа", zap.Error(err))
		return nil, err
	}

	return response, tx.Commit()
}

func (r *DBRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	row := r.db.QueryRowContext(ctx, "SELECT url FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %s", ErrURLNotFound, shortURL)
		}
		r.logger.Error("ошибка при формировании ответа", zap.Error(err))
		return "", fmt.Errorf("ошибка при формировании ответа: %s", err)
	}
	return fullURL, nil
}

func (r *DBRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return nil
}

func (r *DBRepository) Close() error {
	return r.db.Close()
}
