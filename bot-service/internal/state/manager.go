package state

import (
	"sync"
)

type Manager interface {
	SetState(userID int64, state string)
	GetState(userID int64) (string, bool)
	ClearState(userID int64)
}

type memoryManager struct {
	mu     sync.RWMutex
	states map[int64]string
}

func NewMemoryManager() Manager {
	return &memoryManager{
		states: make(map[int64]string),
	}
}

func (m *memoryManager) SetState(userID int64, state string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[userID] = state
}

func (m *memoryManager) GetState(userID int64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.states[userID]
	return s, ok
}

func (m *memoryManager) ClearState(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, userID)
}
