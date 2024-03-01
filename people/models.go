package people

import (
	"context"
	"fmt"
	"time"
)

type Iter[T any] func(context.Context, func(T) bool) error

type ImportOrganizationParams struct {
	Identifiers      Identifiers `json:"identifiers"`
	ParentIdentifier *Identifier `json:"parentIdentifier,omitempty"`
	Names            []Text      `json:"names,omitempty"`
	Ceased           bool        `json:"ceased"`
	CreatedAt        *time.Time  `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time  `json:"updatedAt,omitempty"`
}

type ImportPersonParams struct {
	Identifiers         Identifiers         `json:"identifiers"`
	Name                string              `json:"name"`
	PreferredName       string              `json:"preferredName,omitempty"`
	GivenName           string              `json:"givenName,omitempty"`
	FamilyName          string              `json:"familyName,omitempty"`
	PreferredGivenName  string              `json:"preferredGivenName,omitempty"`
	PreferredFamilyName string              `json:"preferredFamilyName,omitempty"`
	HonorificPrefix     string              `json:"honorificPrefix,omitempty"`
	Email               string              `json:"email,omitempty"`
	Active              bool                `json:"active"`
	Role                string              `json:"role,omitempty"`
	Username            string              `json:"username,omitempty"`
	Attributes          []Attribute         `json:"attributes,omitempty"`
	Tokens              []Token             `json:"tokens,omitempty"`
	Affiliations        []AffiliationParams `json:"affiliations,omitempty"`
	CreatedAt           *time.Time          `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time          `json:"updatedAt,omitempty"`
}

type AffiliationParams struct {
	OrganizationIdentifier Identifier `json:"organizationIdentifier"`
}

type Person struct {
	Identifiers         []Identifier `json:"identifiers"`
	Name                string       `json:"name"`
	PreferredName       string       `json:"preferredName,omitempty"`
	GivenName           string       `json:"givenName,omitempty"`
	FamilyName          string       `json:"familyName,omitempty"`
	PreferredGivenName  string       `json:"preferredGivenName,omitempty"`
	PreferredFamilyName string       `json:"preferredFamilyName,omitempty"`
	HonorificPrefix     string       `json:"honorificPrefix,omitempty"`
	Email               string       `json:"email,omitempty"`
	Active              bool         `json:"active"`
	Username            string       `json:"username,omitempty"`
	Attributes          []Attribute  `json:"attributes,omitempty"`
	CreatedAt           time.Time    `json:"createdAt"`
	UpdatedAt           time.Time    `json:"updatedAt"`
}

func (p *Person) ID() string {
	for _, id := range p.Identifiers {
		if id.Kind == idKind {
			return id.Value
		}
	}
	return ""
}

type Identifiers []Identifier

func (idents Identifiers) Has(kind string) bool {
	for _, ident := range idents {
		if ident.Kind == kind {
			return true
		}
	}
	return false
}

func (idents Identifiers) Get(kind string) string {
	for _, ident := range idents {
		if ident.Kind == kind {
			return ident.Value
		}
	}
	return ""
}

type Identifier struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func (i Identifier) String() string {
	return fmt.Sprintf("%s:%s", i.Kind, i.Value)
}

type Attribute struct {
	Scope string `json:"scope"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Text struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Token struct {
	Kind  string `json:"kind"`
	Value []byte `json:"value"`
}
