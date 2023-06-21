package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/url"
)

type PostgresRepository struct {
	logger *zap.Logger
	cfg    configs.AppConfig
	db     *sqlx.DB
}

func initPG(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener " +
		"(id BIGSERIAL PRIMARY KEY, url VARCHAR(256) UNIQUE NOT NULL, shorten VARCHAR(256) UNIQUE NOT NULL, correlation_id VARCHAR(256))")
	if err != nil {
		return err
	}
	return nil
}

func ConntectPG(dsnString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsnString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewPostgresRepository(logger *zap.Logger, cfg configs.AppConfig) (*PostgresRepository, error) {
	db, err := ConntectPG(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	if err := initPG(db); err != nil {
		return nil, fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return &PostgresRepository{logger: logger, cfg: cfg, db: db}, nil
}

func (r *PostgresRepository) GenerateShortURL(ctx context.Context) (string, error) {
	var id int

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	for {
		shortURL := randomString()
		shortRow := r.db.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
		if err := shortRow.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return shortURL, nil
			}
			return "", err
		}
	}
}

func (r *PostgresRepository) Create(ctx context.Context, fullURL, shortURL string) (string, error) {
	var shortRes string

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	res := r.db.QueryRowContext(ctx, "INSERT INTO shortener (url, shorten) VALUES($1, $2) ON CONFLICT(url) DO UPDATE SET shorten = shortener.shorten RETURNING shorten", fullURL, shortURL)
	if err := res.Scan(&shortRes); err != nil {
		return "", err
	}
	if shortRes != shortURL {
		return shortRes, ErrorUniqueViolation
	}
	return shortRes, nil
}

func (r *PostgresRepository) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
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
	// генерируем список сокращенных урлов
	for _, val := range data {
	generator:
		for {
			shortURL, err := r.GenerateShortURL(ctx)
			if err != nil {
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
		r.logger.Error("ошибка при запросе url", zap.Error(err))
		return nil, err
	}
	// обрабатываем ответ
	for rows.Next() {
		var row = models.BatchCreateResponse{}
		if err := rows.StructScan(&row); err != nil {
			return nil, tx.Rollback()
		}
		row.ShortURL, err = url.JoinPath(r.cfg.BaseURL, row.ShortURL)
		if err != nil {
			r.logger.Error("ошибка при формировании url", zap.Error(err))
			return nil, tx.Rollback()
		}
		response = append(response, row)
	}

	return response, tx.Commit()
}

func (r *PostgresRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	row := r.db.QueryRowContext(ctx, "SELECT url FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL); err != nil {
		return "", fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return fullURL, nil
}

func (r *PostgresRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
