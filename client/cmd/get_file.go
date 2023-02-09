package cmd

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

var GetFileCMd = &cobra.Command{
	Use:   "get [sha256]",
	Short: "Get file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetFile(cmd, args)
	},
}

func GetFile(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	req := &api.GetFileRequest{Sha256: args[0]}
	stream, err := c.GetFile(context.Background(), req)
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
