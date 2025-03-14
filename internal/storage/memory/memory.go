package memory

import (
	"sync"
)

type MemoryStorage struct {
	urls            map[string]string
	originalToShort map[string]string
	mu              sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls:            make(map[string]string),
		originalToShort: make(map[string]string),
	}
}
