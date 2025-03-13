package storage

import (
	"context"
	"fmt"
	"sync"
)

type MemoryStorage struct {
	urls            map[string]string
	originalToShort map[string]string
	mu              sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls:            make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

func (m *MemoryStorage) SaveURL(ctx context.Context, originalURL, shortURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.urls[shortURL] = originalURL
	m.originalToShort[originalURL] = shortURL
	return nil
}

func (m *MemoryStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	original, exists := m.urls[shortURL]
	if !exists {
		return "", fmt.Errorf("url not found")
	}
	return original, nil
}

func (m *MemoryStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	short, exists := m.originalToShort[originalURL]
	if !exists {
		return "", fmt.Errorf("url not found")
	}
	return short, nil
}
