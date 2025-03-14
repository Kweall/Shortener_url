package service

import (
	"context"
	"fmt"
)

func (s *ServiceImpl) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	if shortURL, err := s.storage.GetShortURL(ctx, originalURL); err == nil {
		return shortURL, nil
	}

	var shortURL string
	for {
		var err error
		shortURL, err = s.generateShortURL()
		if err != nil {
			return "", fmt.Errorf("failed to generate, err: %w", err)
		}
		if err = s.storage.SaveURL(ctx, originalURL, shortURL); err == nil {
			break
		}
	}

	return shortURL, nil
}
