package service

import (
	"context"
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
