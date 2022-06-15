package mutantdb

type Type[T any] struct {
	name    string
	factory func() T
}

func NewType[T any](name string, factory func() T) *Type[T] {
	return &Type[T]{
		name:    name,
		factory: factory,
	}
}
func (s *Type[T]) Name() string {
	return s.name
}

func (s *Type[T]) New() T {
	return s.factory()
}
