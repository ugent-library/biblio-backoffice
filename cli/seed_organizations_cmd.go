package cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/people"
)

func init() {
	seedOrganizationsCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedOrganizationsCmd)
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
