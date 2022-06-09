package eventstore

import (
	"encoding/json"
	"fmt"
	"sync"
)

func Handler[T, TT any](fn func(T, TT) (T, error)) func(T, any) (T, error) {
	return func(t T, a any) (T, error) {
		return fn(t, a.(TT))
	}
}

type Processor interface {
	RawApply(json.RawMessage, []Event) (any, error)
}

type ProcessorFor[T any] struct {
	handlers   map[string]func(T, any) (T, error)
	handlersMu sync.RWMutex
}

func (p *ProcessorFor[T]) AddHandler(eventType string, handler func(T, any) (T, error)) {
	p.handlersMu.Lock()
	defer p.handlersMu.Unlock()
	if p.handlers == nil {
		p.handlers = make(map[string]func(T, any) (T, error))
	}
	p.handlers[eventType] = handler
}

func (p *ProcessorFor[T]) RawApply(d json.RawMessage, events []Event) (any, error) {
	var data T
	if d != nil {
		if err := json.Unmarshal(d, data); err != nil {
			return nil, fmt.Errorf("eventstore: can't deserialize into %T: %w", data, err)
		}
	}

	return p.Apply(data, events...)
}

func (p *ProcessorFor[T]) Apply(data T, events ...Event) (T, error) {
	var err error

	for _, e := range events {
		p.handlersMu.RLock()
		handler, ok := p.handlers[e.Type]
		p.handlersMu.RUnlock()

		if !ok {
			return data, fmt.Errorf("eventstore: no handler for %s event %s", e.StreamType, e.Type)
		}

		data, err = handler(data, e.Data)

		if err != nil {
			return data, fmt.Errorf("eventstore: failed to apply %s event %s: %w", e.StreamType, e.Type, err)
		}
	}

	return data, err
}
