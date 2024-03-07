package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/people"
)

type ReindexPeopleArgs struct{}

func (ReindexPeopleArgs) Kind() string { return "reindexPeople" }

func (ReindexPeopleArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

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

func (w *ReindexPeopleWorker) Timeout(*river.Job[ReindexPeopleArgs]) time.Duration {
	return 5 * time.Minute
}
