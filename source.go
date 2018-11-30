package httpproxy

import (
	"sync"

	"github.com/mz-eco/http"
)

type Watcher interface {
	OnListen(txs []*http.Summary, done chan Watcher)
	OnTranslate(tx *http.Summary)
}

type Listener struct {
	s     *Source
	w     Watcher
	done  chan Watcher
	event chan interface{}
}

func (m *Listener) loop() {

	for {
		select {
		case w := <-m.done:
			m.s.removeWatcher(w)
			break
		case e := <-m.event:

			switch x := e.(type) {
			case []*http.Summary:
				m.w.OnListen(x, m.done)
			case *http.Summary:
				m.w.OnTranslate(x)
			}
		}
	}
}

type Source struct {
	lock     sync.Mutex
	size     int
	current  int
	txs      []*http.Translate
	watcher  map[Watcher]*Listener
	done     chan Watcher
	listener chan *http.Translate
}

func (m *Source) loop() {
	for {

		select {
		case w := <-m.done:
			m.removeWatcher(w)
		case x := <-m.listener:
			m.dispatch(x)
		}
	}
}

func (m *Source) GetTranslate(index int) *http.Translate {
	return m.txs[index-m.current]
}

func (m *Source) removeWatcher(w Watcher) {
	m.lock.Lock()
	m.lock.Unlock()

	delete(m.watcher, w)
}

func (m *Source) AddWatcher(w Watcher) {

	m.lock.Lock()
	defer m.lock.Unlock()

	l := &Listener{
		s:     m,
		w:     w,
		done:  make(chan Watcher),
		event: make(chan interface{}),
	}

	l.event <- m.Summary(50)
	m.watcher[w] = l

	go l.loop()
}

func (m *Source) Summary(limit int) []*http.Summary {

	m.lock.Lock()
	defer m.lock.Unlock()

	var (
		summaryList = make([]*http.Summary, 0)
		x           = m.txs
	)

	if limit > len(x) {
		x = x[limit-len(x):]
	}

	for index, tx := range x {

		su := tx.Summary()
		su.Index = index + m.current

		summaryList = append(summaryList, su)
	}

	return summaryList

}

func (m *Source) dispatch(tx *http.Translate) {

	s := tx.Summary()

	for _, l := range m.watcher {
		l.event <- s
	}

}

func (m *Source) Add(tx *http.Translate) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.listener <- tx
	m.txs = append(m.txs, tx)
	m.size++

}

func NewSource() *Source {
	source := &Source{
		txs:     make([]*http.Translate, 0),
		done:    make(chan Watcher),
		watcher: make(map[Watcher]*Listener, 0),
	}

	go source.loop()

	return source
}

func (m *Source) GetSize() int {
	return m.size
}
