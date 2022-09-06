package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
)

type SearchDatasetsCmd struct {
	RootCmd
}

func (c *SearchDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search datasets",
		Run: func(_ *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(args)
			})
		},
	}

	return cmd
}

func (c *SearchDatasetsCmd) Run(args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO these don't work? Test 'em!
	query, _ := c.Command().Flags().GetString("query")
	limit, _ := c.Command().Flags().GetInt32("limit")
	offset, _ := c.Command().Flags().GetInt32("offset")

	req := &api.SearchDatasetsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.Client.SearchDatasets(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	j, err := c.Marshaller.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
