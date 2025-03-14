package memory_test

import (
	"context"
	"ozon/internal/storage/memory"
	"testing"
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {
	storage := memory.NewMemoryStorage()

	ctx := context.Background()

	originalURL := "https://example.com"
	shortURL := "abc123XYZ_"

	err := storage.SaveURL(ctx, originalURL, shortURL)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	retrievedOriginal, err := storage.GetOriginalURL(ctx, shortURL)
	if err != nil {
		t.Fatalf("Failed to get original URL: %v", err)
	}
	if retrievedOriginal != originalURL {
		t.Errorf("Expected original URL %s, got %s", originalURL, retrievedOriginal)
	}

	retrievedShort, err := storage.GetShortURL(ctx, originalURL)
	if err != nil {
		t.Fatalf("Failed to get short URL: %v", err)
	}
	if retrievedShort != shortURL {
		t.Errorf("Expected short URL %s, got %s", shortURL, retrievedShort)
	}

	_, err = storage.GetOriginalURL(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent short URL, got nil")
	}
}
