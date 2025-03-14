package service_test

import (
	"context"
	"testing"

	"ozon/internal/custom_errors"
	"ozon/internal/service"
	"ozon/internal/service/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_ReturnsExisting(t *testing.T) {
	mc := minimock.NewController(t)

	// Создаём мок для интерфейса Storage.
	storageMock := mocks.NewStorageMock(mc)

	ctx := context.Background()
	originalURL := "https://example.com"
	existingShort := "plN_OAp1px"

	// Если в хранилище уже существует короткая ссылка, возвращаем её.
	storageMock.GetShortURLMock.Expect(ctx, originalURL).Return(existingShort, nil)

	// Создаём сервис.
	svc := service.NewService(storageMock)

	// Вызываем метод ShortenURL.
	short, err := svc.ShortenURL(ctx, originalURL)
	assert.NoError(t, err)
	assert.Equal(t, existingShort, short)
}

func TestShortenURL_Generate(t *testing.T) {
	mc := minimock.NewController(t)

	storageMock := mocks.NewStorageMock(mc)

	ctx := context.Background()
	originalURL := "https://example.com"

	storageMock.GetShortURLMock.Expect(ctx, originalURL).Return("", custom_errors.ErrNoRows)

	storageMock.SaveURLMock.Set(func(ctx context.Context, originalURL, shortURL string) error {
		return nil
	})

	baseService := service.NewService(storageMock)

	short, err := baseService.ShortenURL(ctx, originalURL)
	assert.NoError(t, err)
	assert.NotEmpty(t, short)
}
