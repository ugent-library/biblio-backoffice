package projects

import (
	"fmt"
	"time"
)

type Project struct {
	Names           []Text       `json:"names,omitempty"`
	Descriptions    []Text       `json:"descriptions,omitempty"`
	FoundingDate    string       `json:"founding_date,omitempty"`
	DissolutionDate string       `json:"dissolution_date,omitempty"`
	Attributes      []Attribute  `json:"attributes,omitempty"`
	Identifiers     []Identifier `json:"identifiers,omitempty"`
}

type ProjectRecord struct {
	ID              int64        `json:"id,omitempty"`
	Names           []Text       `json:"names,omitempty"`
	Descriptions    []Text       `json:"descriptions,omitempty"`
	FoundingDate    string       `json:"founding_date,omitempty"`
	DissolutionDate string       `json:"dissolution_date,omitempty"`
	Deleted         bool         `json:"deleted,omitempty"`
	Attributes      []Attribute  `json:"attributes,omitempty"`
	Identifiers     []Identifier `json:"identifiers,omitempty"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
}

type Identifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func (i *Identifier) String() string {
	return fmt.Sprintf("%s:%s", i.Type, i.Value)
}

type Attribute struct {
	Scope string `json:"scope,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Text struct {
	Lang  string `json:"lang,omitempty"`
	Value string `json:"value,omitempty"`
}
