package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	api "github.com/ugent-library/biblio-backend/api/v1"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"google.golang.org/grpc"
)

type server struct {
	api.UnimplementedBiblioServer
	services *backends.Services
}

func New(services *backends.Services) *grpc.Server {
	gsrv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	srv := &server{services: services}
	api.RegisterBiblioServer(gsrv, srv)
	return gsrv
}
