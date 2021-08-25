package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-team-api/internal/api"
	"github.com/ozoncp/ocp-team-api/internal/kafka"
	"github.com/ozoncp/ocp-team-api/internal/metrics"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
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
	prometheusPort     = ":9100"

	dsn = "postgres://root:root@localhost:5432/postgres?sslmode=disable"

	shutdownTimeout = 5 * time.Second
)

func createGrpcServer(db *sqlx.DB, producer kafka.Producer) *grpc.Server {
	grpcServer := grpc.NewServer()
	desc.RegisterOcpTeamApiServer(grpcServer, api.NewOcpTeamApi(repo.NewRepo(db), producer))

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

func createMetricsHttpHandler() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    prometheusPort,
		Handler: mux,
	}
}

func createTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: "ocp_team_api",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jaegerMetrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	opentracing.SetGlobalTracer(tracer)

	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}

func db() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	db, err := db()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()

	_, closer, err := createTracer()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer closer.Close()

	kafkaProducer, err := kafka.NewProducer()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	grpcServer := createGrpcServer(db, kafkaProducer)
	httpGateway := createHttpGateway(ctx)
	metricsHttpHandler := createMetricsHttpHandler()

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
	g.Go(func() error {
		log.Info().Msgf("metrics http handler started on port %s", prometheusPort)
		metrics.Register()
		return metricsHttpHandler.ListenAndServe()
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Warn().Msg("received shutdown signal")

	cancel()

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCtxCancel()

	log.Info().Msg("shutdown http gateway")
	err = httpGateway.Shutdown(shutdownCtx)
	if err != nil {
		log.Debug().Msgf("http gateway shutdown failed %v", err)
	}

	log.Info().Msg("shutdown metrics http handler")
	err = metricsHttpHandler.Shutdown(shutdownCtx)
	if err != nil {
		log.Debug().Msgf("metric http handler shutdown failed %v", err)
	}

	log.Info().Msg("shutdown grpc server")
	grpcServer.GracefulStop()

	if err = g.Wait(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msg(err.Error())
	}
}
