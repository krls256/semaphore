package semaphore

import "sync"

type KeyAccessSemaphore struct {
	mu   sync.Mutex
	size int
	sem  map[string]chan struct{}
}

func NewKeyAccessSemaphore(size int) *KeyAccessSemaphore {
	if size <= 0 {
		panic("size must be greater than 1")
	}

	return &KeyAccessSemaphore{size: size, sem: map[string]chan struct{}{}}
}

func (sem *KeyAccessSemaphore) Lock(key string) {
	sem.mu.Lock()
	ch, ok := sem.sem[key]
	if !ok {
		ch = make(chan struct{}, sem.size)
		sem.sem[key] = ch
	}
	sem.mu.Unlock()

	ch <- struct{}{}
}

func (sem *KeyAccessSemaphore) Unlock(key string) {
	<-sem.sem[key]

	sem.mu.Lock()
	select {
	// chan is free
	case sem.sem[key] <- struct{}{}:
		delete(sem.sem, key)
	default: // nothing
	}
	sem.mu.Unlock()
}
