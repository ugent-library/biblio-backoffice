package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

type UpdatePublicationCountArgs struct{}

func (UpdatePublicationCountArgs) Kind() string { return "updatePublicationCount" }

func (UpdatePublicationCountArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}

type UpdatePublicationCountWorker struct {
	river.WorkerDefaults[UpdatePublicationCountArgs]
	repo         *repositories.Repo
	peopleRepo   *people.Repo
	projectsRepo *projects.Repo
}

func NewUpdatePublicationCountWorker(repo *repositories.Repo, peopleRepo *people.Repo, projectsRepo *projects.Repo) *UpdatePublicationCountWorker {
	return &UpdatePublicationCountWorker{
		repo:         repo,
		peopleRepo:   peopleRepo,
		projectsRepo: projectsRepo,
	}
}

func (w *UpdatePublicationCountWorker) Work(ctx context.Context, job *river.Job[UpdatePublicationCountArgs]) error {
	peopleCounts := make(map[string]int)
	projectCounts := make(map[string]int)

	w.repo.EachPublication(func(p *models.Publication) bool {
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

	w.peopleRepo.EachPerson(ctx, func(p *people.Person) bool {
		var count int
		for _, id := range p.Identifiers.GetAll("id") {
			if c, ok := peopleCounts[id]; ok {
				count = c
				break
			}
		}

		id := p.Identifiers.Get("id")
		iterErr = w.peopleRepo.SetPersonPublicationCount(ctx, "id", id, count)

		return iterErr == nil
	})

	if iterErr != nil {
		return iterErr
	}

	w.projectsRepo.EachProject(ctx, func(p *projects.Project) bool {
		var count int
		for _, id := range p.Identifiers.GetAll("iweto") {
			if c, ok := projectCounts[id]; ok {
				count = c
				break
			}
		}

		id := p.Identifiers.Get("id")
		iterErr = w.projectsRepo.SetProjectPublicationCount(ctx, "id", id, count)

		return iterErr == nil
	})

	return iterErr
}

func (w *UpdatePublicationCountWorker) Timeout(*river.Job[UpdatePublicationCountArgs]) time.Duration {
	return 10 * time.Minute
}
