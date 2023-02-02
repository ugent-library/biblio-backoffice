package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/server"
)

type UpdatePublicationCmd struct {
	RootCmd
}

func (c *UpdatePublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dataset",
		Run: func(cmd *cobra.Command, args []string) {
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

	p := &models.Publication{}
	if err := json.Unmarshal(line, p); err != nil {
		log.Fatal(err)
	}

	req := &api.UpdatePublicationRequest{Publication: server.PublicationToMessage(p)}
	if _, err = c.Client.UpdatePublication(ctx, req); err != nil {
		log.Fatal(err)
	}
}
