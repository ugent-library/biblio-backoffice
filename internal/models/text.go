package models

type Text struct {
	Text string `json:"text,omitempty" form:"text"`
	Lang string `json:"lang,omitempty" form:"lang"`
}
