package projects

import (
	"context"
	"time"
)

type Project struct {
	Identifiers     []Identifier `json:"identifiers,omitempty"`
	Names           []Text       `json:"names,omitempty"`
	Descriptions    []Text       `json:"descriptions,omitempty"`
	FoundingDate    string       `json:"founding_date,omitempty"`
	DissolutionDate string       `json:"dissolution_date,omitempty"`
	Attributes      []Attribute  `json:"attributes,omitempty"`
	Deleted         bool         `json:"deleted,omitempty"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
}

func (p *Project) ID() string {
	for _, id := range p.Identifiers {
		if id.Kind == idKind {
			return id.Value
		}
	}
	return ""
}

type Identifier struct {
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
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

type ProjectIter func(context.Context, func(*Project) bool) error
