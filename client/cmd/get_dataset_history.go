package cmd

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

var GetDatasetHistoryCmd = &cobra.Command{
	Use:   "get-history [id]",
	Short: "Get dataset history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetDatasetHistory(cmd, args)
	},
}

func GetDatasetHistory(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	req := &api.GetDatasetHistoryRequest{Id: args[0]}
	stream, err := c.GetDatasetHistory(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		cmd.Printf("%s\n", res.Dataset.Payload)
	}
}
