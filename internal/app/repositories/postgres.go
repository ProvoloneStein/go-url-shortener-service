package repositories

import (
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
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener (id BIGSERIAL PRIMARY KEY, url TEXT UNIQUE NOT NULL, shorten TEXT UNIQUE NOT NULL)")
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

func (r *PostgresRepository) Create(fullURL string) (string, error) {
	for {
		shortURL := randomString()
		var exists bool
		var pgErr *pgconn.PgError

		row := r.db.QueryRow("SELECT id FROM shortener WHERE  shorten = $1", shortURL)
		if err := row.Scan(&exists); err != nil {
			if err == sql.ErrNoRows {
				if _, err := r.db.Exec("INSERT INTO shortener (url, shorten) VALUES($1, $2)", fullURL, shortURL); err != nil {
					if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
						return "", UniqueViolationError
					}
					return "", err
				}
				return shortURL, nil
			}
			return "", err
		}
	}
}

func (r *PostgresRepository) BatchCreate(data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse
	tx, err := r.db.Begin()

	if err != nil {
		return nil, err
	}

	for _, val := range data {
		for {
			shortURL := randomString()
			var exists bool
			row := tx.QueryRow("SELECT id FROM shortener WHERE  shorten = $1", shortURL)
			err = row.Scan(&exists)
			if err != nil {
				if err == sql.ErrNoRows {
					if _, err := tx.Exec("INSERT INTO shortener (url, shorten) VALUES($1, $2)", val.URL, shortURL); err == nil {
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

func (r *PostgresRepository) GetByShort(shortURL string) (string, error) {
	var fullURL string
	row := r.db.QueryRow("SELECT url FROM shortener WHERE  shorten = $1", shortURL)
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
