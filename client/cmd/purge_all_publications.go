package cmd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(PurgeAllPublicationsCmd)
}

var PurgeAllPublicationsCmd = &cobra.Command{
	Use:   "purge-all",
	Short: "Purge all publications",
	Run: func(cmd *cobra.Command, args []string) {
		PurgeAllPublications(cmd, args)
	},
}

func init() {
	PurgeAllPublicationsCmd.Flags().BoolP("yes", "y", false, "are you sure?")
}

func PurgeAllPublications(cmd *cobra.Command, args []string) {
	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.PurgeAllPublicationsRequest{}
	if _, err := c.PurgeAllPublications(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
