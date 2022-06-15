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

type mutator[T, M any] struct {
	entityType *entityType[T]
	name       string
	fn         func(T, M) (T, error)
}

func NewMutator[T, M any](t *entityType[T], name string, fn func(T, M) (T, error)) *mutator[T, M] {
	return &mutator[T, M]{
		entityType: t,
		name:       name,
		fn:         fn,
	}
}

func (m *mutator[T, M]) EntityName() string {
	return m.entityType.name
}

func (m *mutator[T, M]) Name() string {
	return m.name
}

func (m *mutator[T, M]) Apply(d, dd any) (any, error) {
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

func (m *mutator[T, M]) convert(d any) (data M, err error) {
	switch t := d.(type) {
	case nil:
		// do nothing
	case M:
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

func (h *mutator[T, M]) New(data M, meta ...Meta) *mutation[T, M] {
	e := &mutation[T, M]{
		data:    data,
		mutator: h,
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
	EntityType() EntityType
	Name() string
	Data() any
	Meta() Meta
	Apply(any) (any, error)
}

type mutation[T, M any] struct {
	data    M
	meta    Meta
	mutator *mutator[T, M]
}

func (m *mutation[T, M]) EntityType() EntityType {
	return m.mutator.entityType
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

func (m *mutation[T, M]) Apply(d any) (any, error) {
	entityData, err := m.mutator.entityType.convert(d)
	if err != nil {
		return entityData, err
	}

	return m.mutator.fn(entityData, m.data)
}
