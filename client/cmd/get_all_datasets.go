package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type GetAllDatasetsCmd struct {
	RootCmd
}

func (c *GetAllDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get all datasets",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *GetAllDatasetsCmd) Run(args []string) {

	req := &api.GetAllDatasetsRequest{}
	stream, err := c.Client.GetAllDatasets(context.Background(), req)
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

		j, err := c.Marshaller.Marshal(res.Dataset)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	}

}
