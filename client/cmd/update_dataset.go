package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/server"
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

	d := &models.Dataset{}
	if err := json.Unmarshal(line, d); err != nil {
		log.Fatal(err)
	}

	req := &api.UpdateDatasetRequest{Dataset: server.DatasetToMessage(d)}
	if _, err = c.Client.UpdateDataset(ctx, req); err != nil {
		log.Fatal(err)
	}
}
