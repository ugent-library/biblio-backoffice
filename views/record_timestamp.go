package views

import (
	"strings"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
)

func CreatedBy(c *ctx.Ctx, createdAt *time.Time, createdBy *models.Person) string {
	sb := &strings.Builder{}

	sb.WriteString("Created ")
	sb.WriteString(createdAt.In(c.Timezone).Format("2006-01-02 15:04"))

	if createdBy != nil {
		sb.WriteString(" by ")
		addEditorInfo(c, sb, createdBy)
	}

	sb.WriteString(".")

	return sb.String()
}

func UpdatedBy(c *ctx.Ctx, updatedAt *time.Time, updatedBy *models.Person, lastUpdatedBy *models.Person) string {
	sb := &strings.Builder{}

	if updatedBy != nil {
		sb.WriteString("Edited ")
		sb.WriteString(updatedAt.In(c.Timezone).Format("2006-01-02 15:04"))
		sb.WriteString(" by ")
		addEditorInfo(c, sb, updatedBy)
		sb.WriteString(".")
	} else {
		sb.WriteString("System update ")
		sb.WriteString(updatedAt.In(c.Timezone).Format("2006-01-02 15:04"))
		sb.WriteString(".")

		if lastUpdatedBy != nil {
			sb.WriteString(" Last edit by ")
			addEditorInfo(c, sb, lastUpdatedBy)
			sb.WriteString(".")
		}
	}

	return sb.String()
}

func addEditorInfo(c *ctx.Ctx, sb *strings.Builder, person *models.Person) {
	if person.CanCurate() && !c.User.CanCurate() {
		sb.WriteString("a Biblio team member")
	} else {
		sb.WriteString(person.FullName)
	}
}
