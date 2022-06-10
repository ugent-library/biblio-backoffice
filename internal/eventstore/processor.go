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

type processor[T any] struct {
	handlers   map[string]func(T, any) (T, error)
	handlersMu sync.RWMutex
}

func NewProcessor[T any]() *processor[T] {
	return &processor[T]{}
}

func (p *processor[T]) AddHandler(eventType string, h func(T, any) (T, error)) {
	p.handlersMu.Lock()
	defer p.handlersMu.Unlock()
	if p.handlers == nil {
		p.handlers = make(map[string]func(T, any) (T, error))
	}
	p.handlers[eventType] = h
}

func (p *processor[T]) RawApply(d json.RawMessage, events []Event) (json.RawMessage, error) {
	var (
		data T
		err  error
	)

	if d != nil {
		if err = json.Unmarshal(d, data); err != nil {
			return nil, fmt.Errorf("eventstore: failed to deserialize into %T: %w", data, err)
		}
	}

	if data, err = p.Apply(data, events...); err != nil {
		return nil, err
	}

	if d, err = json.Marshal(data); err != nil {
		return nil, fmt.Errorf("eventstore: failed to serialize %T: %w", data, err)
	}

	return d, nil
}

func (p *processor[T]) Apply(data T, events ...Event) (T, error) {
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
