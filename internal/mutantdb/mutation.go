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

func (h *mutator[T, TT]) EntityName() string {
	return h.entityType.name
}

func (h *mutator[T, TT]) Name() string {
	return h.name
}

func (h *mutator[T, TT]) Apply(d, ed any) (any, error) {
	var (
		data         T
		mutationData TT
	)

	switch t := d.(type) {
	case T:
		data = t
	case json.RawMessage:
		if err := json.Unmarshal(t, &data); err != nil {
			return data, fmt.Errorf("mutantdb: failed to deserialize projection data into %T: %w", data, err)
		}
	default:
		return data, fmt.Errorf("mutantdb: invalid projection data type %T", t)
	}

	switch t := ed.(type) {
	case nil:
		// do nothing
	case TT:
		mutationData = t
	case json.RawMessage:
		if err := json.Unmarshal(t, &mutationData); err != nil {
			return data, fmt.Errorf("mutantdb: failed to deserialize mutation data into %T: %w", mutationData, err)
		}
	default:
		return data, fmt.Errorf("mutantdb: invalid mutation data type %T", t)
	}

	return h.fn(data, mutationData)
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

func (e *mutation[T, TT]) EntityID() string {
	return e.entityID
}

func (e *mutation[T, TT]) EntityType() EntityType {
	return e.mutator.entityType
}

func (e *mutation[T, TT]) Name() string {
	return e.mutator.name
}

func (e *mutation[T, TT]) Data() any {
	return e.data
}

func (e *mutation[T, TT]) Meta() Meta {
	return e.meta
}

func (e *mutation[T, TT]) Apply(d any) (any, error) {
	var data T

	switch t := d.(type) {
	case T:
		data = t
	case json.RawMessage:
		if err := json.Unmarshal(t, &data); err != nil {
			return data, fmt.Errorf("mutantdb: failed to deserialize projection data into %T: %w", data, err)
		}
	default:
		return data, fmt.Errorf("mutantdb: invalid projection data type %T", t)
	}

	return e.mutator.fn(data, e.data)
}
