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
	DatasetCmd.AddCommand(GetDatasetCmd)
}

var GetDatasetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get dataset by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetDataset(cmd, args)
	},
}

func GetDataset(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	id := args[0]
	req := &api.GetDatasetRequest{Id: id}
	res, err := c.GetDataset(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Printf("%s\n", res.Dataset.Payload)
}
