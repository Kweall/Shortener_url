package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"ozon/internal/config"
	"ozon/internal/facade"
	"ozon/internal/storage"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupPostgresTestDB() (*sql.DB, string, error) {
	if err := godotenv.Load("C:/Projects/ozon/.env"); err != nil {
		log.Printf("Failed to load from .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, "", err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id SERIAL PRIMARY KEY,
            original_url TEXT UNIQUE,
            short_url VARCHAR(10) UNIQUE
        )
    `)
	if err != nil {
		return nil, "", err
	}

	return db, connStr, nil
}

func TestPostgresStorage_SaveAndGet(t *testing.T) {
	db, connStr, err := setupPostgresTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer db.Close()

	store, err := storage.NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("Failed to open connection to PostgreSQL: %v", err)
	}
	facade := facade.NewDBFacade(store)

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123XYZ_"

	err = facade.SaveURL(ctx, originalURL, shortURL)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	retrievedOriginal, err := facade.GetOriginalURL(ctx, shortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL: %v", err)
	}
	if retrievedOriginal != originalURL {
		t.Errorf("Expected original URL %s, got %s", originalURL, retrievedOriginal)
	}

	retrievedShort, err := facade.GetShortURL(ctx, originalURL)
	if err != nil {
		t.Fatalf("Failed to get short URL: %v", err)
	}
	if retrievedShort != shortURL {
		t.Errorf("Expected short URL %s, got %s", shortURL, retrievedShort)
	}

	_, err = facade.GetOriginalURL(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent short URL, got nil")
	}
}
