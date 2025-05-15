package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/teresa-solution/api-gateway/internal/handler"
	"github.com/teresa-solution/api-gateway/internal/middleware"
	"github.com/teresa-solution/api-gateway/internal/monitoring"
	"google.golang.org/grpc/metadata"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		httpPort = flag.Int("http-port", 8080, "HTTP server port")
		grpcAddr = flag.String("grpc-addr", "127.0.0.1:50051", "gRPC server address")
	)
	flag.Parse()

	monitoring.InitMetrics()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.ToLower(key) == "x-tenant-subdomain" {
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md, _ := metadata.FromIncomingContext(ctx)
			if md == nil {
				md = metadata.MD{}
			}
			if subdomain := r.Header.Get("X-Tenant-Subdomain"); subdomain != "" {
				md.Set("x-tenant-subdomain", subdomain)
			}
			return md
		}),
	)

	if err := handler.RegisterHandlers(ctx, mux, *grpcAddr); err != nil {
		log.Fatal().Err(err).Msg("Failed to register handlers")
	}

	handlerWithAuth := middleware.AuthMiddleware(mux)
	handlerWithRateLimit := middleware.RateLimitMiddleware(handlerWithAuth)
	handlerWithMetrics := middleware.MetricsMiddleware(handlerWithRateLimit)

	finalMux := http.NewServeMux()
	finalMux.Handle("/", handlerWithMetrics)
	finalMux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: finalMux,
	}

	go func() {
		log.Info().Msgf("API Gateway listening on port %d (HTTPS)", *httpPort)
		// Use ListenAndServeTLS for HTTPS
		if err := server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTPS server failed")
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
