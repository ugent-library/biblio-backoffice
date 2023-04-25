package models

type Identifiers map[string][]string

func (i Identifiers) Add(key, val string) {
	i[key] = append(i[key], val)
}

func (i Identifiers) Set(key, val string) {
	i[key] = []string{val}
}

func (i Identifiers) Get(key string) string {
	vals := i[key]
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (i Identifiers) Has(key, val string) bool {
	for _, v := range i[key] {
		if v == val {
			return true
		}
	}
	return false
}
