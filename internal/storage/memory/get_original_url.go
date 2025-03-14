package memory

import (
	"context"
	"ozon/internal/custom_errors"
)

func (m *MemoryStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	original, exists := m.urls[shortURL]
	if !exists {
		return "", custom_errors.ErrNoRows
	}

	return original, nil
}
