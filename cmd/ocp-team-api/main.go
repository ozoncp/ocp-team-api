package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-team-api/internal/api"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
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

var (
	grpcServer  *grpc.Server
	httpGateway *http.Server
)

func runGrpcServer() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer = grpc.NewServer()
	desc.RegisterOcpTeamApiServer(grpcServer, api.NewOcpTeamApi())

	return grpcServer.Serve(listen)
}

func runHttpGateway(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpTeamApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	httpGateway = &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	if err = httpGateway.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(runGrpcServer)
	g.Go(func() error { return runHttpGateway(ctx) })

	select {
	case <-interrupt:
		break
	case <- ctx.Done():
		break
	}

	log.Println("received shutdown signal")

	cancel()

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer shutdownCtxCancel()

	if httpGateway != nil {
		_ = httpGateway.Shutdown(shutdownCtx)
		log.Println("shutdown http gateway")
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
		log.Println("shutdown grpc server")
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}