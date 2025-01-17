package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/joho/godotenv"
	"github.com/peter-pashchenko/grpcExchangeService/config"
	grpcHandler "github.com/peter-pashchenko/grpcExchangeService/internal/application/grpc"
	exchangeRate_v1 "github.com/peter-pashchenko/grpcExchangeService/internal/generated/grpc/exchangeRate"
	apiGarantex "github.com/peter-pashchenko/grpcExchangeService/internal/infrastructure/apiGarantex"
	"github.com/peter-pashchenko/grpcExchangeService/internal/infrastructure/router"
	_ "github.com/peter-pashchenko/grpcExchangeService/internal/models"
	exchangeRateRepo "github.com/peter-pashchenko/grpcExchangeService/internal/modules/repository/exchangeRate"
	exchangeRateService "github.com/peter-pashchenko/grpcExchangeService/internal/modules/services/exchangeRate"
	"github.com/peter-pashchenko/grpcExchangeService/pkg/logger"
	"github.com/peter-pashchenko/grpcExchangeService/pkg/psql"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	exitStatusOk    = 0
	exitStatusError = 1
	idMarket        = "usdtrub"
)

func main() {
	_ = godotenv.Load( /*"../../.env"*/)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config init error %s", err)
	}

	os.Exit(run(cfg))
}

func run(cfg *config.Config) (exitCode int) {
	log := logger.New(cfg.Log.Level)
	defer log.Sync()

	defer func() {
		if paniceErr := recover(); paniceErr != nil {
			log.Error(
				"recover panic",
				zap.Any("error", paniceErr),
			)
			exitCode = exitStatusError
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM)
	defer stop()

	db, err := psql.Connect(
		ctx,
		log,
		psql.WithHost(cfg.PG.Host),
		psql.WithPort(cfg.PG.Port),
		psql.WithUser(cfg.PG.User),
		psql.WithPass(cfg.PG.Pass),
		psql.WithDatabase(cfg.PG.Database),
		psql.WithMigrations("./db/migrations"),
	)

	defer func() {
		if db != nil {
			err = db.Close()
			log.Debug("closing db...")
			if err != nil {
				log.Error(
					"error closing db",
					zap.Error(err))
			}
		}

	}()

	if err != nil {
		log.Error(
			"db init error",
			zap.Any("error", err),
		)
		return exitStatusError
	}

	r := router.New(
		log,
		router.WithPrometheus(),
		router.WithJaeger(cfg.HTTP.Host))

	repo := exchangeRateRepo.New(db, log)
	apiService := apiGarantex.New(idMarket, log)
	service := exchangeRateService.New(repo, apiService, log)
	serviceHandler := grpcHandler.New(
		grpcHandler.WithExchangeRateService(service),
		grpcHandler.WithLogger(log),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpcHandler.MetricsInterceptor,
				grpcHandler.TracingInterceptor),
		),
	)

	exchangeRate_v1.RegisterExchangeRateServiceServer(grpcServer, serviceHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		log.Error(
			"grpc failed to listen",
			zap.Error(err),
		)
		return exitStatusError
	}

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Host,
		Handler:      r.GetRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errChan := make(chan error)

	go func() {
		log.Info("starting grpc server")

		if err = grpcServer.Serve(lis); err != nil {
			log.Error(
				"grpc failed to serve",
				zap.Error(err))
			errChan <- err
		}
	}()

	go func() {
		log.Info("starting http server")

		if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(
				"http server failed to start",
				zap.Error(err))
		}
		errChan <- err
	}()

	defer func() {
		grpcServer.GracefulStop()

		err = r.TPShutdown(context.Background())
		if err != nil {
			log.Error(
				"trace provider shutdown error",
				zap.Error(err))
		}
		err = httpServer.Shutdown(context.Background())

		if err != nil {
			log.Error(
				"http server shutdown error",
				zap.Error(err))
		}
	}()

	select {
	case err = <-errChan:
		log.Error(
			"one of the servers error,shutting down",
			zap.Error(err),
		)
	case <-ctx.Done():
		log.Info("gracefully shutting down")
	}

	return exitStatusOk
}
