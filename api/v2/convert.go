package api

import "github.com/ugent-library/biblio-backoffice/people"

func convertImportOrganizationParams(from ImportOrganizationParams) people.ImportOrganizationParams {
	identifiers := make([]people.Identifier, len(from.Identifiers))
	for i, ident := range from.Identifiers {
		identifiers[i] = people.Identifier(ident)
	}
	names := make([]people.Text, len(from.Names))
	for i, text := range from.Names {
		names[i] = people.Text(text)
	}
	to := people.ImportOrganizationParams{
		Identifiers: identifiers,
		Names:       names,
		Ceased:      from.Ceased.Value,
	}
	if from.ParentIdentifier.Set {
		ident := people.Identifier(from.ParentIdentifier.Value)
		to.ParentIdentifier = &ident
	}
	if from.CreatedAt.Set {
		to.CreatedAt = &from.CreatedAt.Value
	}
	if from.UpdatedAt.Set {
		to.UpdatedAt = &from.UpdatedAt.Value
	}
	return to
}
