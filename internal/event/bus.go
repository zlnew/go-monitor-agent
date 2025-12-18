// Package event
package event

import "sync"

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]func(event any)
}

func New() *Bus {
	return &Bus{
		handlers: make(map[string][]func(event any)),
	}
}

func (b *Bus) Subscribe(eventName string, handler func(event any)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *Bus) Publish(eventName string, event any) {
	b.mu.RLock()
	handlers := b.handlers[eventName]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}
