package mutantdb

import (
	"encoding/json"
	"fmt"
)

type EntityType interface {
	Name() string
	New() any
}

type entityType[T any] struct {
	name    string
	factory func() T
}

func NewType[T any](name string, factory func() T) *entityType[T] {
	return &entityType[T]{
		name:    name,
		factory: factory,
	}
}

func (s *entityType[T]) Name() string {
	return s.name
}

func (s *entityType[T]) New() any {
	return s.factory()
}

func (s *entityType[T]) convert(d any) (data T, err error) {
	switch t := d.(type) {
	case T:
		data = t
	case json.RawMessage:
		if err = json.Unmarshal(t, &data); err != nil {
			err = fmt.Errorf("mutantdb: failed to deserialize entity data into %T: %w", data, err)
		}
	default:
		err = fmt.Errorf("mutantdb: invalid entity data type %T", t)
	}

	return
}
