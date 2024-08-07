package views

import (
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/models"
)

func CreatedBy(c *ctx.Ctx, createdAt *time.Time, createdBy *models.Person) string {
	s := fmt.Sprintf("Created %s", createdAt.In(c.Timezone).Format("2006-01-02 15:04"))
	if createdBy != nil {
		s += " by " + createdBy.FullName
	}
	s += "."
	return s
}

func UpdatedBy(c *ctx.Ctx, updatedAt *time.Time, updatedBy *models.Person, lastUpdatedBy *models.Person) string {
	if updatedBy != nil {
		return fmt.Sprintf("Edited %s by %s.", updatedAt.In(c.Timezone).Format("2006-01-02 15:04"), updatedBy.FullName)
	}

	s := fmt.Sprintf("System update %s.", updatedAt.In(c.Timezone).Format("2006-01-02 15:04"))
	if lastUpdatedBy != nil {
		s += " Last edit by " + lastUpdatedBy.FullName + "."
	}
	return s
}
