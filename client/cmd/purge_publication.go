package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(PurgePublicationCmd)
}

var PurgePublicationCmd = &cobra.Command{
	Use:   "purge [id]",
	Short: "Purge publication",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		PurgePublication(cmd, args)
	},
}

func PurgePublication(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	id := args[0]
	req := &api.PurgePublicationRequest{Id: id}
	if _, err := c.PurgePublication(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
