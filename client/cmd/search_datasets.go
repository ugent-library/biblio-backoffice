package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/server"
)

type SearchDatasetsCmd struct {
	RootCmd
}

func (c *SearchDatasetsCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search datasets",
		Run: func(cmd *cobra.Command, args []string) {
			c.Wrap(func() {
				c.Run(cmd, args)
			})
		},
	}

	cmd.Flags().StringP("query", "q", "", "")
	cmd.Flags().String("limit", "", "")
	cmd.Flags().String("offset", "", "")

	return cmd
}

func (c *SearchDatasetsCmd) Run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt32("limit")
	offset, _ := cmd.Flags().GetInt32("offset")

	req := &api.SearchDatasetsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.Client.SearchDatasets(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	hits := struct {
		Offset, Limit, Total int32
		Hits                 []*models.Dataset
	}{
		Offset: res.Offset,
		Limit:  res.Limit,
		Total:  res.Total,
		Hits:   make([]*models.Dataset, len(res.Hits)),
	}
	for i, d := range res.Hits {
		hits.Hits[i] = server.MessageToDataset(d)
	}

	j, err := json.Marshal(hits)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
