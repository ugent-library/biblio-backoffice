package client

import (
	"context"
	"fmt"
	"time"

	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/client/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Transmit(config Config, t func(c api.BiblioClient) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Set TLS encryption
	var dialOptionSecureConn grpc.DialOption
	if config.Insecure {
		dialOptionSecureConn = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		creds, err := LoadTLSCredentials(config)
		if err != nil {
			return err
		}
		dialOptionSecureConn = grpc.WithTransportCredentials(creds)
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
		return err
	}

	defer conn.Close()

	client := api.NewBiblioClient(conn)

	err = t(client)
	if err != nil {
		return err
	}

	return nil
}
