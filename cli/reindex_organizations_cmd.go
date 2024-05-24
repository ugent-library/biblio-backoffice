package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reindexOrganizations)
}

var reindexOrganizations = &cobra.Command{
	Use:   "reindex-organizations",
	Short: "Reindex organizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()

		services := newServices()

		return services.PeopleIndex.ReindexOrganizations(ctx, services.PeopleRepo.EachOrganization)
	},
}
