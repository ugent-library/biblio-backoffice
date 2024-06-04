package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/people"
	"github.com/ugent-library/biblio-backoffice/projects"
)

func init() {
	rootCmd.AddCommand(seedCmd)

	seedOrganizationsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedOrganizationsCmd)

	seedPeopleCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedPeopleCmd)

	seedProjectsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedProjectsCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed",
}

var seedOrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "seed organizations data",
	RunE: func(cmd *cobra.Command, args []string) error {
		peopleRepo := newServices().PeopleRepo

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := peopleRepo.CountOrganizations(cmd.Context())
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		iter := func(ctx context.Context, fn func(people.ImportOrganizationParams) bool) error {
			dec := json.NewDecoder(os.Stdin)
			for {
				params := people.ImportOrganizationParams{}
				if err := dec.Decode(&params); errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					return err
				}
				if !fn(params) {
					break
				}
				logger.Info("imported organization", "identifier", params.Identifiers.Get("biblio"))
			}
			return nil
		}

		if err := peopleRepo.ImportOrganizations(context.TODO(), iter); err != nil {
			return err
		}

		return newServices().PeopleIndex.ReindexOrganizations(context.TODO(), peopleRepo.EachOrganization)
	},
}

var seedPeopleCmd = &cobra.Command{
	Use:   "people",
	Short: "seed people data",
	RunE: func(cmd *cobra.Command, args []string) error {
		peopleRepo := newServices().PeopleRepo

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := peopleRepo.CountPeople(cmd.Context())
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		dec := json.NewDecoder(os.Stdin)
		for {
			params := people.ImportPersonParams{}
			if err := dec.Decode(&params); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return err
			}
			if err := peopleRepo.ImportPerson(context.TODO(), params); err != nil {
				continue
			}
			logger.Info("imported person", "username", params.Username)
		}

		return newServices().PeopleIndex.ReindexPeople(context.TODO(), peopleRepo.EachPerson)
	},
}

var seedProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "seed projects' data",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectsRepo := newServices().ProjectsRepo

		if force, _ := cmd.Flags().GetBool("force"); !force {
			count, err := projectsRepo.CountProjects(cmd.Context())
			if err != nil {
				return err
			}
			if count > 0 {
				logger.Warn("not seeding dummy data because the database is not empty")
				return nil
			}
		}

		dec := json.NewDecoder(os.Stdin)
		for {
			params := projects.ImportProjectParams{}
			if err := dec.Decode(&params); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return err
			}
			if err := projectsRepo.ImportProject(context.TODO(), params); err != nil {
				continue
			}
			logger.Info("imported project", "iwetoID", params.Identifiers.Get("iweto"))
		}

		return newServices().ProjectsIndex.ReindexProjects(context.TODO(), projectsRepo.EachProject)
	},
}
