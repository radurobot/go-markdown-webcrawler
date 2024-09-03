package storage

import "sync"

type InMemoryStore struct {
	mu     sync.RWMutex
	hashes map[string]struct{}
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		hashes: make(map[string]struct{}),
	}
}

func (s *InMemoryStore) Add(hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashes[hash] = struct{}{}
}

func (s *InMemoryStore) Exists(hash string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.hashes[hash]
	return exists
}

func (s *InMemoryStore) Close() error {
	return nil
}
