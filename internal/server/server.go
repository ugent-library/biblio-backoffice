package server

import (
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/viper"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/internal/backends"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

type server struct {
	api.UnimplementedBiblioServer
	services *backends.Services
}

func listUsers() Users {
	var u Users

	u = append(u, &User{
		Username: viper.GetString("admin_username"),
		Password: viper.GetString("admin_password"),
		Role:     "admin",
	})

	u = append(u, &User{
		Username: viper.GetString("curator_username"),
		Password: viper.GetString("curator_password"),
		Role:     "curator",
	})

	return u
}

func listPermissions() map[string][]string {
	const biblioServicePath = "/biblio.v1.Biblio/"

	return map[string][]string{
		biblioServicePath + "AddDatasets":           {"admin"},
		biblioServicePath + "AddFile":               {"admin"},
		biblioServicePath + "AddPublications":       {"admin"},
		biblioServicePath + "CleanupPublications":   {"admin"},
		biblioServicePath + "ExistsFile":            {"admin", "curator"},
		biblioServicePath + "GetAllDatasets":        {"admin", "curator"},
		biblioServicePath + "GetAllPublications":    {"admin", "curator"},
		biblioServicePath + "GetDataset":            {"admin", "curator"},
		biblioServicePath + "GetDatasetHistory":     {"admin", "curator"},
		biblioServicePath + "GetFile":               {"admin", "curator"},
		biblioServicePath + "GetPublication":        {"admin", "curator"},
		biblioServicePath + "GetPublicationHistory": {"admin", "curator"},
		biblioServicePath + "ImportDatasets":        {"admin"},
		biblioServicePath + "ImportPublications":    {"admin"},
		biblioServicePath + "PurgeAllDatasets":      {"admin"},
		biblioServicePath + "PurgeAllPublications":  {"admin"},
		biblioServicePath + "PurgeDataset":          {"admin"},
		biblioServicePath + "PurgePublication":      {"admin"},
		biblioServicePath + "ReindexDatasets":       {"admin"},
		biblioServicePath + "ReindexPublications":   {"admin"},
		biblioServicePath + "Relate":                {"admin"},
		biblioServicePath + "SearchDatasets":        {"admin", "curator"},
		biblioServicePath + "SearchPublications":    {"admin", "curator"},
		biblioServicePath + "TransferPublications":  {"admin"},
		biblioServicePath + "UpdateDataset":         {"admin"},
		biblioServicePath + "UpdatePublication":     {"admin"},
		biblioServicePath + "ValidateDatasets":      {"admin", "curator"},
		biblioServicePath + "ValidatePublications":  {"admin", "curator"},
	}
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

	users := listUsers()
	permissions := listPermissions()
	basicAuthInterceptor := NewBasicAuthInterceptor(users, permissions)

	gsrv := grpc.NewServer(
		tlsOption,
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger.Desugar(), zap_opt),
			grpc.StreamServerInterceptor(basicAuthInterceptor.Stream()),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger.Desugar(), zap_opt),
			grpc.UnaryServerInterceptor(basicAuthInterceptor.Unary()),
		),
	)

	srv := &server{
		services: services,
	}

	// Enable the gRPC reflection API
	// e.g. grpcurl host:port list -> list all available services & methods
	reflection.Register(gsrv)

	api.RegisterBiblioServer(gsrv, srv)
	return gsrv
}
