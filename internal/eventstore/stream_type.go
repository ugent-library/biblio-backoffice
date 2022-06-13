package eventstore

type StreamType interface {
	Name() string
	New() any
}

type streamType[T any] struct {
	name    string
	factory func() T
}

func NewStreamType[T any](name string, factory func() T) *streamType[T] {
	return &streamType[T]{
		name:    name,
		factory: factory,
	}
}

func (s *streamType[T]) Name() string {
	return s.name
}

func (s *streamType[T]) New() any {
	return s.factory()
}
