package engine

type Query struct {
	QueryString string `form:"q,omitempty"`
	Page        int    `form:"page"`
}

func (q Query) WithPage(p int) Query {
	q.Page = p
	return q
}
