package cmd

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
)

type UpdatePublicationCmd struct {
	RootCmd
}

func (c *UpdatePublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dataset",
		Run: func(cmd *cobra.Command, args []string) {
			log.SetOutput(cmd.OutOrStdout())

			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *UpdatePublicationCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	p := &api.Publication{
		Payload: line,
	}

	req := &api.UpdatePublicationRequest{Publication: p}
	if _, err = c.Client.UpdatePublication(ctx, req); err != nil {
		log.Fatal(err)
	}
}
