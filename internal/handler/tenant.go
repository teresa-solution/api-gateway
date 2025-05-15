package handler

import (
	"context"
	_ "net/http"
	"time" // Add this import

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	tenantpb "github.com/teresa-solution/api-gateway/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RegisterHandlers registers the HTTP handlers for the API Gateway
func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, grpcAddr string) error {
	// Add a timeout to the context for dialing
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Wait for the connection to be established
	)
	if err != nil {
		log.Error().Err(err).Str("grpc_addr", grpcAddr).Msg("Failed to dial gRPC server")
		return err
	}

	// Register the gRPC-Gateway handler
	if err := tenantpb.RegisterTenantServiceHandler(ctx, mux, conn); err != nil {
		log.Error().Err(err).Msg("Failed to register gRPC gateway handlers")
		return err
	}

	// Log successful connection
	log.Info().Str("grpc_addr", grpcAddr).Msg("Successfully connected to gRPC server")
	return nil
}
