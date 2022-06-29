package ulid

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

var defaultGenerator = NewGenerator()

type generator struct {
	mu sync.Mutex
	r  io.Reader
}

func Generate() (string, error) {
	return defaultGenerator.Generate()
}

func MustGenerate() string {
	return defaultGenerator.MustGenerate()
}

func NewGenerator() *generator {
	return &generator{
		r: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
}

func (g *generator) Generate() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id, err := ulid.New(ulid.Now(), g.r)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (g *generator) MustGenerate() string {
	id, err := g.Generate()
	if err != nil {
		panic(err)
	}
	return id
}
