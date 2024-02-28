package cli

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/jobs"
)

func init() {
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.AddCommand(jobsStartCmd)
}

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Biblio backoffice jobs server",
}

var jobsStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start jobs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		services := newServices()

		// start job server
		riverWorkers := river.NewWorkers()
		// river.AddWorker(riverWorkers, jobs.NewDeactivatePeopleWorker(repo))
		river.AddWorker(riverWorkers, jobs.NewReindexPeopleWorker(services.PeopleRepo, services.PeopleIndex))
		river.AddWorker(riverWorkers, jobs.NewReindexProjectsWorker(services.ProjectsRepo, services.ProjectsIndex))
		riverClient, err := river.NewClient(riverpgxv5.New(services.PgxPool), &river.Config{
			Logger:  logger,
			Workers: riverWorkers,
			Queues: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: 100},
			},
			PeriodicJobs: []*river.PeriodicJob{
				// river.NewPeriodicJob(
				// 	river.PeriodicInterval(10*time.Minute),
				// 	func() (river.JobArgs, *river.InsertOpts) {
				// 		return jobs.DeactivatePeopleArgs{}, nil
				// 	},
				// 	&river.PeriodicJobOpts{RunOnStart: true},
				// ),
				river.NewPeriodicJob(
					river.PeriodicInterval(30*time.Minute),
					func() (river.JobArgs, *river.InsertOpts) {
						return jobs.ReindexPeopleArgs{}, nil
					},
					&river.PeriodicJobOpts{RunOnStart: true},
				),
				river.NewPeriodicJob(
					river.PeriodicInterval(30*time.Minute),
					func() (river.JobArgs, *river.InsertOpts) {
						return jobs.ReindexProjectsArgs{}, nil
					},
					&river.PeriodicJobOpts{RunOnStart: true},
				),
			},
		})
		if err != nil {
			return err
		}

		logger.Info("Starting jobs server...")
		if err := riverClient.Start(context.TODO()); err != nil {
			return err
		}

		sigintOrTerm := make(chan os.Signal, 1)
		signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigintOrTerm
			logger.Info("Received SIGINT/SIGTERM; initiating soft stop (try to wait for jobs to finish)")

			softStopCtx, softStopCtxCancel := context.WithTimeout(ctx, 10*time.Second)
			defer softStopCtxCancel()

			go func() {
				select {
				case <-sigintOrTerm:
					logger.Info("Received SIGINT/SIGTERM again; initiating hard stop (cancel everything)")
					softStopCtxCancel()
				case <-softStopCtx.Done():
					logger.Info("Soft stop timeout; initiating hard stop (cancel everything)")
				}
			}()

			err := riverClient.Stop(softStopCtx)
			if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
				panic(err)
			}
			if err == nil {
				logger.Info("Soft stop succeeded")
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
				logger.Info("Hard stop timeout; ignoring stop procedure and exiting unsafely")
			} else if err != nil {
				panic(err)
			}

			// hard stop succeeded
		}()

		<-riverClient.Stopped()

		return nil
	},
}
