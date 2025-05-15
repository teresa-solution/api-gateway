package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/teresa-solution/api-gateway/internal/handler"
	"github.com/teresa-solution/api-gateway/internal/middleware"
	"github.com/teresa-solution/api-gateway/internal/monitoring"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		httpPort = flag.Int("http-port", 8080, "HTTP server port")
		grpcAddr = flag.String("grpc-addr", "127.0.0.1:50051", "gRPC server address") // Define grpcAddr
	)
	flag.Parse()

	// Initialize metrics
	monitoring.InitMetrics()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	if err := handler.RegisterHandlers(ctx, mux, *grpcAddr); err != nil { // Pass grpcAddr
		log.Fatal().Err(err).Msg("Failed to register handlers")
	}

	// Apply middleware
	handlerWithAuth := middleware.AuthMiddleware(mux)
	handlerWithRateLimit := middleware.RateLimitMiddleware(handlerWithAuth)
	handlerWithMetrics := middleware.MetricsMiddleware(handlerWithRateLimit)

	// Create HTTP server with metrics endpoint
	finalMux := http.NewServeMux()
	finalMux.Handle("/", handlerWithMetrics)
	finalMux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: finalMux,
	}

	go func() {
		log.Info().Msgf("API Gateway listening on port %d", *httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down API Gateway...")

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}
	log.Info().Msg("API Gateway exiting")
}
