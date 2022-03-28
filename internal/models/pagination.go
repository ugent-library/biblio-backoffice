package models

import "math"

type Pagination struct {
	Offset, Limit, Total int
}

func (p Pagination) Page() int {
	if p.Limit == 0 {
		return 1
	}
	return p.Offset/p.Limit + 1
}

func (p Pagination) LastPage() int {
	if p.Limit == 0 {
		return 1
	}
	return int(math.Ceil(float64(p.Total) / float64(p.Limit)))
}

func (p Pagination) PreviousPage() bool {
	return p.Page() > 1
}

func (p Pagination) NextPage() bool {
	return p.Page() < p.LastPage()
}

func (p Pagination) FirstOnPage() int {
	if p.Limit == 0 || p.Total == 0 {
		return 0
	}
	return p.Offset + 1
}

func (p Pagination) LastOnPage() int {
	if p.NextPage() {
		return p.Offset + p.Limit
	}
	return p.Total
}
