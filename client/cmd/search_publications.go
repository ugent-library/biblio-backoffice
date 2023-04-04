package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	cnx "github.com/ugent-library/biblio-backoffice/client/connection"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	PublicationCmd.AddCommand(SearchPublicationsCmd)
	SearchPublicationsCmd.Flags().StringP("query", "q", "", "")
	SearchPublicationsCmd.Flags().StringP("limit", "", "", "")
	SearchPublicationsCmd.Flags().StringP("offset", "", "", "")
}

var SearchPublicationsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search publications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return SearchPublications(cmd, args)
	},
}

func SearchPublications(cmd *cobra.Command, args []string) error {
	return cnx.Handle(config, func(c api.BiblioClient) error {
		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt32("limit")
		offset, _ := cmd.Flags().GetInt32("offset")

		req := &api.SearchPublicationsRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		}
		res, err := c.SearchPublications(context.Background(), req)
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
		for i, h := range res.Hits {
			p := &models.Publication{}
			if err := json.Unmarshal(h.Payload, p); err != nil {
				hits.Hits[i] = p
			}
		}

		j, err := json.Marshal(hits)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)

		return nil
	})
}
