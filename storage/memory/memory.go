package memory

import (
	"crawler/storage"
	"errors"
	"net/url"
	"sync"
)

// Errors returned from memory storage.
var (
	ErrDuplicateKey = errors.New("key already exists")
	ErrInvalidKey   = errors.New("key does not exist")
)

// Memory is a storage method which stores storage.Page's in-memory.
type Memory struct {
	pages map[url.URL]storage.Page
	sync.RWMutex
}

// New instantiates Memory.
func New() *Memory {
	return &Memory{
		pages: make(map[url.URL]storage.Page),
	}
}

// Get fetches a storage.Page for a given url.URL. Error if no key can be found.
func (m *Memory) Get(u url.URL) (storage.Page, error) {
	m.RLock()
	defer m.RUnlock()
	s, ok := m.pages[u]
	if !ok {
		return storage.Page{}, ErrInvalidKey
	}

	return s, nil
}

// Insert creates a storage.Page in-memory. Error if key already exists.
func (m *Memory) Insert(page storage.Page) error {
	m.Lock()
	defer m.Unlock()
	_, ok := m.pages[page.URL]
	if ok {
		return ErrDuplicateKey
	}

	m.pages[page.URL] = page

	return nil
}
