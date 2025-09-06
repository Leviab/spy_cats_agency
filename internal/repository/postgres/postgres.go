package postgres

import (
	"fmt"
	"spy_cats_agency/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver
)

// DB is a wrapper for the sqlx.DB that provides database connection.
type DB struct {
	*sqlx.DB
}

// New creates a new database connection.
func New(cfg config.Config) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{db}, nil
}
