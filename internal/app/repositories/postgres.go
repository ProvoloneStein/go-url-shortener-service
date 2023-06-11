package repositories

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsnString string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dsnString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Create(fullURL string) (string, error) {
	return "", errors.New("service not working")
}

func (r *PostgresRepository) GetByShort(shortURL string) (string, error) {
	return "", errors.New("service not working")
}

func (r *PostgresRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
