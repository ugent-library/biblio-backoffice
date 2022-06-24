package publicationviewing

import (
	"github.com/ugent-library/biblio-backend/internal/render/fields"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/helpers"
)

func detailFieldsBook(ctx Context) []*fields.Fields {
	p := ctx.Publication
	return []*fields.Fields{
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label: ctx.T("builder.type"),
					Value: ctx.TS("publication_types", p.Type),
				},
				&fields.Text{
					Label:         ctx.T("builder.doi"),
					Value:         ctx.Publication.DOI,
					ValueTemplate: "format/doi",
				},
				&fields.Text{
					Label: ctx.T("builder.classification"),
					Value: ctx.TS("publication_classifications", p.Classification),
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:    ctx.T("builder.title"),
					Value:    p.Title,
					Required: true,
				},
				&fields.Text{
					Label:  ctx.T("builder.alternative_title"),
					List:   true,
					Values: p.AlternativeTitle,
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label:  ctx.T("builder.language"),
					List:   true,
					Values: p.Language,
				},
				&fields.Text{
					Label: ctx.T("builder.publication_status"),
					Value: ctx.TS("publication_publishing_statuses", p.PublicationStatus),
				},
				&fields.Text{
					Label: ctx.T("builder.extern"),
					Value: helpers.FormatBool(p.Extern, "✓", "-"),
				},
				&fields.Text{
					Label:    ctx.T("builder.year"),
					Value:    p.Year,
					Required: true,
				},
				&fields.Text{
					Label: ctx.T("builder.place_of_publication"),
					Value: p.PlaceOfPublication,
				},
				&fields.Text{
					Label: ctx.T("builder.publisher"),
					Value: p.Publisher,
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label: ctx.T("builder.page_count"),
					Value: p.PageCount,
				},
			},
		},
		{
			Theme: "default",
			Fields: []fields.Field{
				&fields.Text{
					Label: ctx.T("builder.wos_type"),
					Value: p.WOSType,
				},
				&fields.Text{
					Label: ctx.T("builder.wos_id"),
					Value: p.WOSID,
				},
				&fields.Text{
					Label:  ctx.T("builder.issn"),
					List:   true,
					Values: p.ISSN,
				},
				&fields.Text{
					Label:  ctx.T("builder.eissn"),
					List:   true,
					Values: p.EISSN,
				},
				&fields.Text{
					Label:  ctx.T("builder.isbn"),
					List:   true,
					Values: p.ISBN,
				},
				&fields.Text{
					Label:  ctx.T("builder.eisbn"),
					List:   true,
					Values: p.EISBN,
				},
			},
		},
	}
}
