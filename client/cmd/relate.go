package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type PublicationRelateDatasetCmd struct {
	RootCmd
}

func (c *PublicationRelateDatasetCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relate-dataset [id] [dataset-id]",
		Short: "Add related dataset to publication",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *PublicationRelateDatasetCmd) Run(cmd *cobra.Command, args []string) {
	req := &api.RelateRequest{
		One: &api.RelateRequest_PublicationOne{PublicationOne: args[0]},
		Two: &api.RelateRequest_DatasetTwo{DatasetTwo: args[1]},
	}
	if _, err := c.Client.Relate(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
