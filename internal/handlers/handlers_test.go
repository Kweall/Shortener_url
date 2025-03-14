package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ozon/internal/custom_errors"
	"ozon/internal/handlers/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler_SuccessNewURL(t *testing.T) {
	// arrange
	mc := minimock.NewController(t)

	storeMock := mocks.NewServiceMock(mc)

	storeMock.ShortenURLMock.Set(func(ctx context.Context, originalURL string) (string, error) {
		return "plN_OAp1px", nil
	})

	reqBody := `{"original_url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	// act
	ShortenURLHandler(rr, req, storeMock)
	// assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response URLResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.ShortURL, 10)
}

func TestShortenURLHandler_ExistingURL(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewServiceMock(mc)

	storeMock.ShortenURLMock.Set(func(ctx context.Context, originalURL string) (string, error) {
		return "plN_OAp1px", nil
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
	assert.Equal(t, "plN_OAp1px", response.ShortURL)
}

func TestShortenURLHandler_InvalidMethod(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewServiceMock(mc)

	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rr := httptest.NewRecorder()

	ShortenURLHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "Method not allowed\n", rr.Body.String())
}

func TestShortenURLHandler_InvalidBody(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewServiceMock(mc)

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

	storeMock := mocks.NewServiceMock(mc)

	storeMock.RedirectMock.Expect(context.Background(), "plN_OAp1px").Return("https://example.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/plN_OAp1px", nil)
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

	storeMock := mocks.NewServiceMock(mc)

	storeMock.RedirectMock.Expect(context.Background(), "plN_OAp1px").Return("", custom_errors.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/plN_OAp1px", nil)
	rr := httptest.NewRecorder()

	RedirectHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "URL not found\n", rr.Body.String())
}

func TestRedirectHandler_InvalidMethod(t *testing.T) {
	mc := minimock.NewController(t)

	storeMock := mocks.NewServiceMock(mc)

	req := httptest.NewRequest(http.MethodPost, "/plN_OAp1px", nil)
	rr := httptest.NewRecorder()

	RedirectHandler(rr, req, storeMock)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "Method not allowed\n", rr.Body.String())
}
