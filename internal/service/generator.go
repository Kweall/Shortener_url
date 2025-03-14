package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	shortURLLength = 10
)

func (s *ServiceImpl) generateShortURL() (string, error) {
	result := make([]byte, shortURLLength)
	for i := range result {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}
		result[i] = charset[idx.Int64()]
	}

	return string(result), nil
}
