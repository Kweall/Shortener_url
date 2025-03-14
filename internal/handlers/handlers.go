package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ozon/internal/custom_errors"
	"strings"
)

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

//go:generate minimock -g -i *  -o ./mocks/ -s ".mock.go"
type Service interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)
	OriginalURL(ctx context.Context, shortURL string) (string, error)
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request, service Service) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OriginalURL == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	shortURL, err := Service.ShortenURL(service, ctx, req.OriginalURL)
	if err != nil {
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(URLResponse{ShortURL: shortURL})
}

func OriginalURLHandler(w http.ResponseWriter, r *http.Request, service Service) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := service.OriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "service.Redirect", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request, service Service) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	shortURL := strings.TrimPrefix(r.URL.Path, "/redirect/")
	originalURL, err := service.OriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
