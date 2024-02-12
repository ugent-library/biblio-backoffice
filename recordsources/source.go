package recordsources

import (
	"context"
	"fmt"
	"sync"
)

type Source interface {
	GetRecords(context.Context, func(Record) error) error
}

type Record struct {
	SourceName     string
	SourceID       string
	SourceMetadata []byte
}

type Factory func(string) (Source, error)

var factories = make(map[string]Factory)
var mu sync.RWMutex

func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	factories[name] = factory
}

func New(name, conn string) (Source, error) {
	mu.RLock()
	factory, ok := factories[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown source '%s'", name)
	}
	return factory(conn)
}
