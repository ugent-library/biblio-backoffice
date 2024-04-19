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
	seedPeopleCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.AddCommand(seedPeopleCmd)
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
