package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/people"
)

type ReindexOrganizationsArgs struct{}

func (ReindexOrganizationsArgs) Kind() string { return "reindexOrganizations" }

func (ReindexOrganizationsArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

type ReindexOrganizationsWorker struct {
	river.WorkerDefaults[ReindexOrganizationsArgs]
	repo  *people.Repo
	index *people.Index
}

func NewReindexOrganizationsWorker(repo *people.Repo, index *people.Index) *ReindexOrganizationsWorker {
	return &ReindexOrganizationsWorker{repo: repo, index: index}
}

func (w *ReindexOrganizationsWorker) Work(ctx context.Context, job *river.Job[ReindexOrganizationsArgs]) error {
	return w.index.ReindexOrganizations(ctx, w.repo.EachOrganization)
}
