package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
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
	return nil
}
