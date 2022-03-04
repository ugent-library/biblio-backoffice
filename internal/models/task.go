package models

import "math"

type TaskStatus int

const (
	Waiting TaskStatus = iota
	Running
	Done
	Failed
)

type TaskState struct {
	Status      TaskStatus
	Message     string
	Numerator   int
	Denominator int
}

func (t TaskState) Waiting() bool {
	return t.Status == Waiting
}

func (t TaskState) Running() bool {
	return t.Status == Running
}

func (t TaskState) Done() bool {
	return t.Status == Done
}

func (t TaskState) Failed() bool {
	return t.Status == Failed
}

func (t TaskState) Percent() int {
	if t.Denominator == 0 {
		return 0
	}
	return int(math.Round((float64(t.Numerator) * float64(100)) / float64(t.Denominator)))
}
