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

type UpdatePublicationCmd struct {
	RootCmd
}

func (c *UpdatePublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dataset",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *UpdatePublicationCmd) Run(args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	pub := &api.Publication{}
	if err := c.Unmarshaller.Unmarshal(line, pub); err != nil {
		log.Fatal(err)
	}

	req := &api.UpdatePublicationRequest{Publication: pub}
	if _, err = c.Client.UpdatePublication(ctx, req); err != nil {
		log.Fatal(err)
	}
}
