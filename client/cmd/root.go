package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/client/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaller   = protojson.MarshalOptions{UseProtoNames: true}
	unmarshaller = protojson.UnmarshalOptions{}
)

type RootCmd struct {
	Client       api.BiblioClient
	Marshaller   protojson.MarshalOptions
	Unmarshaller protojson.UnmarshalOptions
}

func (c *RootCmd) Wrap(fn func()) {
	// Set marshaller
	c.Marshaller = marshaller
	// Set unmarshaller
	c.Unmarshaller = unmarshaller

	// Set encryption
	var dialOptionSecureConn grpc.DialOption
	if viper.GetBool("insecure") {
		dialOptionSecureConn = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		dialOptionSecureConn = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	}

	// Set up the connection and the API client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	log.Println(addr)

	conn, err := grpc.DialContext(ctx, addr,
		dialOptionSecureConn,
		grpc.WithPerRPCCredentials(auth.BasicAuth{
			User:     viper.GetString("username"),
			Password: viper.GetString("password"),
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
