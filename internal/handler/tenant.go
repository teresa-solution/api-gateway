package handler

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	tenantpb "github.com/teresa-solution/api-gateway/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RegisterHandlers registers the HTTP handlers for the API Gateway
func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, grpcAddr string) error {
	// Register the gRPC-Gateway handler
	if err := tenantpb.RegisterTenantServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		log.Error().Err(err).Msg("Failed to register gRPC gateway handlers")
		return err
	}

	log.Info().Str("grpc_addr", grpcAddr).Msg("Successfully connected to gRPC server")
	return nil
}
