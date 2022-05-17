// TODO check pool load on add
// TODO ignore to frequent progress updates
// TODO context
// TODO gracefully stop pool
// TODO give get status a timeout
package tasks

import (
	"errors"
	"math"
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
	id         string
	progressCh chan progressCmd
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
	statuses   map[string]Status
	pool       *pond.WorkerPool
	taskCh     chan taskCmd
	stateCh    chan stateCmd
	progressCh chan progressCmd
	statusCh   chan statusCmd
}

func NewHub() *Hub {
	return &Hub{
		statuses:   make(map[string]Status),
		pool:       pond.New(250, 5000),
		taskCh:     make(chan taskCmd),
		stateCh:    make(chan stateCmd),
		progressCh: make(chan progressCmd),
		statusCh:   make(chan statusCmd),
	}
}

func (h *Hub) Add(id string, fn func(Task) error) string {
	h.taskCh <- taskCmd{id: id, fn: fn}
	return id
}

func (h *Hub) Status(id string) Status {
	ch := make(chan Status)
	defer close(ch)
	h.statusCh <- statusCmd{id: id, ch: ch}
	return <-ch
}

func (h *Hub) addTaskToPool(id string, fn func(Task) error) {
	h.pool.Submit(func() {
		h.stateCh <- stateCmd{id: id, state: Running, start: time.Now()}
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
				h.stateCh <- stateCmd{id: id, state: Failed, err: err, end: time.Now()}
			}
		}()
		if err := fn(Task{id: id, progressCh: h.progressCh}); err != nil {
			h.stateCh <- stateCmd{id: id, state: Failed, err: err, end: time.Now()}
			return
		}
		h.stateCh <- stateCmd{id: id, state: Done, end: time.Now()}
	})
}

func (h *Hub) Run() {
	retentionTime := 5 * time.Minute
	cleanupTicker := time.NewTicker(retentionTime)
	defer cleanupTicker.Stop()

	for {
		select {
		case now := <-cleanupTicker.C:
			for id, status := range h.statuses {
				if (status.State == Done || status.State == Failed) && now.Sub(status.EndTime) > retentionTime {
					delete(h.statuses, id)
				}
			}
		case cmd := <-h.taskCh:
			// no overlapping tasks
			if _, ok := h.statuses[cmd.id]; !ok {
				h.statuses[cmd.id] = Status{}
				h.addTaskToPool(cmd.id, cmd.fn)
			}
		case cmd := <-h.progressCh:
			if s, ok := h.statuses[cmd.id]; ok {
				s.Progress.Numerator = cmd.num
				s.Progress.Denominator = cmd.denom
			}
		case cmd := <-h.stateCh:
			if s, ok := h.statuses[cmd.id]; ok {
				s.State = cmd.state
				s.Error = cmd.err
				if !cmd.start.IsZero() {
					s.StartTime = cmd.start
				}
				if !cmd.end.IsZero() {
					s.EndTime = cmd.end
				}
			}
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
	t.progressCh <- progressCmd{id: t.id, num: num, denom: denom}
}
