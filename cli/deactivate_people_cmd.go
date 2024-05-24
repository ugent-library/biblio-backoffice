package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deactivatePeople)
}

var deactivatePeople = &cobra.Command{
	Use:   "deactivate-people",
	Short: "Deactivate people",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()

		services := newServices()

		return services.PeopleRepo.DeactivatePeople(ctx)
	},
}
