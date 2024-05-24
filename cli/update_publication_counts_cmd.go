package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
)

func init() {
	rootCmd.AddCommand(updatePublicationCounts)
}

var updatePublicationCounts = &cobra.Command{
	Use:   "update-publication-counts",
	Short: "Update publication counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()

		services := newServices()

		peopleCounts := make(map[string]int)
		projectCounts := make(map[string]int)

		services.Repo.EachPublication(func(p *models.Publication) bool {
			if p.Status != "public" {
				return true
			}
			for _, rp := range p.RelatedProjects {
				projectCounts[rp.ProjectID] += 1
			}
			for _, a := range p.Author {
				if a.PersonID != "" {
					peopleCounts[a.PersonID] += 1
				}
			}
			if p.Type == "book_editor" || p.Type == "issue_editor" {
				for _, a := range p.Editor {
					if a.PersonID != "" {
						peopleCounts[a.PersonID] += 1
					}
				}
			}
			return true
		})

		var iterErr error

		services.PeopleRepo.EachPerson(ctx, func(p *people.Person) bool {
			var count int
			for _, id := range p.Identifiers.GetAll("id") {
				if c, ok := peopleCounts[id]; ok {
					count = c
					break
				}
			}

			id := p.Identifiers.Get("id")
			iterErr = services.PeopleRepo.SetPersonPublicationCount(ctx, "id", id, count)

			return iterErr == nil
		})

		if iterErr != nil {
			return iterErr
		}

		services.ProjectsRepo.EachProject(ctx, func(p *projects.Project) bool {
			var count int
			for _, id := range p.Identifiers.GetAll("iweto") {
				if c, ok := projectCounts[id]; ok {
					count = c
					break
				}
			}

			id := p.Identifiers.Get("id")
			iterErr = services.ProjectsRepo.SetProjectPublicationCount(ctx, "id", id, count)

			return iterErr == nil
		})

		return iterErr
	},
}
