// TODO check pool load on add
// TODO ignore to frequent progress updates
// TODO context
// TODO gracefully stop pool
// TODO give get status a timeout
package tasks

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/alitto/pond"
)

type State int

const (
	Waiting State = iota
	Running
	Done
	Failed
)

type Progress struct {
	Numerator, Denominator int
}

type Status struct {
	State     State
	Progress  Progress
	Error     error
	StartTime time.Time
	EndTime   time.Time
}

type Task struct {
	id  string
	hub *Hub
}

type taskCmd struct {
	id string
	fn func(Task) error
}

type stateCmd struct {
	id    string
	state State
	err   error
	start time.Time
	end   time.Time
}

type progressCmd struct {
	id         string
	num, denom int
}

type statusCmd struct {
	id string
	ch chan Status
}

type Hub struct {
	statuses         map[string]*Status
	pool             *pond.WorkerPool
	mu               sync.RWMutex
	minRetentionTime time.Duration
}

func NewHub() *Hub {
	return &Hub{
		statuses:         make(map[string]*Status),
		pool:             pond.New(250, 5000),
		minRetentionTime: 5 * time.Minute,
	}
}

func (h *Hub) Add(id string, fn func(Task) error) {
	if h.addTask(id) {
		h.pool.Submit(func() {
			h.setRunning(id)
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch rt := r.(type) {
					case string:
						err = errors.New(rt)
					case error:
						err = rt
					default:
						err = errors.New("unknown panic")
					}
					h.setFailed(id, err)
				}
			}()
			if err := fn(Task{hub: h, id: id}); err != nil {
				h.setFailed(id, err)
				return
			}
			h.setDone(id)
		})
	}
}

func (h *Hub) Status(id string) Status {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if s, ok := h.statuses[id]; ok {
		return *s
	}
	return Status{}
}

func (h *Hub) addTask(id string) bool {
	h.cleanup()
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.statuses[id]; ok {
		return false
	}
	h.statuses[id] = &Status{}
	return true
}

func (h *Hub) setRunning(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.statuses[id]
	s.State = Running
	s.StartTime = time.Now()
}

func (h *Hub) setDone(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.statuses[id]
	s.State = Done
	s.EndTime = time.Now()
}

func (h *Hub) setFailed(id string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.statuses[id]
	s.State = Failed
	s.EndTime = time.Now()
	s.Error = err
}

func (h *Hub) setProgress(id string, num, denom int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.statuses[id]
	s.Progress.Numerator = num
	s.Progress.Denominator = denom
}

func (h *Hub) cleanup() {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	for id, status := range h.statuses {
		if (status.State == Done || status.State == Failed) && now.Sub(status.EndTime) > h.minRetentionTime {
			delete(h.statuses, id)
		}
	}
}

func (s Status) Waiting() bool {
	return s.State == Waiting
}

func (s Status) Running() bool {
	return s.State == Running
}

func (s Status) Done() bool {
	return s.State == Done
}

func (s Status) Failed() bool {
	return s.State == Failed
}

func (p Progress) Percent() int {
	if p.Denominator == 0 {
		return 0
	}
	return int(math.Round((float64(p.Numerator) * float64(100)) / float64(p.Denominator)))
}

func (t Task) Progress(num int, denom int) {
	t.hub.setProgress(t.id, num, denom)
}
