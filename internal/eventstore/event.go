package eventstore

import (
	"encoding/json"
	"fmt"
)

type Meta map[string]string

type Event interface {
	StreamID() string
	StreamType() string
	Name() string
	Data() any
	Meta() Meta
	Apply(any) (any, error)
}

type event[T, TT any] struct {
	streamID string
	data     TT
	meta     Meta
	handler  *eventHandler[T, TT]
}

func (e *event[T, TT]) StreamID() string {
	return e.streamID
}

func (e *event[T, TT]) StreamType() string {
	return e.handler.streamType
}

func (e *event[T, TT]) Name() string {
	return e.handler.name
}

func (e *event[T, TT]) Data() any {
	return e.data
}

func (e *event[T, TT]) Meta() Meta {
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
		if err := json.Unmarshal(t, &data); err != nil {
			return data, fmt.Errorf("eventstore: failed to deserialize projection data into %T: %w", data, err)
		}
	default:
		return data, fmt.Errorf("eventstore: invalid projection data type %T", t)
	}

	return e.handler.fn(data, e.data)
}
