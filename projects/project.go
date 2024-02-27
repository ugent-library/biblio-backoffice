package projects

import (
	"fmt"
	"time"
)

type Project struct {
	Identifiers     []Identifier `json:"identifiers,omitempty"`
	Names           []Text       `json:"names,omitempty"`
	Descriptions    []Text       `json:"descriptions,omitempty"`
	FoundingDate    string       `json:"founding_date,omitempty"`
	DissolutionDate string       `json:"dissolution_date,omitempty"`
	Attributes      []Attribute  `json:"attributes,omitempty"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
}

type Identifier struct {
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
}

func (i *Identifier) String() string {
	return fmt.Sprintf("%s:%s", i.Kind, i.Value)
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
