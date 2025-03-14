package service_test

import (
	"context"
	"ozon/internal/service"
	"ozon/internal/service/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestRedirect_ReturnsOriginal(t *testing.T) {
	mc := minimock.NewController(t)

	storageMock := mocks.NewStorageMock(mc)

	ctx := context.Background()
	shortURL := "plN_OAp1px"
	originalURL := "https://example.com"

	// Ожидаем, что при вызове GetOriginalURL вернётся оригинальный URL.
	storageMock.GetOriginalURLMock.Expect(ctx, shortURL).Return(originalURL, nil)

	svc := service.NewService(storageMock)
	result, err := svc.OriginalURL(ctx, shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)
}
