package cmutex

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
)

func TestMutexPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("unlock of unlocked mutex did not panic")
		}
	}()

	var mu Mutex
	mu.Lock(context.Background())
	mu.Unlock()
	mu.Unlock() //panic here
}

// borrowed from https://golang.org/src/sync/mutex_test.go
func hammerMutex(ctx context.Context, m *Mutex, loops int, cdone chan bool, canc *int32) {

	defer func() {
		// report as done
		cdone <- true
	}()

	for i := 0; i < loops; i++ {
		err := m.Lock(ctx)
		if err != nil {
			// mark as canceled
			atomic.AddInt32(canc, 1)
			return
		}
		m.Unlock()
	}
}

func TestMutex(t *testing.T) {
	m := &Mutex{}
	c := make(chan bool)
	ctx, cf := context.WithCancel(context.Background())

	var numCan int32
	for i := 0; i < 100; i++ {
		go hammerMutex(ctx, m, 5000, c, &numCan)
	}

	for i := 0; i < 100; i++ {
		<-c
		cf()
	}

	if numCan == 0 && runtime.GOMAXPROCS(0) > 1 {
		t.Fatal("no cancelations recorded")
	} else {
		t.Log("cancellations recored: ", numCan)
	}

	if m.c != 0 {
		t.Fatal("mutex should be unlocked after hammer")
	}

}
