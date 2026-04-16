// Package multicast fans out a single vault event to multiple handlers.
package multicast

import "sync"

// Handler is a function that receives a named event payload.
type Handler func(event string, payload any)

// Broadcaster delivers events to a registered set of handlers concurrently.
type Broadcaster struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New returns an empty Broadcaster.
func New() *Broadcaster {
	return &Broadcaster{handlers: make(map[string]Handler)}
}

// Register adds a named handler. Registering the same name overwrites the
// previous handler.
func (b *Broadcaster) Register(name string, h Handler) {
	if name == "" || h == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
}

// Deregister removes a handler by name. It is a no-op if the name is unknown.
func (b *Broadcaster) Deregister(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, name)
}

// Len returns the number of registered handlers.
func (b *Broadcaster) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}

// Broadcast delivers event and payload to every registered handler in
// separate goroutines and waits for all to finish.
func (b *Broadcaster) Broadcast(event string, payload any) {
	b.mu.RLock()
	snap := make(map[string]Handler, len(b.handlers))
	for k, v := range b.handlers {
		snap[k] = v
	}
	b.mu.RUnlock()

	var wg sync.WaitGroup
	for _, h := range snap {
		wg.Add(1)
		go func(fn Handler) {
			defer wg.Done()
			fn(event, payload)
		}(h)
	}
	wg.Wait()
}
