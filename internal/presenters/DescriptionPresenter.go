package presenters

import (
	"github.com/ugent-library/biblio-backend/internal/DAO"
	"github.com/ugent-library/biblio-backend/internal/fields"
)

type DescriptionPresenter struct {
	Description *DAO.DescriptionDAO
}

// Processes the Description, encapsulates all business logic
// Results in a tree of elements which can be rendered

func (dp *DescriptionPresenter) Process() map[string]fields.FieldSet {
	view := make(map[string]fields.FieldSet)

	items := []fields.Field{}

	items = append(items, &fields.TextField{
		Label: "Title",
		Value: dp.Description.Title,
	})

	items = append(items, &fields.TextField{
		Label: "DOI",
		Value: dp.Description.DOI,
	})

	items = append(items, &fields.TextField{
		Label: "Publication Type",
		Value: dp.Description.PublicationType,
	})

	view["publication_details"] = fields.FieldSet{
		Label: "Publication Details",
		Items: items,
	}

	if dp.Description.PublicationType == "journal_article" {
		items := []fields.Field{}

		items = append(items, &fields.TextField{
			Label: "Conference",
			Value: dp.Description.Conference,
		})

		items = append(items, &fields.TextField{
			Label: "Conference location",
			Value: dp.Description.ConferenceLocation,
		})

		view["conference_details"] = fields.FieldSet{
			Label: "Conference Details",
			Items: items,
		}
	}

	return view
}
