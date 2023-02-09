package cmd

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

var UpdateDatasetCmd = &cobra.Command{
	Use:   "update",
	Short: "Update publication",
	Run: func(cmd *cobra.Command, args []string) {
		UpdateDataset(cmd, args)
	},
}

func UpdateDataset(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx)
	defer cnx.Close()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	d := &api.Dataset{
		Payload: line,
	}

	req := &api.UpdateDatasetRequest{Dataset: d}
	if _, err = c.UpdateDataset(ctx, req); err != nil {
		log.Fatal(err)
	}
}
