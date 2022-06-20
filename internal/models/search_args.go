package models

type SearchArgs struct {
	Query    string              `query:"q,omitempty"`
	Filters  map[string][]string `query:"f,omitempty"`
	Page     int                 `query:"page"`
	Sort     []string            `query:"sort,omitempty"`
	PageSize int                 `query:"page-size"`
}

func NewSearchArgs() *SearchArgs {
	return &SearchArgs{Filters: map[string][]string{}, Page: 1, PageSize: 20}
}

func (s *SearchArgs) Clone() *SearchArgs {
	filters := make(map[string][]string, len(s.Filters))
	for field, terms := range s.Filters {
		t := make([]string, len(terms))
		copy(t, terms)
		filters[field] = t
	}

	sort := make([]string, len(s.Sort))
	copy(sort, s.Sort)

	return &SearchArgs{
		Query:    s.Query,
		Filters:  filters,
		Page:     s.Page,
		Sort:     sort,
		PageSize: s.PageSize,
	}
}

func (s *SearchArgs) WithQuery(q string) *SearchArgs {
	s.Query = q
	return s
}

func (s *SearchArgs) WithFilter(field string, terms ...string) *SearchArgs {
	if s.Filters == nil {
		s.Filters = map[string][]string{}
	}
	s.Filters[field] = terms
	return s
}

func (s *SearchArgs) WithPage(p int) *SearchArgs {
	s.Page = p
	return s
}

func (s *SearchArgs) WithSort(sort string) *SearchArgs {
	if !s.HasSort(sort) {
		s.Sort = append(s.Sort, sort)
	}
	return s
}

func (s *SearchArgs) HasFilter(field string, terms ...string) bool {
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

func (s *SearchArgs) HasSort(sort string) bool {
	for _, s := range s.Sort {
		if s == sort {
			return true
		}
	}

	return false
}

func (s *SearchArgs) FiltersFor(field string) []string {
	return s.Filters[field]
}

func (s *SearchArgs) FilterFor(field string) string {
	filters := s.Filters[field]
	if len(filters) > 0 {
		return filters[0]
	}
	return ""
}

func (s *SearchArgs) Limit() int {
	return s.PageSize
}

func (s *SearchArgs) Offset() int {
	return (s.Page - 1) * s.PageSize
}
