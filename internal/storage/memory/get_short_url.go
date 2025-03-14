package memory

import (
	"context"
	"fmt"
)

func (m *MemoryStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	short, exists := m.originalToShort[originalURL]
	if !exists {
		return "", fmt.Errorf("url not found")
	}
	return short, nil
}
