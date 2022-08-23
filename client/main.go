package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultAPIHost = ""
	defaultAPIPort = 30000
)

var client api.BiblioClient

func main() {
	viper.SetEnvPrefix("biblio-backend")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("api-host", defaultAPIHost)
	viper.SetDefault("api-port", defaultAPIPort)

	rootCmd.Flags().String("api-host", defaultAPIHost, "api server host")
	rootCmd.Flags().Int("api-port", defaultAPIPort, "api server port")

	rootCmd.AddCommand(publicationCmd)
	publicationCmd.AddCommand(publicationGetCmd)
	publicationCmd.AddCommand(publicationGetAllCmd)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := fmt.Sprintf("%s:%d", viper.GetString("api-host"), viper.GetInt("api-port"))
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client = api.NewBiblioClient(conn)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// TODO we shouldn't do this for all flags, only ones that have a config equivalent
var rootCmd = &cobra.Command{
	Use:   "biblio [command]",
	Short: "biblio client",
	// flags override env vars
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				viper.Set(f.Name, f.Value.String())
			}
		})
		return nil
	},
}

var publicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}

var publicationGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publication by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		id := args[0]
		req := &api.GetPublicationRequest{Id: id}
		res, err := client.GetPublication(ctx, req)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", res)
	},
}

var publicationGetAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all publications",
	Run: func(cmd *cobra.Command, args []string) {
		req := &api.GetAllPublicationsRequest{}
		stream, err := client.GetAllPublications(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		marshaller := protojson.MarshalOptions{
			UseProtoNames: true,
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}

			j, err := marshaller.Marshal(res.Publication)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", j)
		}
	},
}
