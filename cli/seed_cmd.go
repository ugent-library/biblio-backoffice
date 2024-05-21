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
				zapLogger.Warnf("Not seeding dummy data because the database is not empty")
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
					zapLogger.Errorf("unable to decode json: %w", err)
					return err
				}
				if !fn(params) {
					break
				}
				zapLogger.Infof("imported organization %s", params.Identifiers.Get("biblio"))
			}
			return nil
		}

		if err := peopleRepo.ImportOrganizations(context.TODO(), iter); err != nil {
			return err
		}

		return nil
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
				zapLogger.Warnf("Not seeding dummy data because the database is not empty")
				return nil
			}
		}

		dec := json.NewDecoder(os.Stdin)
		for {
			params := people.ImportPersonParams{}
			if err := dec.Decode(&params); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				zapLogger.Errorf("unable to decode json: %w", err)
				return err
			}
			if err := peopleRepo.ImportPerson(context.TODO(), params); err != nil {
				zapLogger.Errorf("unable to import person %s: %w", params.Username, err)
				continue
			}
			zapLogger.Infof("imported person %s", params.Username)
		}
		return nil
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
				zapLogger.Warnf("Not seeding dummy data because the database is not empty")
				return nil
			}
		}

		dec := json.NewDecoder(os.Stdin)
		for {
			params := projects.ImportProjectParams{}
			if err := dec.Decode(&params); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				zapLogger.Errorf("unable to decode json: %w", err)
				return err
			}
			if err := projectsRepo.ImportProject(context.TODO(), params); err != nil {
				zapLogger.Errorf("unable to import project %s: %w", params.Identifiers.Get("iweto"), err)
				continue
			}
			zapLogger.Infof("imported project %s", params.Identifiers.Get("iweto"))
		}
		return nil
	},
}
