package models

type ExternalFields map[string][]string

func (ex ExternalFields) Add(key string, values ...string) ExternalFields {
	ex[key] = append(ex[key], values...)
	return ex
}

func (ex ExternalFields) Set(key string, values ...string) ExternalFields {
	ex[key] = values
	return ex
}

func (ex ExternalFields) Get(key string) []string {
	return ex[key]
}

func (ex ExternalFields) Delete(typ string) ExternalFields {
	delete(ex, typ)
	return ex
}

func (ex ExternalFields) Clear() ExternalFields {
	for k := range ex {
		delete(ex, k)
	}
	return ex
}
