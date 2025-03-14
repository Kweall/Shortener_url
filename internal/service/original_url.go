package service

import "context"

func (s *ServiceImpl) OriginalURL(ctx context.Context, shortURL string) (string, error) {
	return s.storage.GetOriginalURL(ctx, shortURL)
}
