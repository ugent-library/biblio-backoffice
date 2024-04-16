package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/people"
)

type ReindexOrganizationsPeriodicArgs struct{}

func (ReindexOrganizationsPeriodicArgs) Kind() string { return "reindexOrganizationsPeriodic" }

func (ReindexOrganizationsPeriodicArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

type ReindexOrganizationsPeriodicWorker struct {
	river.WorkerDefaults[ReindexOrganizationsPeriodicArgs]
	repo  *people.Repo
	index *people.Index
}

func NewReindexOrganizationsPeriodicWorker(repo *people.Repo, index *people.Index) *ReindexOrganizationsPeriodicWorker {
	return &ReindexOrganizationsPeriodicWorker{repo: repo, index: index}
}

func (w *ReindexOrganizationsPeriodicWorker) Work(ctx context.Context, job *river.Job[ReindexOrganizationsPeriodicArgs]) error {
	return w.index.ReindexOrganizations(ctx, w.repo.EachOrganization)
}

type ReindexOrganizationsArgs struct{}

func (ReindexOrganizationsArgs) Kind() string { return "reindexOrganizations" }

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
