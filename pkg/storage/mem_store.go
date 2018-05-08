package storage

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

// MemStore is an in-memory Store implementation.
type MemStore struct {
	lock sync.RWMutex
	mem  map[string][]byte
}

// NewMemStore returns a new MemoryStore.
func NewMemStore() *MemStore {
	return &MemStore{
		mem: make(map[string][]byte),
	}
}

// Get implements the Store interface.
func (s *MemStore) Get(key []byte) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if val, ok := s.mem[string(key)]; ok {
		return copyBytes(val), nil
	}
	return nil, errNotFound
}

// Put implements the Store interface.
func (s *MemStore) Put(key, val []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.mem[string(key)] = copyBytes(val)
	return nil
}

// Has implements the Store interface.
func (s *MemStore) Has(key []byte) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.mem[string(key)]
	return ok
}

// Len makes us use the Memory as a Batcher.
func (s *MemStore) Len() int { return 0 }

// copyBytes will preven nasty stuff from happening.
func copyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	dest := make([]byte, len(b))
	copy(dest, b)
	return dest
}
