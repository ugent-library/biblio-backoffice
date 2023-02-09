package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
)

func init() {
	PublicationCmd.AddCommand(GetPublicationCmd)
}

var GetPublicationCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(cmd.OutOrStdout())
		GetPublication(cmd, args)
	},
}

func GetPublication(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx := client.Create(ctx, config)
	defer cnx.Close()

	id := args[0]
	req := &api.GetPublicationRequest{Id: id}
	res, err := c.GetPublication(ctx, req)
	if err != nil {
		cmd.Println(err)
		// log.Fatal(err)
	} else {
		cmd.Printf("%s\n", res.Publication.Payload)
	}
}
