package service

import (
	"context"
	"fmt"
)

//go:generate minimock -g -i *  -o ./mocks/ -s ".mock.go"

type Storage interface {
	SaveURL(ctx context.Context, originalURL, shortURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
}

type ServiceImpl struct {
	storage Storage
}

func NewService(storage Storage) *ServiceImpl {
	return &ServiceImpl{storage: storage}
}

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

func (s *ServiceImpl) Redirect(ctx context.Context, shortURL string) (string, error) {
	return s.storage.GetOriginalURL(ctx, shortURL)
}
