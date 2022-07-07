package mutantdb

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Listener[T any] struct {
	conn       *pgx.Conn
	entityType *Type[T]
	listening  bool
	chunks     map[string]*bytes.Buffer
}

func NewListener[T any](conn *pgx.Conn, t *Type[T]) *Listener[T] {
	return &Listener[T]{
		conn:       conn,
		entityType: t,
		chunks:     make(map[string]*bytes.Buffer),
	}
}

func (l *Listener[T]) WaitForProjection(ctx context.Context) (Projection[T], error) {
	p := Projection[T]{}

	if !l.listening {
		chanName := l.entityType.Name() + "_projections"
		if _, err := l.conn.Exec(ctx, "listen $1", chanName); err != nil {
			return p, err
		}
		l.listening = true
	}

	for {
		// if ctx is done, err will be non nil
		n, err := l.conn.WaitForNotification(ctx)
		if err != nil {
			return p, err
		}

		// characters before first pipe are notification id
		// characters between first and second pipe are chunk counter
		// characters after second pipe are up to 4000 bytes of json
		// payload is complete if counter is EOF
		pipe1 := strings.Index(n.Payload, "|")
		pipe2 := strings.Index(n.Payload[pipe1+1:], "|") + pipe1 + 1
		id := n.Payload[:pipe1]
		counter := n.Payload[pipe1+1 : pipe2]
		chunk := n.Payload[pipe2+1:]

		buf, ok := l.chunks[id]
		if !ok {
			buf = bytes.NewBuffer([]byte{})
			l.chunks[id] = buf
		}

		if counter != "EOF" {
			buf.WriteString(chunk)
			continue
		}

		delete(l.chunks, id)

		err = json.Unmarshal(buf.Bytes(), &p)

		return p, err
	}
}
