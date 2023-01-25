package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/models"
	"github.com/ugent-library/biblio-backoffice/internal/server"
)

type SearchPublicationsCmd struct {
	RootCmd
}

func (c *SearchPublicationsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search publications",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	cmd.Flags().StringP("query", "q", "", "")
	cmd.Flags().StringP("limit", "", "", "")
	cmd.Flags().StringP("offset", "", "", "")

	return cmd
}

func (c *SearchPublicationsCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt32("limit")
	offset, _ := cmd.Flags().GetInt32("offset")

	req := &api.SearchPublicationsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.Client.SearchPublications(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	hits := struct {
		Offset, Limit, Total int32
		Hits                 []*models.Publication
	}{
		Offset: res.Offset,
		Limit:  res.Limit,
		Total:  res.Total,
		Hits:   make([]*models.Publication, len(res.Hits)),
	}
	for i, p := range res.Hits {
		hits.Hits[i] = server.MessageToPublication(p)
	}

	j, err := json.Marshal(hits)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
