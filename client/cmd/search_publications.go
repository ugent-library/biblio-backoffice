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
	PublicationCmd.AddCommand(SearchPublicationsCmd)
	SearchPublicationsCmd.Flags().StringP("query", "q", "", "")
	SearchPublicationsCmd.Flags().StringP("limit", "", "", "")
	SearchPublicationsCmd.Flags().StringP("offset", "", "", "")
}

var SearchPublicationsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search publications",
	Run: func(cmd *cobra.Command, args []string) {
		SearchPublications(cmd, args)
	},
}

func SearchPublications(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt32("limit")
	offset, _ := cmd.Flags().GetInt32("offset")

	req := &api.SearchPublicationsRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	res, err := c.SearchPublications(ctx, req)
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
}
