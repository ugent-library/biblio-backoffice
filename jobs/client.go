package jobs

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
)

type JobsConfig struct {
	PgxPool       *pgxpool.Pool
	PeopleRepo    *people.Repo
	PeopleIndex   *people.Index
	ProjectsRepo  *projects.Repo
	ProjectsIndex *projects.Index
	Logger        *slog.Logger
}

func Start(ctx context.Context, c JobsConfig) error {
	// start job server
	riverWorkers := river.NewWorkers()
	// river.AddWorker(riverWorkers, NewDeactivatePeopleWorker(repo))
	river.AddWorker(riverWorkers, NewReindexPeopleWorker(c.PeopleRepo, c.PeopleIndex))
	river.AddWorker(riverWorkers, NewReindexProjectsWorker(c.ProjectsRepo, c.ProjectsIndex))
	riverClient, err := river.NewClient(riverpgxv5.New(c.PgxPool), &river.Config{
		Logger:  c.Logger,
		Workers: riverWorkers,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		PeriodicJobs: []*river.PeriodicJob{
			// river.NewPeriodicJob(
			// 	river.PeriodicInterval(10*time.Minute),
			// 	func() (river.JobArgs, *river.InsertOpts) {
			// 		return DeactivatePeopleArgs{}, nil
			// 	},
			// 	&river.PeriodicJobOpts{RunOnStart: true},
			// ),
			river.NewPeriodicJob(
				river.PeriodicInterval(30*time.Minute),
				func() (river.JobArgs, *river.InsertOpts) {
					return ReindexPeopleArgs{}, nil
				},
				&river.PeriodicJobOpts{RunOnStart: true},
			),
			river.NewPeriodicJob(
				river.PeriodicInterval(30*time.Minute),
				func() (river.JobArgs, *river.InsertOpts) {
					return ReindexProjectsArgs{}, nil
				},
				&river.PeriodicJobOpts{RunOnStart: true},
			),
		},
	})
	if err != nil {
		return err
	}

	c.Logger.Info("Starting jobs server...")
	if err := riverClient.Start(context.TODO()); err != nil {
		return err
	}

	sigintOrTerm := make(chan os.Signal, 1)
	signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigintOrTerm
		c.Logger.Info("Received SIGINT/SIGTERM; initiating soft stop (try to wait for jobs to finish)")

		softStopCtx, softStopCtxCancel := context.WithTimeout(ctx, 10*time.Second)
		defer softStopCtxCancel()

		go func() {
			select {
			case <-sigintOrTerm:
				c.Logger.Info("Received SIGINT/SIGTERM again; initiating hard stop (cancel everything)")
				softStopCtxCancel()
			case <-softStopCtx.Done():
				c.Logger.Info("Soft stop timeout; initiating hard stop (cancel everything)")
			}
		}()

		err := riverClient.Stop(softStopCtx)
		if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			panic(err)
		}
		if err == nil {
			c.Logger.Info("Soft stop succeeded")
			return
		}

		hardStopCtx, hardStopCtxCancel := context.WithTimeout(ctx, 10*time.Second)
		defer hardStopCtxCancel()

		// As long as all jobs respect context cancellation, StopAndCancel will
		// always work. However, in the case of a bug where a job blocks despite
		// being cancelled, it may be necessary to either ignore River's stop
		// result (what's shown here) or have a supervisor kill the process.
		err = riverClient.StopAndCancel(hardStopCtx)
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			c.Logger.Info("Hard stop timeout; ignoring stop procedure and exiting unsafely")
		} else if err != nil {
			panic(err)
		}

		// hard stop succeeded
	}()

	return nil
}
