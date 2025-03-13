package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ozon/internal/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler_SuccessNewURL(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	storeMock.GetShortURLMock.Set(func(ctx context.Context, originalURL string) (string, error) {
		return "", errors.New("not found")
	})

	storeMock.SaveURLMock.Set(func(ctx context.Context, originalURL string, shortURL string) error {
		return nil
	})

	reqBody := `{"original_url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	ShortenURLHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response URLResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ShortURL)
}

func TestShortenURLHandler_ExistingURL(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	storeMock.GetShortURLMock.Set(func(ctx context.Context, originalURL string) (string, error) {
		return "abc123", nil
	})

	reqBody := `{"original_url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	ShortenURLHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response URLResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", response.ShortURL)
}

func TestShortenURLHandler_InvalidMethod(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rr := httptest.NewRecorder()

	ShortenURLHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "Method not allowed\n", rr.Body.String())
}

func TestShortenURLHandler_InvalidBody(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	reqBody := `{"invalid_field": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	ShortenURLHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Original URL is required\n", rr.Body.String())
}

func TestRedirectHandler_Success(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	storeMock.GetOriginalURLMock.Set(func(ctx context.Context, shortURL string) (string, error) {
		return "https://example.com", nil
	})

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rr := httptest.NewRecorder()

	RedirectHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", response["original_url"])
}

func TestRedirectHandler_NotFound(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	storeMock.GetOriginalURLMock.Set(func(ctx context.Context, shortURL string) (string, error) {
		return "", errors.New("not found")
	})

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rr := httptest.NewRecorder()

	RedirectHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "URL not found\n", rr.Body.String())
}

func TestRedirectHandler_InvalidMethod(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewStorageMock(mc)

	req := httptest.NewRequest(http.MethodPost, "/abc123", nil)
	rr := httptest.NewRecorder()

	RedirectHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "Method not allowed\n", rr.Body.String())
}
