package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"ozon/internal/storage"
	"ozon/internal/utils"
)

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {

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

	if shortURL, err := store.GetShortURL(ctx, req.OriginalURL); err == nil {
		json.NewEncoder(w).Encode(URLResponse{ShortURL: shortURL})
		return
	}

	var shortURL string
	for {
		var err error
		shortURL, err = utils.GenerateShortURL()
		if err != nil {
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}
		if err = store.SaveURL(ctx, req.OriginalURL, shortURL); err == nil {
			break
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(URLResponse{ShortURL: shortURL})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := store.GetOriginalURL(ctx, shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}
