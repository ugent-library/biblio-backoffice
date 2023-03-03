package cmd

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(GetAllPublicationsCmd)
}

var GetAllPublicationsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all publications",
	Run: func(cmd *cobra.Command, args []string) {
		GetAllPublications(cmd, args)
	},
}

func GetAllPublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.GetAllPublicationsRequest{}
	stream, err := c.GetAllPublications(context.Background(), req)
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

		if ge := res.GetError(); ge != nil {
			sre := status.FromProto(ge)
			cmd.Printf("%s\n", sre.Message())
		}

		if rr := res.GetPublication(); rr != nil {
			cmd.Printf("%s\n", rr.GetPayload())
		}
	}
}
