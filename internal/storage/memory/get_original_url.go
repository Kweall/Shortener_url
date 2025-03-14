package memory

import (
	"context"
	"fmt"
)

func (m *MemoryStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	original, exists := m.urls[shortURL]
	if !exists {
		return "", fmt.Errorf("url not found")
	}
	return original, nil
}
