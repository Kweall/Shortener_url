package utils

import (
	"strings"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	shortURL, err := GenerateShortURL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(shortURL) != 10 {
		t.Errorf("Expected short URL length to be 10, got %d", len(shortURL))
	}

	for _, char := range shortURL {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_", char) {
			t.Errorf("Invalid character in short URL: %c", char)
		}
	}
}
