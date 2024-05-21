package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reindexProjects)
}

var reindexProjects = &cobra.Command{
	Use:   "reindex-projects",
	Short: "Reindex projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()

		services := newServices()

		return services.ProjectsIndex.ReindexProjects(ctx, services.ProjectsRepo.EachProject)
	},
}
