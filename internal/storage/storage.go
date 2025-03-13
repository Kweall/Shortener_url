package storage

import "context"

type Storage interface {
	SaveURL(ctx context.Context, originalURL, shortURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
}
