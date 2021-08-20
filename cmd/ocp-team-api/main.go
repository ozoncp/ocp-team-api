package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ozoncp/ocp-team-api/internal/api"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
)

const (
	grpcPort = ":8082"
	grpcServerEndpoint = "localhost:8082"
	httpPort = ":8080"
)

func runJSON(stopSignal <-chan os.Signal, stopStruct chan<- struct{}) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpTeamApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	s := &http.Server{
		Addr: httpPort,
		Handler: mux,
	}

	go func() {
		err = s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-stopSignal
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5 * time.Second)
	defer shutdownCancel()

	if err = s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}

	stopStruct<- struct{}{}
}

func runGrpc(stopStruct <-chan struct{}) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterOcpTeamApiServer(s, api.NewOcpTeamApi())

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stopStruct
	s.GracefulStop()

	return nil
}

func main() {
	stopGateway := make(chan os.Signal, 1)
	stopGrpcServer := make(chan struct{}, 1)

	signal.Notify(stopGateway, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go runJSON(stopGateway, stopGrpcServer)

	if err := runGrpc(stopGrpcServer); err != nil {
		log.Fatal(err)
	}
}