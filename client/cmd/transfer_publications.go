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
	PublicationCmd.AddCommand(TransferPublicationsCmd)
}

var TransferPublicationsCmd = &cobra.Command{
	Use:   "transfer UID UID [PUBID]",
	Short: "Transfer publications between people",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		TransferPublications(cmd, args)
	},
}

func TransferPublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	source := args[0]
	dest := args[1]

	pubid := ""
	if len(args) > 2 {
		pubid = args[2]
	}

	req := &api.TransferPublicationsRequest{
		Src:           source,
		Dest:          dest,
		Publicationid: pubid,
	}
	stream, err := c.TransferPublications(context.Background(), req)
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
