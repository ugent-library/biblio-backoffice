package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

var PublicationRelateDatasetCmd = &cobra.Command{
	Use:   "relate-dataset [id] [dataset-id]",
	Short: "Add related dataset to publication",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		PublicationRelateDataset(cmd, args)
	},
}

func PublicationRelateDataset(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	req := &api.RelateRequest{
		One: &api.RelateRequest_PublicationOne{PublicationOne: args[0]},
		Two: &api.RelateRequest_DatasetTwo{DatasetTwo: args[1]},
	}
	if _, err := c.Relate(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
