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

type UpdateDatasetCmd struct {
	RootCmd
}

func (c *UpdateDatasetCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update publication",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	return cmd
}

func (c *UpdateDatasetCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	d := &api.Dataset{
		Payload: line,
	}


	req := &api.UpdateDatasetRequest{Dataset: d}
	if _, err = c.Client.UpdateDataset(ctx, req); err != nil {
		log.Fatal(err)
	}
}
