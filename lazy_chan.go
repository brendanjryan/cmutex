package cmutex

import "sync"

// lazychan is a simple wrapper around a generic channel
// which can be initialized lazily.
type lazyChan struct {
	init sync.Once
	ch   chan struct{}
}

func (s *lazyChan) get() chan struct{} {
	s.init.Do(func() {
		s.ch = make(chan struct{})
	})

	return s.ch
}
