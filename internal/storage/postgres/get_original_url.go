package postgres

import (
	"context"
	"database/sql"
	"errors"
	"ozon/internal/custom_errors"
)

func (p *PostgresStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var original string
	err := p.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&original)
	if errors.Is(err, sql.ErrNoRows) {
		return "", custom_errors.ErrNoRows
	}

	return original, err
}
