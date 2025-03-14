package memory

import (
	"context"
	"ozon/internal/custom_errors"
)

func (m *MemoryStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	short, exists := m.originalToShort[originalURL]
	if !exists {
		return "", custom_errors.ErrNoRows
	}

	return short, nil
}
