package store

import (
	"sync"

	"weight-tracker/src/models"
)

// Store defines the operations used by handlers.
type Store interface {
	List() []models.WeightEntry
	Create(models.WeightEntry) models.WeightEntry
	Get(string) (models.WeightEntry, bool)
	Update(string, models.WeightEntry) (models.WeightEntry, bool)
	Delete(string) bool
	Reset()
}

type memoryStore struct {
	mu   sync.RWMutex
	data map[string]models.WeightEntry
}

func NewMemoryStore() Store {
	return &memoryStore{data: map[string]models.WeightEntry{}}
}

var Default Store = NewMemoryStore()

func (m *memoryStore) List() []models.WeightEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.WeightEntry, 0, len(m.data))
	for _, v := range m.data {
		out = append(out, v)
	}
	return out
}

func (m *memoryStore) Create(entry models.WeightEntry) models.WeightEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[entry.ID] = entry
	return entry
}

func (m *memoryStore) Get(id string) (models.WeightEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.data[id]
	return e, ok
}

func (m *memoryStore) Update(id string, entry models.WeightEntry) (models.WeightEntry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[id]; !ok {
		return models.WeightEntry{}, false
	}
	entry.ID = id
	m.data[id] = entry
	return entry, true
}

func (m *memoryStore) Delete(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[id]; !ok {
		return false
	}
	delete(m.data, id)
	return true
}

func (m *memoryStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = map[string]models.WeightEntry{}
}
