package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reindexPeople)
}

var reindexPeople = &cobra.Command{
	Use:   "reindex-people",
	Short: "Reindex people",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()

		services := newServices()

		return services.PeopleIndex.ReindexPeople(ctx, services.PeopleRepo.EachPerson)
	},
}
