package people

import (
	"fmt"
	"time"
)

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
	// Active              bool         `json:"active"`
	// Username    string       `json:"username,omitempty"`
	Attributes []Attribute `json:"attributes"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
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
