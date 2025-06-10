package storage

import (
	"sync"

	"github.com/paxaf/workmateTest/internal/entity"
)

type HashTable struct {
	mutex sync.RWMutex
	data  map[string]entity.Task
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]entity.Task),
	}
}

func (h *HashTable) Set(key string, value entity.Task) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.data[key] = value
}

func (h *HashTable) Del(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.data, key)
}

func (h *HashTable) Get(key string) (entity.Task, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	value, found := h.data[key]
	return value, found
}
