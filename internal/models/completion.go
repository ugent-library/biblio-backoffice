package models

type Completion struct {
	ID           string `json:"id"`
	Heading      string `json:"heading"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url"`
}
