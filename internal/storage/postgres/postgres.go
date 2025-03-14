package postgres

import (
	"database/sql"
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
