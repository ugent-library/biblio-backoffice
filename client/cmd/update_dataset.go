package cmd

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type UpdateDatasetCmd struct {
	RootCmd
}

func (c *UpdateDatasetCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update publication",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *UpdateDatasetCmd) Run(args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	dataset := &api.Dataset{}
	if err := c.Unmarshaller.Unmarshal(line, dataset); err != nil {
		log.Fatal(err)
	}

	req := &api.UpdateDatasetRequest{Dataset: dataset}
	if _, err = c.Client.UpdateDataset(ctx, req); err != nil {
		log.Fatal(err)
	}
}
