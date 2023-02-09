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

func init() {
	DatasetCmd.AddCommand(GetAllDatasetsCmd)
}

var GetAllDatasetsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		GetAllDatasets(cmd, args)
	},
}

func GetAllDatasets(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx, config)
	defer cnx.Close()

	req := &api.GetAllDatasetsRequest{}
	stream, err := c.GetAllDatasets(context.Background(), req)
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
