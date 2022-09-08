package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/client/auth"
	"github.com/ugent-library/biblio-backend/client/tls"
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

	dialOptionSecureConn := grpc.WithTransportCredentials(insecure.NewCredentials())
	if !viper.GetBool("api-tls-disabled") {
		tlsCredentials, err := tls.LoadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}
		dialOptionSecureConn = grpc.WithTransportCredentials(tlsCredentials)
	}

	// Set up the connection and the API client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", viper.GetString("api-host"), viper.GetInt("api-port"))
	log.Println(addr)

	conn, err := grpc.DialContext(ctx, addr,
		dialOptionSecureConn,
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
