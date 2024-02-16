package recordsources

import (
	"context"
	"fmt"
	"sync"

	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/models"
)

type Source interface {
	GetRecords(context.Context, func(Record) error) error
}

type Record interface {
	SourceName() string
	SourceID() string
	ToCandidateRecord(*backends.Services) (*models.CandidateRecord, error)
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
