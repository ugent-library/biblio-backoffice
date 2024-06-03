package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	api "github.com/ugent-library/biblio-backoffice/api/v1"
	"github.com/ugent-library/biblio-backoffice/backends"
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

func New(services *backends.Services, users Users) *grpc.Server {
	logger, _ := zap.NewProduction()

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

	permissions := listPermissions()
	basicAuthInterceptor := NewBasicAuthInterceptor(users, permissions)

	gsrv := grpc.NewServer(
		grpc.Creds(nil),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, zap_opt),
			grpc.StreamServerInterceptor(basicAuthInterceptor.Stream()),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, zap_opt),
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
