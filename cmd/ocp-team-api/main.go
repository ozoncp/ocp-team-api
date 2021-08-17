package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/ozoncp/ocp-team-api/internal/api"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
)

const (
	grpcPort = ":8082"
	grpcServerEndpoint = "localhost:8082"
	httpPort = ":8080"
)

func runJSON() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpTeamApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(httpPort, mux)
	if err != nil {
		panic(err)
	}
}

func runGrpc() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterOcpTeamApiServer(s, api.NewOcpTeamApi())

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}


func main() {
	go runJSON()

	if err := runGrpc(); err != nil {
		log.Fatal(err)
	}
}
