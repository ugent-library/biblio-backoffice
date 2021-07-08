package engine

type Publication map[string]interface{}

type PublicationHits struct {
	Total        int           `json:"total"`
	Start        int           `json:"start"`
	Limit        int           `json:"limit"`
	PageSize     int           `json:"page_size"`
	Page         int           `json:"page"`
	LastPage     int           `json:"last_page"`
	PreviousPage bool          `json:"previous_page"`
	NextPage     bool          `json:"next_page"`
	Hits         []Publication `json:"hits"`
}

func (r Publication) ID() string {
	return r["_id"].(string)
}
