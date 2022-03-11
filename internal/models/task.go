package models

import "math"

type TaskStatus int

const (
	Waiting TaskStatus = iota
	Running
	Done
	Failed
)

type Progress struct {
	Numerator, Denominator int
}

func (p Progress) Percent() int {
	if p.Denominator == 0 {
		return 0
	}
	return int(math.Round((float64(p.Numerator) * float64(100)) / float64(p.Denominator)))
}

type TaskState struct {
	Status  TaskStatus
	Message string
	Progress
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
