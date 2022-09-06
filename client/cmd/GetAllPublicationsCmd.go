package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type GetAllPublicationsCmd struct {
	RootCmd
}

func (c *GetAllPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get all publications",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *GetAllPublicationsCmd) Run(args []string) {
	req := &api.GetAllPublicationsRequest{}
	stream, err := c.Client.GetAllPublications(context.Background(), req)
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

		j, err := c.Marshaller.Marshal(res.Publication)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	}
}
