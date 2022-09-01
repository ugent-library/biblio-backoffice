package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
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
	fileBufSize    = 524288
)

var (
	client api.BiblioClient

	marshaller = protojson.MarshalOptions{
		UseProtoNames: true,
	}

	unmarshaller = protojson.UnmarshalOptions{}
)

func main() {
	viper.SetEnvPrefix("biblio-backend")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("api-host", defaultAPIHost)
	viper.SetDefault("api-port", defaultAPIPort)

	rootCmd.PersistentFlags().String("api-host", defaultAPIHost, "api server host")
	rootCmd.PersistentFlags().Int("api-port", defaultAPIPort, "api server port")

	searchPublicationsCmd.Flags().StringP("query", "q", "", "")
	searchPublicationsCmd.Flags().StringP("limit", "", "", "")
	searchPublicationsCmd.Flags().StringP("offset", "", "", "")

	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(getFileCmd)
	fileCmd.AddCommand(addFileCmd)

	rootCmd.AddCommand(datasetCmd)
	datasetCmd.AddCommand(getDatasetCmd)
	datasetCmd.AddCommand(getAllDatasetsCmd)
	datasetCmd.AddCommand(searchDatasetsCmd)
	datasetCmd.AddCommand(updateDatasetCmd)
	datasetCmd.AddCommand(addDatasetsCmd)

	rootCmd.AddCommand(publicationCmd)
	publicationCmd.AddCommand(getPublicationCmd)
	publicationCmd.AddCommand(getAllPublicationsCmd)
	publicationCmd.AddCommand(searchPublicationsCmd)
	publicationCmd.AddCommand(updatePublicationCmd)
	publicationCmd.AddCommand(addPublicationsCmd)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := fmt.Sprintf("%s:%d", viper.GetString("api-host"), viper.GetInt("api-port"))
	log.Println(addr)
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

var fileCmd = &cobra.Command{
	Use:   "file [command]",
	Short: "File commands",
}

var getFileCmd = &cobra.Command{
	Use:   "get [sha256]",
	Short: "Get file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := &api.GetFileRequest{Sha256: args[0]}
		stream, err := client.GetFile(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}

			if _, err := os.Stdout.Write(res.Chunk); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var addFileCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file",
	Run: func(cmd *cobra.Command, args []string) {
		stream, err := client.AddFile(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, fileBufSize)

		for {
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("cannot read chunk to buffer: ", err)
			}

			req := &api.AddFileRequest{Chunk: buf[:n]}

			if err = stream.Send(req); err != nil {
				log.Fatal("cannot send chunk to server: ", err)
			}
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.WriteString(res.Sha256 + "\n")
	},
}

var datasetCmd = &cobra.Command{
	Use:   "dataset [command]",
	Short: "Dataset commands",
}

var getDatasetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get dataset by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		id := args[0]
		req := &api.GetDatasetRequest{Id: id}
		res, err := client.GetDataset(ctx, req)
		if err != nil {
			log.Fatal(err)
		}

		j, err := marshaller.Marshal(res.Dataset)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	},
}

var getAllDatasetsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		req := &api.GetAllDatasetsRequest{}
		stream, err := client.GetAllDatasets(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}

			j, err := marshaller.Marshal(res.Dataset)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", j)
		}
	},
}

var searchDatasetsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search datasets",
	Run: func(cmd *cobra.Command, args []string) {
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
		res, err := client.SearchDatasets(ctx, req)
		if err != nil {
			log.Fatal(err)
		}

		j, err := marshaller.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	},
}

var updateDatasetCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dataset",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}

		dataset := &api.Dataset{}
		if err := unmarshaller.Unmarshal(line, dataset); err != nil {
			log.Fatal(err)
		}

		req := &api.UpdateDatasetRequest{Dataset: dataset}
		if _, err = client.UpdateDataset(ctx, req); err != nil {
			log.Fatal(err)
		}
	},
}

var addDatasetsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add datasets",
	Run: func(cmd *cobra.Command, args []string) {
		stream, err := client.AddDatasets(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		waitc := make(chan struct{})

		go func() {
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}
				if err != nil {
					log.Fatal(err)
				}
				log.Println(res.Messsage)
			}
		}()

		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			dataset := &api.Dataset{}
			if err := unmarshaller.Unmarshal(line, dataset); err != nil {
				log.Fatal(err)
			}

			req := &api.AddDatasetsRequest{Dataset: dataset}
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
			}
		}

		stream.CloseSend()
		<-waitc
	},
}

var publicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}

var getPublicationCmd = &cobra.Command{
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

		j, err := marshaller.Marshal(res.Publication)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	},
}

var getAllPublicationsCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all publications",
	Run: func(cmd *cobra.Command, args []string) {
		req := &api.GetAllPublicationsRequest{}
		stream, err := client.GetAllPublications(context.Background(), req)
		if err != nil {
			log.Fatal(err)
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

var searchPublicationsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search publications",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt32("limit")
		offset, _ := cmd.Flags().GetInt32("offset")

		req := &api.SearchPublicationsRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		}
		res, err := client.SearchPublications(ctx, req)
		if err != nil {
			log.Fatal(err)
		}

		j, err := marshaller.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", j)
	},
}

var updatePublicationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update publication",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}

		pub := &api.Publication{}
		if err := unmarshaller.Unmarshal(line, pub); err != nil {
			log.Fatal(err)
		}

		req := &api.UpdatePublicationRequest{Publication: pub}
		if _, err = client.UpdatePublication(ctx, req); err != nil {
			log.Fatal(err)
		}
	},
}

var addPublicationsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add publications",
	Run: func(cmd *cobra.Command, args []string) {
		stream, err := client.AddPublications(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		waitc := make(chan struct{})

		go func() {
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}
				if err != nil {
					log.Fatal(err)
				}
				log.Println(res.Messsage)
			}
		}()

		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			pub := &api.Publication{}
			if err := unmarshaller.Unmarshal(line, pub); err != nil {
				log.Fatal(err)
			}

			req := &api.AddPublicationsRequest{Publication: pub}
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
			}
		}

		stream.CloseSend()
		<-waitc
	},
}
