package models

import (
	"slices"

	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
)

type Text struct {
	Text string `json:"text,omitempty"`
	Lang string `json:"lang,omitempty"`
	ID   string `json:"id,omitempty"`
}

func (t Text) Validate() error {
	errs := okay.NewErrors()

	if t.ID == "" {
		errs.Add(okay.NewError("/id", "id.required"))
	}

	if t.Lang == "" {
		errs.Add(okay.NewError("/lang", "lang.required"))
	} else if !slices.Contains(vocabularies.Map["language_codes"], t.Lang) {
		errs.Add(okay.NewError("/lang", "lang.invalid"))
	}

	if t.Text == "" {
		errs.Add(okay.NewError("/text", "text.required"))
	}

	return errs.ErrorOrNil()
}
