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

func init() {
	PublicationCmd.AddCommand(UpdatePublicationCmd)
}

var UpdatePublicationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dataset",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(cmd.OutOrStdout())
		UpdatePublication(cmd, args)
	},
}

func UpdatePublication(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx, config)
	defer cnx.Close()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	p := &api.Publication{
		Payload: line,
	}

	req := &api.UpdatePublicationRequest{Publication: p}
	if _, err = c.UpdatePublication(ctx, req); err != nil {
		log.Fatal(err)
	}
}
