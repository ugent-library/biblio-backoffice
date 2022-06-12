package eventstore

import (
	"encoding/json"
	"fmt"
)

type EventHandler[T any] interface {
	Name() string
	// Apply(T, any) (T, error)
}

type eventHandler[T, TT any] struct {
	name string
	fn   func(T, TT) (T, error)
}

func NewEventHandler[T, TT any](name string, fn func(T, TT) (T, error)) *eventHandler[T, TT] {
	return &eventHandler[T, TT]{
		name: name,
		fn:   fn,
	}
}

func (h *eventHandler[T, TT]) Name() string {
	return h.name
}

func (h *eventHandler[T, TT]) Apply(data T, d any) (T, error) {
	var eventData TT

	switch t := d.(type) {
	case nil:
		// do nothing
	case TT:
		eventData = t
	case json.RawMessage:
		if err := json.Unmarshal(t, eventData); err != nil {
			return data, fmt.Errorf("eventstore: failed to deserialize event data into %T: %w", eventData, err)
		}
	default:
		return data, fmt.Errorf("eventstore: invalid event data type %T", t)
	}

	return h.fn(data, eventData)
}

func (h *eventHandler[T, TT]) NewEvent(streamType, streamID string, data TT, meta ...map[string]string) *event[T, TT] {
	e := &event[T, TT]{
		streamType: streamType,
		streamID:   streamID,
		data:       data,
		handler:    h,
	}
	for _, meta := range meta {
		if e.meta == nil {
			e.meta = meta
		} else {
			for k, v := range meta {
				e.meta[k] = v
			}
		}
	}
	return e
}

type Event interface {
	StreamID() string
	StreamType() string
	Name() string
	Data() any
	Meta() map[string]string
	Apply(any) (any, error)
}

type event[T, TT any] struct {
	streamID   string
	streamType string
	data       TT
	meta       map[string]string
	handler    *eventHandler[T, TT]
}

func (e *event[T, TT]) StreamID() string {
	return e.streamID
}

func (e *event[T, TT]) StreamType() string {
	return e.streamType
}

func (e *event[T, TT]) Name() string {
	return e.handler.name
}

func (e *event[T, TT]) Data() any {
	return e.data
}

func (e *event[T, TT]) Meta() map[string]string {
	return e.meta
}

func (e *event[T, TT]) Apply(d any) (any, error) {
	var data T

	switch t := d.(type) {
	case nil:
		// TODO remove this when we have a factory for nil data
	case T:
		data = t
	case json.RawMessage:
		if err := json.Unmarshal(t, data); err != nil {
			return data, fmt.Errorf("eventstore: failed to deserialize projection data into %T: %w", data, err)
		}
	default:
		return data, fmt.Errorf("eventstore: invalid projection data type %T", t)
	}

	return e.handler.fn(data, e.data)
}
