package mutantdb

// TODO add (de)serialization

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
