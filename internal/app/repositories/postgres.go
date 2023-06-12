package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/url"
)

type PostgresRepository struct {
	cfg configs.AppConfig
	db  *sql.DB
}

func initPG(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener " +
		"(id BIGSERIAL PRIMARY KEY, url TEXT UNIQUE NOT NULL, shorten TEXT UNIQUE NOT NULL)")
	if err != nil {
		return err
	}
	return nil
}

func ConntectPG(dsnString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsnString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewPostgresRepository(cfg configs.AppConfig, db *sql.DB) (*PostgresRepository, error) {
	if err := initPG(db); err != nil {
		return nil, err
	}
	return &PostgresRepository{cfg, db}, nil
}

const UniqueViolation = "23505"

func (r *PostgresRepository) Create(ctx context.Context, fullURL string) (string, error) {
	var id int
	var pgErr *pgconn.PgError

	for {
		shortURL := randomString()

		shortRow := r.db.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
		if err := shortRow.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				if _, err := r.db.ExecContext(ctx, "INSERT INTO shortener (url, shorten) VALUES($1, $2)", fullURL, shortURL); err != nil {
					if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
						row := r.db.QueryRowContext(ctx, "SELECT shorten FROM shortener WHERE  url = $1", fullURL)
						if err := row.Scan(&shortURL); err != nil {
							return "", err
						}
						return shortURL, ErrorUniqueViolation
					}
					return "", err
				}
				return shortURL, nil
			}
			return "", err
		}
	}
}

func (r *PostgresRepository) BatchCreate(ctx context.Context, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	for _, val := range data {
		for {
			shortURL := randomString()
			var id int
			row := tx.QueryRowContext(ctx, "SELECT id FROM shortener WHERE  shorten = $1", shortURL)
			err = row.Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					if _, err := tx.ExecContext(ctx, "INSERT INTO shortener (url, shorten) VALUES($1, $2)", val.URL, shortURL); err == nil {
						resShortURL, err := url.JoinPath(r.cfg.BaseURL, shortURL)
						if err == nil {
							response = append(response, models.BatchCreateResponse{URL: resShortURL, UUID: val.UUID})
						}
					}
				}
				break
			}
		}
	}
	return response, tx.Commit()
}

func (r *PostgresRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	row := r.db.QueryRowContext(ctx, "SELECT url FROM shortener WHERE  shorten = $1", shortURL)
	if err := row.Scan(&fullURL); err != nil {
		return "", err
	}
	return fullURL, nil
}

func (r *PostgresRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
