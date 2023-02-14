package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func Create(ctx context.Context, config Config) (api.BiblioClient, *grpc.ClientConn, error) {
	// Set encryption
	var dialOptionSecureConn grpc.DialOption
	if viper.GetBool("insecure") {
		dialOptionSecureConn = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		dialOptionSecureConn = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	}

	// Set up the connection and the API client with Basic Authentication
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := grpc.DialContext(ctx, addr,
		dialOptionSecureConn,
		grpc.WithPerRPCCredentials(auth.BasicAuth{
			User:     config.Username,
			Password: config.Password,
		}),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, nil, err
	}

	client := api.NewBiblioClient(conn)

	return client, conn, nil
}
