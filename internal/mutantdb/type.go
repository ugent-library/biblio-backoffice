package mutantdb

type Type[T any] struct {
	name      string
	factory   func() T
	validator func(T) error
}

func NewType[T any](name string, factory func() T) *Type[T] {
	return &Type[T]{
		name:      name,
		factory:   factory,
		validator: func(T) error { return nil },
	}
}

func (t *Type[T]) WithValidator(fn func(T) error) *Type[T] {
	t.validator = fn
	return t
}

func (t *Type[T]) Name() string {
	return t.name
}

func (t *Type[T]) New() T {
	return t.factory()
}

func (t *Type[T]) Validate(d T) error {
	return t.validator(d)
}
