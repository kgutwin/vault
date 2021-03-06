package physical

import (
	"fmt"
	"sync"
)

type InmemHABackend struct {
	InmemBackend
	locks map[string]string
	l     sync.Mutex
	cond  *sync.Cond
}

// NewInmemHA constructs a new in-memory HA backend. This is only for testing.
func NewInmemHA() *InmemHABackend {
	in := &InmemHABackend{
		InmemBackend: *NewInmem(),
		locks:        make(map[string]string),
	}
	in.cond = sync.NewCond(&in.l)
	return in
}

// LockWith is used for mutual exclusion based on the given key.
func (i *InmemHABackend) LockWith(key, value string) (Lock, error) {
	l := &InmemLock{
		in:    i,
		key:   key,
		value: value,
	}
	return l, nil
}

// InmemLock is an in-memory Lock implementation for the HABackend
type InmemLock struct {
	in    *InmemHABackend
	key   string
	value string

	held     bool
	leaderCh chan struct{}
	l        sync.Mutex
}

func (i *InmemLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	i.l.Lock()
	defer i.l.Unlock()
	if i.held {
		return nil, fmt.Errorf("lock already held")
	}

	// Attempt an async acquisition
	didLock := make(chan struct{})
	releaseCh := make(chan bool, 1)
	go func() {
		// Wait to acquire the lock
		i.in.l.Lock()
		_, ok := i.in.locks[i.key]
		for ok {
			i.in.cond.Wait()
			_, ok = i.in.locks[i.key]
		}
		i.in.locks[i.key] = i.value
		i.in.l.Unlock()

		// Signal that lock is held
		close(didLock)

		// Handle an early abort
		release := <-releaseCh
		if release {
			i.in.l.Lock()
			delete(i.in.locks, i.key)
			i.in.l.Unlock()
			i.in.cond.Broadcast()
		}
	}()

	// Wait for lock acquisition or shutdown
	select {
	case <-didLock:
		releaseCh <- false
	case <-stopCh:
		releaseCh <- true
		return nil, nil
	}

	// Create the leader channel
	i.held = true
	i.leaderCh = make(chan struct{})
	return i.leaderCh, nil
}

func (i *InmemLock) Unlock() error {
	i.l.Lock()
	defer i.l.Unlock()

	if !i.held {
		return nil
	}

	close(i.leaderCh)
	i.leaderCh = nil
	i.held = false

	i.in.l.Lock()
	delete(i.in.locks, i.key)
	i.in.l.Unlock()
	i.in.cond.Broadcast()
	return nil
}

func (i *InmemLock) Value() (bool, string, error) {
	i.in.l.Lock()
	val, ok := i.in.locks[i.key]
	i.in.l.Unlock()
	return ok, val, nil
}
