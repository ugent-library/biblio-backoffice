package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/client"
	"google.golang.org/grpc/status"
)

func init() {
	PublicationCmd.AddCommand(GetPublicationCmd)
}

var GetPublicationCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)
		log.SetOutput(cmd.OutOrStdout())
		GetPublication(cmd, args)
	},
}

func GetPublication(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, cnx, err := client.Create(ctx, config)
	defer cnx.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("ContextDeadlineExceeded: true")
	}

	id := args[0]
	req := &api.GetPublicationRequest{Id: id}
	res, err := c.GetPublication(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatal(err)
		}
		cmd.Println(st.Message())
	} else {
		cmd.Printf("%s\n", res.Publication.Payload)
	}
}
