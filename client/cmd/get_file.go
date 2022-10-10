package cmd

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type GetFileCMd struct {
	RootCmd
}

func (c *GetFileCMd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [sha256]",
		Short: "Get file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *GetFileCMd) Run(cmd *cobra.Command, args []string) {
	req := &api.GetFileRequest{Sha256: args[0]}
	stream, err := c.Client.GetFile(context.Background(), req)
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

		if _, err := os.Stdout.Write(res.Chunk); err != nil {
			log.Fatal(err)
		}
	}
}
