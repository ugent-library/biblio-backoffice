package engine

type Filters map[string][]string

type SearchArgs struct {
	Query   string  `form:"q,omitempty"`
	Filters Filters `form:"f,omitempty"`
	Page    int     `form:"page"`
}

func NewSearchArgs() *SearchArgs {
	return &SearchArgs{Filters: Filters{}, Page: 1}
}

func (s *SearchArgs) WithQuery(q string) *SearchArgs {
	s.Query = q
	return s
}

func (s *SearchArgs) WithFilter(field string, terms ...string) *SearchArgs {
	if s.Filters == nil {
		s.Filters = Filters{}
	}
	s.Filters[field] = terms
	return s
}

func (s *SearchArgs) WithPage(p int) *SearchArgs {
	s.Page = p
	return s
}

func (s *SearchArgs) HasFilter(field string, terms ...string) bool {
	if s.Filters == nil {
		return false
	}
	filter, ok := s.Filters[field]
	if !ok {
		return false
	}

	for _, term := range terms {
		var contains bool
		for _, t := range filter {
			if t == term {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	return true
}
