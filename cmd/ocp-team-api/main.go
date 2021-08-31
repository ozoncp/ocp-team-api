package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-team-api/internal/api"
	"github.com/ozoncp/ocp-team-api/internal/config"
	"github.com/ozoncp/ocp-team-api/internal/kafka"
	"github.com/ozoncp/ocp-team-api/internal/metrics"
	"github.com/ozoncp/ocp-team-api/internal/repo"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"
	"io"
	"sync/atomic"

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

func createGrpcServer(db *sqlx.DB, producer kafka.Producer) *grpc.Server {
	grpcServer := grpc.NewServer()
	desc.RegisterOcpTeamApiServer(grpcServer, api.NewOcpTeamApi(repo.NewRepo(db), producer))

	return grpcServer
}

func createHttpGateway(ctx context.Context) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpTeamApiHandlerFromEndpoint(
		ctx,
		mux,
		config.GetInstance().Server.Host+config.GetInstance().Server.GrpcPort,
		opts,
	)

	if err != nil {
		log.Fatal().Msgf("cannot register http handlers %v", err)
	}

	return &http.Server{
		Addr:    config.GetInstance().Server.HttpPort,
		Handler: mux,
	}
}

func createStatusServer() *http.Server {
	isReady := &atomic.Value{}
	isReady.Store(false)

	go func() {
		log.Debug().Msg("Ready probe is negative by default...")
		time.Sleep(time.Duration(config.GetInstance().Server.StartupTime) * time.Second)
		isReady.Store(true)
		log.Debug().Msg("Ready probe is positive.")
	}()

	mux := http.DefaultServeMux

	mux.HandleFunc(config.GetInstance().Status.HealthHandler, health)
	mux.HandleFunc(config.GetInstance().Status.ReadyHandler, ready(isReady))

	return &http.Server{
		Addr:    config.GetInstance().Status.Port,
		Handler: mux,
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ready(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func createMetricsHttpHandler() *http.Server {
	mux := http.NewServeMux()
	mux.Handle(config.GetInstance().Metrics.Handler, promhttp.Handler())

	return &http.Server{
		Addr:    config.GetInstance().Metrics.Port,
		Handler: mux,
	}
}

func createTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: config.GetInstance().Jaeger.ServiceName,
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
	db, err := sqlx.Connect("pgx", config.GetInstance().Database.DSN)

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
	log.Info().Msg("connection with DB established")
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
	statusServer := createStatusServer()

	g.Go(func() error {
		listen, err := net.Listen("tcp", config.GetInstance().Server.GrpcPort)
		if err != nil {
			log.Fatal().Msgf("failed to listen: %v", err)
		}

		log.Info().Msgf("grpc server started on port %s", config.GetInstance().Server.GrpcPort)
		return grpcServer.Serve(listen)
	})
	g.Go(func() error {
		log.Info().Msgf("http gateway started on port %s", config.GetInstance().Server.HttpPort)
		return httpGateway.ListenAndServe()
	})
	g.Go(func() error {
		log.Info().Msgf("metrics http handler started on port %s", config.GetInstance().Metrics.Port)
		metrics.Register()
		return metricsHttpHandler.ListenAndServe()
	})
	g.Go(func() error {
		log.Info().Msgf("status server started on port %s", config.GetInstance().Status.Port)
		return statusServer.ListenAndServe()
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Warn().Msg("received shutdown signal")

	cancel()

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.GetInstance().Server.ShutdownTime)*time.Second,
	)
	defer shutdownCtxCancel()

	log.Info().Msg("shutdown status server")
	err = statusServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Debug().Msgf("http status server shutdown failed %v", err)
	}

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
