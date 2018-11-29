package httpproxy

import (
	"sync"

	"github.com/mz-eco/http"
)

type Source struct {
	lock     sync.Mutex
	size     int
	contexts []*http.Translate
}

func (m *Source) Add(tx *http.Translate) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.contexts = append(m.contexts, tx)
	m.size++
}

func NewSource() *Source {
	return &Source{
		contexts: make([]*http.Translate, 0),
	}
}

func (m *Source) GetSize() int {
	return m.size
}
