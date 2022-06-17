package mutantdb

import (
	"encoding/json"
	"fmt"
)

type Mutator[T any] interface {
	Name() string
	Apply(T, any) (T, error)
}

type mutator[T, M any] struct {
	name string
	fn   func(T, M) (T, error)
}

func NewMutator[T, M any](name string, fn func(T, M) (T, error)) *mutator[T, M] {
	return &mutator[T, M]{
		name: name,
		fn:   fn,
	}
}

func (m *mutator[T, M]) Name() string {
	return m.name
}

func (m *mutator[T, M]) Apply(data T, md any) (T, error) {
	var mutationData M

	switch t := md.(type) {
	case nil:
		// do nothing
	case M:
		mutationData = t
	case json.RawMessage:
		if err := json.Unmarshal(t, &mutationData); err != nil {
			return data, fmt.Errorf("mutantdb: failed to deserialize mutation data into %T: %w", mutationData, err)
		}
	default:
		return data, fmt.Errorf("mutantdb: invalid mutation data type %T", t)
	}

	return m.fn(data, mutationData)
}

func (m *mutator[T, M]) New(data M, meta ...Meta) *mutation[T, M] {
	e := &mutation[T, M]{
		data:    data,
		mutator: m,
	}
	for _, meta := range meta {
		if e.meta == nil {
			e.meta = make(Meta)
		}
		for k, v := range meta {
			e.meta[k] = v
		}
	}
	return e
}

type Meta map[string]string

type Mutation[T any] interface {
	Name() string
	Data() any
	Meta() Meta
	Apply(T) (T, error)
}

// Apply is a convenience function that applies mutations to a value.
func Apply[T any](d T, mutations ...Mutation[T]) (T, error) {
	for _, m := range mutations {
		d, err := m.Apply(d)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

type mutation[T, M any] struct {
	data    M
	meta    Meta
	mutator *mutator[T, M]
}

func (m *mutation[T, M]) Name() string {
	return m.mutator.name
}

func (m *mutation[T, M]) Data() any {
	return m.data
}

func (m *mutation[T, M]) Meta() Meta {
	return m.meta
}

func (m *mutation[T, M]) Apply(data T) (T, error) {
	return m.mutator.fn(data, m.data)
}
