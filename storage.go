package main

import (
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Save(short, original string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[short] = original
}

func (s *Storage) Get(short string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[short]
	return val, ok
}

func (s *Storage) Exists(short string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[short]
	return ok
}
