package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type PurgeAllPublicationsCmd struct {
	RootCmd
}

func (c *PurgeAllPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge-all",
		Short: "Purge all publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *PurgeAllPublicationsCmd) Run(cmd *cobra.Command, args []string) {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		return
	}

	req := &api.PurgeAllPublicationsRequest{}
	if _, err := c.Client.PurgeAllPublications(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
