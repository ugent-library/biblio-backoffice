package task

import (
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rabbitmq/amqp091-go"
)

type Task struct {
	ID                     string
	UserID                 string
	Numerator, Denominator int
	Error                  error
	Result                 interface{}
}

func (t Task) Running() bool {
	return t.Error == nil && t.Denominator > 0 && t.Numerator != t.Denominator
}

func (t Task) Done() bool {
	return t.Denominator > 0 && t.Numerator == t.Denominator
}

func (t Task) Percent() int {
	if t.Denominator == 0 {
		return 0
	}
	return int(math.Round((float64(t.Numerator) * float64(100)) / float64(t.Denominator)))
}

type Hub struct {
	mqCh  *amqp091.Channel
	tasks map[string]Task
	sync.RWMutex
}

func NewHub(mqCh *amqp091.Channel) *Hub {
	return &Hub{
		mqCh:  mqCh,
		tasks: make(map[string]Task),
	}
}

func (h *Hub) Get(id string) (Task, bool) {
	h.RLock()
	defer h.RUnlock()
	t, ok := h.tasks[id]
	return t, ok
}

func (h *Hub) Add(key, userID string, payload interface{}) (string, error) {
	id := newCorrelationID()

	msg := struct {
		CorrelationID string      `json:"correlation_id"`
		UserID        string      `json:"user_id"`
		Payload       interface{} `json:"payload"`
	}{
		id,
		userID,
		payload,
	}

	msgJSON, _ := json.Marshal(msg)

	err := h.mqCh.Publish(
		"tasks",      // exchange
		"tasks."+key, // routing key
		false,        // mandatory
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         msgJSON,
		},
	)
	if err != nil {
		return "", err
	}

	h.Lock()
	defer h.Unlock()
	h.tasks[id] = Task{ID: id, UserID: userID}

	return id, nil
}

func (h *Hub) SetError(id string, e error) {
	h.Lock()
	defer h.Unlock()
	if task, ok := h.tasks[id]; ok {
		task.Error = e
		h.tasks[id] = task
	}
}

func (h *Hub) SetProgress(id string, n, d int) {
	h.Lock()
	defer h.Unlock()
	if task, ok := h.tasks[id]; ok {
		task.Numerator = n
		task.Denominator = d
		h.tasks[id] = task
	}
}

func newCorrelationID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
