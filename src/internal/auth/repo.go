package auth

import "sync"

type TokenStore interface {
	Save(token, userID string)
	IsValid(token string) bool
	GetUserID(token string) (string, bool)
	Delete(token string)
}

type InMemoryTokenStore struct {
	tokens map[string]string
	mu     sync.RWMutex
}

func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{tokens: make(map[string]string)}
}

func (s *InMemoryTokenStore) Save(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = userID
}

func (s *InMemoryTokenStore) IsValid(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.tokens[token]
	return ok
}

func (s *InMemoryTokenStore) GetUserID(token string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uid, ok := s.tokens[token]
	return uid, ok
}

func (s *InMemoryTokenStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}
