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
	PublicationCmd.AddCommand(CleanupPublicationsCmd)
}

var CleanupPublicationsCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Make publications consistent, clean up data anomalies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CleanupPublications(cmd, args)
	},
}

func CleanupPublications(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	req := &api.CleanupPublicationsRequest{}
	stream, err := c.CleanupPublications(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error while reading stream: %v", err)
		}

		cmd.Printf("%s\n", res.Message)
	}

	return nil
}
