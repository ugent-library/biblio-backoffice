package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(GetPublicationHistoryCmd)
}

var GetPublicationHistoryCmd = &cobra.Command{
	Use:   "get-history [id]",
	Short: "Get publication history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetPublicationHistory(cmd, args)
	},
}

func GetPublicationHistory(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx, config)
	defer cnx.Close()

	req := &api.GetPublicationHistoryRequest{Id: args[0]}
	stream, err := c.GetPublicationHistory(context.Background(), req)
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

		fmt.Printf("%s\n", res.Publication.Payload)
	}
}
