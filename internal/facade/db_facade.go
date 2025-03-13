package facade

import (
	"context"
	"ozon/internal/storage"
)

type DBFacade interface {
	SaveURL(ctx context.Context, originalURL, shortURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
}

type dbFacade struct {
	storage storage.Storage
}

func NewDBFacade(storage storage.Storage) DBFacade {
	return &dbFacade{storage: storage}
}

func (f *dbFacade) SaveURL(ctx context.Context, originalURL, shortURL string) error {
	return f.storage.SaveURL(ctx, originalURL, shortURL)
}

func (f *dbFacade) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return f.storage.GetOriginalURL(ctx, shortURL)
}

func (f *dbFacade) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	return f.storage.GetShortURL(ctx, originalURL)
}
