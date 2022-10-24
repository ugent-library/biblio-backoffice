package pagination

import (
	"math"
)

type Pagination struct {
	Offset, Limit, Total int
}

func (p Pagination) MaxPages() int {
	return 500
}

func (p Pagination) MaxVisiblePages() int {
	return 5
}

func (p Pagination) TotalPages() int {
	if p.Limit == 0 {
		return 1
	}
	total := int(math.Ceil(float64(p.Total) / float64(p.Limit)))
	if total > p.MaxPages() {
		return p.MaxPages()
	}
	return total
}

func (p Pagination) Page() int {
	if p.Limit == 0 {
		return 1
	}
	return p.Offset/p.Limit + 1
}

func (p Pagination) HasPreviousPage() bool {
	return p.Page() > 1
}

func (p Pagination) HasNextPage() bool {
	return p.Page() < p.TotalPages()
}

func (p Pagination) PreviousPage() int {
	return p.Page() - 1
}

func (p Pagination) NextPage() int {
	if p.HasNextPage() {
		return p.Page() + 1
	}
	return 0
}

func (p Pagination) FirstOnPage() int {
	if p.Limit == 0 || p.Total == 0 {
		return 0
	}
	return p.Offset + 1
}

func (p Pagination) LastOnPage() int {
	if p.HasNextPage() {
		return p.Offset + p.Limit
	}
	return p.Total
}

// ported from https://metacpan.org/dist/Data-SpreadPagination/source/lib/Data/SpreadPagination.pm
func (p Pagination) PagesWithEllipsis() []int {
	if p.Total <= p.Limit {
		return []int{1}
	}

	pages := []int{}

	ranges := p.pageRanges()

	if ranges[0] == nil && p.Page() > 1 {
		pages = append(pages, 0)
	} else if ranges[0] != nil {
		pages = append(pages, makeRange(ranges[0][0], ranges[0][1])...)
		if ranges[1] != nil && ranges[1][0]-ranges[0][1] > 1 {
			pages = append(pages, 0)
		}
	}

	if ranges[1] != nil {
		pages = append(pages, makeRange(ranges[1][0], ranges[1][1])...)
	}

	pages = append(pages, p.Page())

	if ranges[2] != nil {
		pages = append(pages, makeRange(ranges[2][0], ranges[2][1])...)
	}

	if ranges[3] == nil && p.Page() < p.TotalPages() {
		pages = append(pages, 0)
	} else if ranges[3] != nil {
		if ranges[2] != nil && ranges[3][0]-ranges[2][1] > 1 {
			pages = append(pages, 0)
		}
		pages = append(pages, makeRange(ranges[3][0], ranges[3][1])...)
	}

	return pages
}

// TODO enforce Page <= MaxPages
func (p Pagination) pageRanges() [][]int {
	page := p.Page()
	totalPages := p.TotalPages()

	visiblePages := 0
	if p.MaxVisiblePages() < (totalPages - 1) {
		visiblePages = p.MaxVisiblePages()
	} else {
		visiblePages = totalPages - 1
	}

	var qSize []int

	if totalPages-1 <= p.MaxVisiblePages() {
		qSize = []int{page - 1, 0, 0, totalPages - page}
	} else {
		qSize = []int{
			int(math.Floor(float64(visiblePages) / 4)),
			int(math.Round(float64(visiblePages) / 4)),
			int(math.Ceil(float64(visiblePages) / 4)),
			int(math.Round((float64(visiblePages) - math.Round(float64(visiblePages)/4)) / 3)),
		}
		if page-qSize[0] < 1 {
			addPages := qSize[0] + qSize[1] - page + 1
			qSize = []int{
				page - 1,
				0,
				qSize[2] + int(math.Ceil(float64(addPages)/2)),
				qSize[3] + int(math.Floor(float64(addPages)/2)),
			}
		} else if page-qSize[1]-int(math.Ceil(float64(qSize[1])/3)) <= qSize[0] {
			adj := int(math.Ceil(float64(3*(page-qSize[0]-1)) / 4))
			addPages := qSize[1] - adj
			qSize = []int{
				qSize[0],
				adj,
				qSize[2] + int(math.Ceil(float64(addPages)/2)),
				qSize[3] + int(math.Floor(float64(addPages)/2)),
			}
		} else if page+qSize[3] >= totalPages {
			addPages := qSize[2] + qSize[3] - totalPages + page
			qSize = []int{
				qSize[0] + int(math.Floor(float64(addPages)/2)),
				qSize[1] + int(math.Ceil(float64(addPages)/2)),
				0,
				totalPages - page,
			}
		} else if page+qSize[2] >= totalPages-qSize[3] {
			adj := int(math.Ceil(float64(3*(totalPages-page-qSize[3])) / 4))
			addPages := qSize[2] - adj
			qSize = []int{
				qSize[0] + int(math.Floor(float64(addPages)/2)),
				qSize[1] + int(math.Ceil(float64(addPages)/2)),
				adj,
				qSize[3],
			}
		}
	}

	pageRanges := make([][]int, 4)

	if qSize[0] != 0 {
		pageRanges[0] = []int{1, qSize[0]}
	}
	if qSize[1] != 0 {
		pageRanges[1] = []int{page - qSize[1], page - 1}
	}
	if qSize[2] != 0 {
		pageRanges[2] = []int{page + 1, page + qSize[2]}
	}
	if qSize[3] != 0 {
		pageRanges[3] = []int{totalPages - qSize[3] + 1, totalPages}
	}

	return pageRanges
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
