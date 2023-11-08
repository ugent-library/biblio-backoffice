package models

type Values map[string][]string

func (v Values) Add(key, val string) {
	v[key] = append(v[key], val)
}

func (v Values) Set(key, val string) {
	v[key] = []string{val}
}

func (v Values) SetAll(key string, vals ...string) {
	v[key] = vals
}

func (v Values) Get(key string) string {
	vals := v[key]
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (v Values) GetAll(key string) []string {
	return v[key]
}

func (v Values) Has(key, val string) bool {
	for _, v := range v[key] {
		if v == val {
			return true
		}
	}
	return false
}
