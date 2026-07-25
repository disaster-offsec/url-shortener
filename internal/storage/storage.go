package storage

import (
	"sync"
)

type Storage struct {
	mu        sync.RWMutex
	codeToURL map[string]string // short_url -> original_url (амортизированное O(1))
	urlToCode map[string]string // original_url -> short_url (амортизированное O(1))
}

func NewStorage() *Storage {
	return &Storage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *Storage) Save(short, original string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codeToURL[short] = original
	s.urlToCode[original] = short
}

// Возвращаем original_url по short_url
func (s *Storage) GetByCode(short string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.codeToURL[short]
	return val, ok
}

// Возвращаем short_url по original_url
func (s *Storage) GetByOriginal(original string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.urlToCode[original]
	return val, ok
}

// Проверяем на наличие по short_url
func (s *Storage) ExistCode(short string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.codeToURL[short]
	return ok
}

// Проверяем на наличие по original_url
func (s *Storage) ExistURL(original string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.urlToCode[original]
	return ok
}
