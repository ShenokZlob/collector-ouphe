package usecase

import "sync"

type inMemoryStorage struct {
	users map[int64]string
	mu    sync.RWMutex
}

func newInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		users: make(map[int64]string),
		mu:    sync.RWMutex{},
	}
}

func (s *inMemoryStorage) AddUser(id int64, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[id] = token
}

func (s *inMemoryStorage) GetUser(id int64) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, ok := s.users[id]
	return token, ok
}
