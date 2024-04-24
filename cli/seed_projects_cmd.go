package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/projects"
)

func init() {
	seedProjectsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedProjectsCmd)
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
