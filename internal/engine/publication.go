package engine

type Publication map[string]interface{}

type PublicationHits struct {
	Total int           `json:"total,omitempty"`
	Hits  []Publication `json:"hits,omitempty"`
}

func (r Publication) ID() string {
	return r["_id"].(string)
}
