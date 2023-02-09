package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	DatasetCmd.AddCommand(ImportDatasetsCmd)
	PurgeAllDatasetsCmd.Flags().BoolP("yes", "y", false, "are you sure?")
}

var PurgeAllDatasetsCmd = &cobra.Command{
	Use:   "purge-all",
	Short: "Purge all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		PurgeAllDatasets(cmd, args)
	},
}

func PurgeAllDatasets(cmd *cobra.Command, args []string) {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx, config)
	defer cnx.Close()

	req := &api.PurgeAllDatasetsRequest{}
	if _, err := c.PurgeAllDatasets(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
