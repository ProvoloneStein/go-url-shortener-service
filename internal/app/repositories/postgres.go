package repositories

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func initPG(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS shortener (id BIGSERIAL PRIMARY KEY, url TEXT NOT NULL, shorten TEXT UNIQUE NOT NULL)")
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

func NewPostgresRepository(db *sql.DB) (*PostgresRepository, error) {
	if err := initPG(db); err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Create(fullURL string) (string, error) {
	for {
		shortURL := randomString()
		var exists bool
		row := r.db.QueryRow("SELECT id FROM shortener WHERE  shorten = $1", shortURL)
		if err := row.Scan(&exists); err != nil {
			if err == sql.ErrNoRows {
				if _, err := r.db.Exec("INSERT INTO shortener (url, shorten) VALUES($1, $2)", fullURL, shortURL); err != nil {
					return "", err
				}
				return shortURL, nil
			}
			return "", err
		}
	}
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
