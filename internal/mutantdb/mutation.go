package mutantdb

import (
	"encoding/json"
	"fmt"
)

type Mutator interface {
	EntityName() string
	Name() string
	Apply(any, any) (any, error)
}

type mutator[T, TT any] struct {
	entityType *entityType[T]
	name       string
	fn         func(T, TT) (T, error)
}

func NewMutator[T, TT any](t *entityType[T], name string, fn func(T, TT) (T, error)) *mutator[T, TT] {
	return &mutator[T, TT]{
		entityType: t,
		name:       name,
		fn:         fn,
	}
}

func (m *mutator[T, TT]) EntityName() string {
	return m.entityType.name
}

func (m *mutator[T, TT]) Name() string {
	return m.name
}

func (m *mutator[T, TT]) Apply(d, dd any) (any, error) {
	entityData, err := m.entityType.convert(d)
	if err != nil {
		return entityData, err
	}

	mutationData, err := m.convert(dd)
	if err != nil {
		return entityData, err
	}

	return m.fn(entityData, mutationData)
}

func (m *mutator[T, TT]) convert(d any) (data TT, err error) {
	switch t := d.(type) {
	case nil:
		// do nothing
	case TT:
		data = t
	case json.RawMessage:
		if err = json.Unmarshal(t, &data); err != nil {
			err = fmt.Errorf("mutantdb: failed to deserialize mutation data into %T: %w", data, err)
		}
	default:
		err = fmt.Errorf("mutantdb: invalid mutation data type %T", t)
	}

	return
}

func (h *mutator[T, TT]) New(entityID string, data TT, meta ...Meta) *mutation[T, TT] {
	e := &mutation[T, TT]{
		entityID: entityID,
		data:     data,
		mutator:  h,
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

type Mutation interface {
	EntityID() string
	EntityType() EntityType
	Name() string
	Data() any
	Meta() Meta
	Apply(any) (any, error)
}

type mutation[T, TT any] struct {
	entityID string
	data     TT
	meta     Meta
	mutator  *mutator[T, TT]
}

func (m *mutation[T, TT]) EntityID() string {
	return m.entityID
}

func (m *mutation[T, TT]) EntityType() EntityType {
	return m.mutator.entityType
}

func (m *mutation[T, TT]) Name() string {
	return m.mutator.name
}

func (m *mutation[T, TT]) Data() any {
	return m.data
}

func (m *mutation[T, TT]) Meta() Meta {
	return m.meta
}

func (m *mutation[T, TT]) Apply(d any) (any, error) {
	entityData, err := m.mutator.entityType.convert(d)
	if err != nil {
		return entityData, err
	}

	return m.mutator.fn(entityData, m.data)
}
