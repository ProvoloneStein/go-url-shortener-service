package repositories

import (
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(dsnString string) (*PostgresRepository, error) {
	db, err := sqlx.Open("pgx", dsnString)
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
