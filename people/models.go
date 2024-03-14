package people

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")

type InvalidIdentifierError struct {
	Identifier string
}

func (e *InvalidIdentifierError) Error() string {
	return fmt.Sprintf("%q is not a valid identifier", e.Identifier)
}

type DuplicateError struct {
	Identifier string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("identifier %s already exists", e.Identifier)
}

type Iter[T any] func(context.Context, func(T) bool) error

type ImportOrganizationParams struct {
	Identifiers      Identifiers `json:"identifiers"`
	ParentIdentifier *Identifier `json:"parentIdentifier,omitempty"`
	Names            Texts       `json:"names,omitempty"`
	Ceased           bool        `json:"ceased"`
	CeasedOn         *time.Time  `json:"ceasedOn,omitempty"`
	CreatedAt        *time.Time  `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time  `json:"updatedAt,omitempty"`
}

type Organization struct {
	Identifiers Identifiers          `json:"identifiers"`
	Names       Texts                `json:"names,omitempty"`
	Ceased      bool                 `json:"ceased"`
	CeasedOn    *time.Time           `json:"ceasedOn,omitempty"`
	Position    int                  `json:"position"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Parents     []ParentOrganization `json:"parents"` // TODO just use *Organization?
}

type ParentOrganization struct {
	Identifiers Identifiers `json:"identifiers"`
	Names       Texts       `json:"names,omitempty"`
	Ceased      bool        `json:"ceased"`
	CeasedOn    *time.Time  `json:"ceasedOn,omitempty"`
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
	Attributes          Attributes          `json:"attributes,omitempty"`
	Tokens              Tokens              `json:"tokens,omitempty"`
	Affiliations        []AffiliationParams `json:"affiliations,omitempty"`
	CreatedAt           *time.Time          `json:"createdAt,omitempty"` // TODO pointers not needed, use IsZero
	UpdatedAt           *time.Time          `json:"updatedAt,omitempty"`
}

type AffiliationParams struct {
	OrganizationIdentifier Identifier `json:"organizationIdentifier"`
}

type Person struct {
	Identifiers         Identifiers   `json:"identifiers"`
	Name                string        `json:"name"`
	PreferredName       string        `json:"preferredName,omitempty"`
	GivenName           string        `json:"givenName,omitempty"`
	FamilyName          string        `json:"familyName,omitempty"`
	PreferredGivenName  string        `json:"preferredGivenName,omitempty"`
	PreferredFamilyName string        `json:"preferredFamilyName,omitempty"`
	HonorificPrefix     string        `json:"honorificPrefix,omitempty"`
	Email               string        `json:"email,omitempty"`
	Active              bool          `json:"active"`
	Role                string        `json:"role,omitempty"`
	Username            string        `json:"username,omitempty"`
	Attributes          Attributes    `json:"attributes,omitempty"`
	Tokens              Tokens        `json:"tokens,omitempty"`
	Affiliations        []Affiliation `json:"affiliations,omitempty"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
}

type Affiliation struct {
	Organization *Organization `json:"organization"`
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

func (idents Identifiers) GetAll(kind string) (ids []string) {
	for _, ident := range idents {
		if ident.Kind == kind {
			ids = append(ids, ident.Value)
		}
	}
	return
}

type Identifier struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func NewIdentifier(str string) (Identifier, error) {
	k, v, ok := strings.Cut(str, ":")
	if !ok || k == "" || v == "" {
		return Identifier{}, &InvalidIdentifierError{Identifier: str}
	}
	return Identifier{Kind: k, Value: v}, nil
}

func (i Identifier) String() string {
	return fmt.Sprintf("%s:%s", i.Kind, i.Value)
}

type Attributes []Attribute

func (a Attributes) Get(scope, key string) string {
	for _, attr := range a {
		if attr.Scope == scope && attr.Key == key {
			return attr.Value
		}
	}
	return ""
}

type Attribute struct {
	Scope string `json:"scope"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Texts []Text

func (t Texts) Get(lang string) string {
	for _, text := range t {
		if text.Lang == lang {
			return text.Value
		}
	}
	return ""
}

type Text struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Tokens []Token

func (t Tokens) Get(kind string) []byte {
	for _, token := range t {
		if token.Kind == kind {
			return token.Value
		}
	}
	return nil
}

type Token struct {
	Kind  string `json:"kind"`
	Value []byte `json:"value"`
}

type SearchParams struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type SearchResults[T any] struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Total  int    `json:"total"`
	Hits   []T    `json:"hits"`
}
