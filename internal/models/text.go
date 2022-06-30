package models

import "github.com/ugent-library/biblio-backend/internal/validation"

type Text struct {
	Text string `json:"text,omitempty" form:"text"`
	Lang string `json:"lang,omitempty" form:"lang"`
}

func (t Text) Validate() (errs validation.Errors) {
	if t.Lang == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/lang",
			Code:    "lang.required",
		})
	}
	if t.Text == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/text",
			Code:    "text.required",
		})
	}
	return
}
