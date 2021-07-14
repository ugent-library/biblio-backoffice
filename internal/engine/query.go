package engine

type Filters map[string][]string

type Query struct {
	QueryString string  `form:"q,omitempty"`
	Filters     Filters `form:"f,omitempty"`
	Page        int     `form:"page"`
}

func NewQuery() *Query {
	return &Query{Filters: Filters{}, Page: 1}
}

func (q *Query) WithQueryString(qs string) *Query {
	q.QueryString = qs
	return q
}

func (q *Query) WithFilter(field string, terms ...string) *Query {
	if q.Filters == nil {
		q.Filters = Filters{}
	}
	q.Filters[field] = terms
	return q
}

func (q *Query) WithPage(p int) *Query {
	q.Page = p
	return q
}
