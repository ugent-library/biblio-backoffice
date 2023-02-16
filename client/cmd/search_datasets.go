package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	DatasetCmd.AddCommand(SearchDatasetsCmd)

	SearchDatasetsCmd.Flags().StringP("query", "q", "", "")
	SearchDatasetsCmd.Flags().String("limit", "", "")
	SearchDatasetsCmd.Flags().String("offset", "", "")
}

var SearchDatasetsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search datasets",
	Run: func(cmd *cobra.Command, args []string) {
		SearchDatasets(cmd, args)
	},
}

func SearchDatasets(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt32("limit")
	offset, _ := cmd.Flags().GetInt32("offset")

	req := &api.SearchDatasetsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.SearchDatasets(ctx, req)
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
	for i, h := range res.Hits {
		d := &models.Dataset{}
		if err := json.Unmarshal(h.Payload, d); err != nil {
			hits.Hits[i] = d
		}
	}

	j, err := json.Marshal(hits)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", j)
}
