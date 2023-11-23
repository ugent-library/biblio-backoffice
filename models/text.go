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
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/id",
			Rule: "id.required",
		})
	}

	if t.Lang == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/lang",
			Rule: "lang.required",
		})
	} else if !slices.Contains(vocabularies.Map["language_codes"], t.Lang) {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/lang",
			Rule: "lang.invalid",
		})
	}

	if t.Text == "" {
		errs.Errors = append(errs.Errors, &okay.Error{
			Key:  "/text",
			Rule: "text.required",
		})
	}

	return errs.ErrorOrNil()
}
