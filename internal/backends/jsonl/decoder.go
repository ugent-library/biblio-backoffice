package jsonl

import (
	"encoding/json"
	"io"

	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
)

type Decoder struct {
	jsonDecoder *json.Decoder
}

func NewDecoder(r io.Reader) backends.PublicationDecoder {
	return &Decoder{json.NewDecoder(r)}
}

func (d *Decoder) Decode(p *models.Publication) error {
	return d.jsonDecoder.Decode(p)
}
