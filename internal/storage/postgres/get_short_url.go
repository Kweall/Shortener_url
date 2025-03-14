package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func (p *PostgresStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	var short string
	err := p.db.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&short)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("url not found, err: %w", err)
	}

	return short, err
}
