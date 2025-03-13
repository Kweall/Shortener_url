package storage_test

import (
	"context"
	"ozon/internal/storage"
	"testing"
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {
	store := storage.NewMemoryStorage()

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123XYZ_"

	err := store.SaveURL(ctx, originalURL, shortURL)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	retrievedOriginal, err := store.GetOriginalURL(ctx, shortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL: %v", err)
	}
	if retrievedOriginal != originalURL {
		t.Errorf("Expected original URL %s, got %s", originalURL, retrievedOriginal)
	}

	retrievedShort, err := store.GetShortURL(ctx, originalURL)
	if err != nil {
		t.Fatalf("Failed to get short URL: %v", err)
	}
	if retrievedShort != shortURL {
		t.Errorf("Expected short URL %s, got %s", shortURL, retrievedShort)
	}

	_, err = store.GetOriginalURL(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent short URL, got nil")
	}
}
