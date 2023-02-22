package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	DatasetCmd.AddCommand(ReindexDatasetCmd)
}

var ReindexDatasetCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		ReindexDatasets(cmd, args)
	},
}

func ReindexDatasets(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.ReindexDatasetsRequest{}
	stream, err := c.ReindexDatasets(context.Background(), req)
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

		fmt.Printf("%s\n", res.Message)
	}
}
