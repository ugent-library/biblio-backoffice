package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type PurgeAllDatasetsCmd struct {
	RootCmd
}

func (c *PurgeAllDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge-all",
		Short: "Purge all datasets",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "are you sure?")

	return cmd
}

func (c *PurgeAllDatasetsCmd) Run(cmd *cobra.Command, args []string) {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		return
	}

	req := &api.PurgeAllDatasetsRequest{}
	if _, err := c.Client.PurgeAllDatasets(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
