package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/people"
)

type DeactivatePeopleArgs struct{}

func (DeactivatePeopleArgs) Kind() string { return "deactivatePeople" }

func (DeactivatePeopleArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 10 * time.Minute,
		},
	}
}

type DeactivatePeopleWorker struct {
	river.WorkerDefaults[DeactivatePeopleArgs]
	repo *people.Repo
}

func NewDeactivatePeopleWorker(repo *people.Repo) *DeactivatePeopleWorker {
	return &DeactivatePeopleWorker{repo: repo}
}

func (w *DeactivatePeopleWorker) Work(ctx context.Context, job *river.Job[DeactivatePeopleArgs]) error {
	return w.repo.DeactivatePeople(ctx)
}
