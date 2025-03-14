package memory

import "context"

func (m *MemoryStorage) SaveURL(ctx context.Context, originalURL, shortURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.urls[shortURL] = originalURL
	m.originalToShort[originalURL] = shortURL

	return nil
}
