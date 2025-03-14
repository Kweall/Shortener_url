package postgres

import "context"

func (p *PostgresStorage) SaveURL(ctx context.Context, originalURL, shortURL string) error {
	_, err := p.db.Exec(
		"INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING",
		originalURL, shortURL,
	)

	return err
}
