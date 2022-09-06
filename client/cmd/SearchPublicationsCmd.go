package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type SearchPublicationsCmd struct {
	RootCmd
}

func (c *SearchPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search publications",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *SearchPublicationsCmd) Run(args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO Fix me
	query, _ := c.Command().Flags().GetString("query")
	limit, _ := c.Command().Flags().GetInt32("limit")
	offset, _ := c.Command().Flags().GetInt32("offset")

	req := &api.SearchPublicationsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.Client.SearchPublications(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	j, err := c.Marshaller.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
