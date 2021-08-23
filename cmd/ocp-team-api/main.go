package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-team-api/internal/api"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPort           = ":8082"
	grpcServerEndpoint = "localhost:8082"
	httpPort           = ":8080"
)

func createGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	desc.RegisterOcpTeamApiServer(grpcServer, api.NewOcpTeamApi())

	return grpcServer
}

func createHttpGateway(ctx context.Context) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpTeamApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)

	if err != nil {
		log.Fatal().Msgf("cannot register http handlers %v", err)
	}

	return &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	grpcServer := createGrpcServer()
	httpGateway := createHttpGateway(ctx)

	g.Go(func() error {
		listen, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatal().Msgf("failed to listen: %v", err)
		}

		log.Info().Msgf("grpc server started on port %s", grpcPort)
		return grpcServer.Serve(listen)
	})
	g.Go(func() error {
		log.Info().Msgf("http gateway started on port %s", httpPort)
		return httpGateway.ListenAndServe()
	})

	select {
	case <-interrupt:
		break
	case <- ctx.Done():
		break
	}

	log.Warn().Msg("received shutdown signal")

	cancel()

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer shutdownCtxCancel()

	log.Info().Msg("shutdown http gateway")
	err := httpGateway.Shutdown(shutdownCtx)
	if err != nil {
		log.Debug().Msgf("http gateway shutdown failed %v", err)
	}

	log.Info().Msg("shutdown grpc server")
	grpcServer.GracefulStop()

	if err = g.Wait(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msg(err.Error())
	}
}