package cmutex

import (
	"context"
	"sync/atomic"
)

// CMutex is a context-aware synchronization primitive.
type CMutex struct {
	c  int32
	ch lazyChan
}

// Lock attempts to lock the mutex.
// Tis method will block until the requested resource is available or
// until the context expires.
//
// If the take is not successful this method will return an error.
func (m *CMutex) Lock(ctx context.Context) error {

	// increment the counter on the mutex
	v := atomic.AddInt32(&m.c, 1)
	if v == 1 {
		return nil
	}
	select {
	case <-m.ch.get():
		// Lock grabbed.
		return nil
	case <-ctx.Done():
		go func() {
			// drain the chan to ensure we don't block
			<-m.ch.get()
			m.Unlock()
		}()
		return ctx.Err()
	}
}

// Unlock the mutex.
// Similar to sync.CMutex, this method will panic if
// the mutex is already unlocked.
func (m *CMutex) Unlock() {
	v := atomic.AddInt32(&m.c, -1)
	if v < 0 {
		panic("unlock of an already unlocked mutex")
	}

	if v > 0 {
		m.ch.get() <- struct{}{}
	}
}
