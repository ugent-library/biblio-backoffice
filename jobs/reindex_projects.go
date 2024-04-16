package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/ugent-library/biblio-backoffice/projects"
)

type ReindexProjectsPeriodicArgs struct{}

func (ReindexProjectsPeriodicArgs) Kind() string { return "reindexProjectsPeriodic" }

func (ReindexProjectsPeriodicArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

type ReindexProjectsPeriodicWorker struct {
	river.WorkerDefaults[ReindexProjectsPeriodicArgs]
	repo  *projects.Repo
	index *projects.Index
}

func NewReindexProjectsPeriodicWorker(repo *projects.Repo, index *projects.Index) *ReindexProjectsPeriodicWorker {
	return &ReindexProjectsPeriodicWorker{repo: repo, index: index}
}

func (w *ReindexProjectsPeriodicWorker) Work(ctx context.Context, job *river.Job[ReindexProjectsPeriodicArgs]) error {
	return w.index.ReindexProjects(ctx, w.repo.EachProject)
}

type ReindexProjectsArgs struct{}

func (ReindexProjectsArgs) Kind() string { return "reindexProjects" }

type ReindexProjectsWorker struct {
	river.WorkerDefaults[ReindexProjectsArgs]
	repo  *projects.Repo
	index *projects.Index
}

func NewReindexProjectsWorker(repo *projects.Repo, index *projects.Index) *ReindexProjectsWorker {
	return &ReindexProjectsWorker{repo: repo, index: index}
}

func (w *ReindexProjectsWorker) Work(ctx context.Context, job *river.Job[ReindexProjectsArgs]) error {
	return w.index.ReindexProjects(ctx, w.repo.EachProject)
}
