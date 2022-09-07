package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/client/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type RootCmd struct {
	Client       api.BiblioClient
	Marshaller   protojson.MarshalOptions
	Unmarshaller protojson.UnmarshalOptions
}

func (c *RootCmd) Wrap(fn func()) {
	// Set Marshaller
	c.Marshaller = protojson.MarshalOptions{
		UseProtoNames: true,
	}

	// Set Unmarshaller
	c.Unmarshaller = protojson.UnmarshalOptions{}

	// Set up the connection and the API client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", viper.GetString("api-host"), viper.GetInt("api-port"))
	log.Println(addr)

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth.BasicAuth{
			User:     viper.GetString("api-username"),
			Password: viper.GetString("api-password"),
		}),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	c.Client = api.NewBiblioClient(conn)

	// Run the command
	fn()
}
