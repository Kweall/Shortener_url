package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open connection to PostgreSQL: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping PostgreSQL: %v", err)
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id SERIAL PRIMARY KEY,
            original_url TEXT UNIQUE,
            short_url VARCHAR(10) UNIQUE
        )
    `)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) SaveURL(ctx context.Context, originalURL, shortURL string) error {
	_, err := p.db.Exec(
		"INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING",
		originalURL, shortURL,
	)
	return err
}

func (p *PostgresStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var original string
	err := p.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&original)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("url not found, err: %w", err)
	}
	return original, err
}

func (p *PostgresStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	var short string
	err := p.db.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&short)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("url not found, err: %w", err)
	}
	return short, err
}
