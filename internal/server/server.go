package server

import (
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type server struct {
	api.UnimplementedBiblioServer
	services *backends.Services
}

func New(services *backends.Services, logger *zap.SugaredLogger) *grpc.Server {
	zap_opt := grpc_zap.WithLevels(
		func(c codes.Code) zapcore.Level {
			var l zapcore.Level
			switch c {
			case codes.OK:
				l = zapcore.InfoLevel

			case codes.Internal:
				l = zapcore.ErrorLevel

			default:
				l = zapcore.DebugLevel
			}
			return l
		},
	)

	// Defaults to an insecure connection
	tlsOption := grpc.Creds(nil)

	// If set, enabled server-side TLS secure connection
	if viper.GetBool("api-tls-enabled") {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}

		tlsOption = grpc.Creds(tlsCredentials)
	}

	gsrv := grpc.NewServer(
		tlsOption,
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger.Desugar(), zap_opt),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger.Desugar(), zap_opt),
			grpc_auth.UnaryServerInterceptor(auth),
		),
	)

	srv := &server{services: services}
	api.RegisterBiblioServer(gsrv, srv)
	return gsrv
}
