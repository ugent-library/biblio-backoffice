package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/people"
)

type ReindexPeoplePeriodicArgs struct{}

func (ReindexPeoplePeriodicArgs) Kind() string { return "reindexPeoplePeriodic" }

func (ReindexPeoplePeriodicArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

type ReindexPeoplePeriodicWorker struct {
	river.WorkerDefaults[ReindexPeoplePeriodicArgs]
	repo  *people.Repo
	index *people.Index
}

func NewReindexPeoplePeriodicWorker(repo *people.Repo, index *people.Index) *ReindexPeoplePeriodicWorker {
	return &ReindexPeoplePeriodicWorker{repo: repo, index: index}
}

func (w *ReindexPeoplePeriodicWorker) Work(ctx context.Context, job *river.Job[ReindexPeoplePeriodicArgs]) error {
	return w.index.ReindexPeople(ctx, w.repo.EachPerson)
}

type ReindexPeopleArgs struct{}

func (ReindexPeopleArgs) Kind() string { return "reindexPeople" }

type ReindexPeopleWorker struct {
	river.WorkerDefaults[ReindexPeopleArgs]
	repo  *people.Repo
	index *people.Index
}

func NewReindexPeopleWorker(repo *people.Repo, index *people.Index) *ReindexPeopleWorker {
	return &ReindexPeopleWorker{repo: repo, index: index}
}

func (w *ReindexPeopleWorker) Work(ctx context.Context, job *river.Job[ReindexPeopleArgs]) error {
	return w.index.ReindexPeople(ctx, w.repo.EachPerson)
}
