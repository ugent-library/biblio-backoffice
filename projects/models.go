package projects

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

type ProjectIter func(context.Context, func(*Project) bool) error

type ImportProjectParams struct {
	Identifiers  Identifiers
	Names        Texts
	Descriptions Texts
	StartDate    string
	EndDate      string
	Deleted      bool
	Attributes   Attributes
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type AddProjectParams struct {
	Identifiers  Identifiers
	Names        Texts
	Descriptions Texts
	StartDate    string
	EndDate      string
	Deleted      bool
	Attributes   Attributes
}

type Project struct {
	Identifiers  Identifiers `json:"identifiers,omitempty"`
	Names        Texts       `json:"names,omitempty"`
	Descriptions Texts       `json:"descriptions,omitempty"`
	StartDate    string      `json:"startDate,omitempty"`
	EndDate      string      `json:"endDate,omitempty"`
	Attributes   Attributes  `json:"attributes,omitempty"`
	Deleted      bool        `json:"deleted,omitempty"`
	CreatedAt    time.Time   `json:"created_at,omitempty"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
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
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
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
	Scope string `json:"scope,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
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
	Lang  string `json:"lang,omitempty"`
	Value string `json:"value,omitempty"`
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
