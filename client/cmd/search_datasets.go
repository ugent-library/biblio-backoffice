package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return SearchDatasets(cmd, args)
	},
}

func SearchDatasets(cmd *cobra.Command, args []string) error {
	err := cnx.Handle(config, func(c api.BiblioClient) error {
		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt32("limit")
		offset, _ := cmd.Flags().GetInt32("offset")

		req := &api.SearchDatasetsRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		}
		res, err := c.SearchDatasets(context.Background(), req)
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

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	return err
}
