package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type PurgePublicationCmd struct {
	RootCmd
}

func (c *PurgePublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge [id]",
		Short: "Purge publication",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *PurgePublicationCmd) Run(cmd *cobra.Command, args []string) {
	id := args[0]
	req := &api.PurgePublicationRequest{Id: id}
	if _, err := c.Client.PurgePublication(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
