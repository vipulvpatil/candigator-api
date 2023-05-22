package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/config"
	"github.com/vipulvpatil/candidate-tracker-go/internal/health"
	"github.com/vipulvpatil/candidate-tracker-go/internal/server"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/tls"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
	"github.com/vipulvpatil/candidate-tracker-go/internal/workers"
	pb "github.com/vipulvpatil/candidate-tracker-go/protos"
	"google.golang.org/grpc"
)

const WORKER_NAMESPACE = "candidate_tracker_go"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg, errs := config.NewConfigFromEnvVars()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		log.Fatal("Unable to load config. Required Env vars are missing")
	}

	logger, deferFunc, err := utilities.InitLogger(utilities.LoggerParams{
		Mode: cfg.LoggerMode,
		SentryParams: struct {
			Dsn         string
			Environment string
		}{
			Dsn:         cfg.SentryDsn,
			Environment: cfg.Environment,
		},
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	if deferFunc != nil {
		defer deferFunc(2 * time.Second)
	}

	db, err := storage.InitDb(cfg, logger)
	if err != nil {
		log.Fatalf("Unable to initialize database: %v", err)
	}

	dbStorage, err := storage.NewDbStorage(
		storage.StorageOptions{
			Db: db,
		},
	)
	if err != nil {
		log.Fatalf("Unable to initialize storage: %v", err)
	}

	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(cfg.RedisUrl)
		},
	}
	_ = workers.NewJobStarter(WORKER_NAMESPACE, redisPool)

	serverDeps := server.ServerDependencies{
		Storage:      dbStorage,
		OpenAiClient: openai.NewClient(openai.OpenAiClientOptions{ApiKey: cfg.OpenAiApiKey}, logger),
		Config:       cfg,
		Logger:       logger,
	}

	s, err := server.NewServer(serverDeps)
	if err != nil {
		log.Fatalf("Unable to initialize new server: %v", err)
	}
	grpcServer := setupGrpcServer(s, cfg, logger)

	workerPooldeps := workers.PoolDependencies{
		RedisPool:    redisPool,
		Namespace:    WORKER_NAMESPACE,
		Storage:      dbStorage,
		OpenAiApiKey: cfg.OpenAiApiKey,
		Logger:       logger,
	}
	workerPool := workers.NewPool(workerPooldeps)
	workerPool.Start()

	var wg sync.WaitGroup
	startGrpcServerAsync("candidate tracker go", &wg, grpcServer, "9000", logger)
	httpHealthServer := startHTTPHealthServer(&wg, logger)

	osTermSig := make(chan os.Signal, 1)
	signal.Notify(osTermSig, syscall.SIGINT, syscall.SIGTERM)

	logger.LogMessageln("Everything started correctly")

	<-osTermSig

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpHealthServer.Shutdown(ctx); err != nil {
		panic(err)
	}
	grpcServer.GracefulStop()
	workerPool.Stop()
	wg.Wait()
	logger.LogMessageln("Stopping Service")
}

func setupGrpcServer(s *server.CandidateTrackerGoService, cfg *config.Config, logger utilities.Logger) *grpc.Server {
	serverOpts := make([]grpc.ServerOption, 0)
	tlsServerOpts := tlsGrpcServerOptions(cfg, logger)
	if tlsServerOpts != nil {
		serverOpts = append(serverOpts, tlsServerOpts)
	}
	serverOpts = append(
		serverOpts,
		grpc.ChainUnaryInterceptor(
			s.RequestingUserInterceptor,
		),
	)
	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterCandidateTrackerGoServer(grpcServer, s)
	return grpcServer
}

func startHTTPHealthServer(wg *sync.WaitGroup, logger utilities.Logger) *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", health.HealthCheckHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.LogMessageln("Starting HTTP Health Check")
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("HTTP Health Check failed with: %v", err)
		}
		logger.LogMessageln("Stopping HTTP Health Check")
	}()
	return srv
}

func startGrpcServerAsync(name string, wg *sync.WaitGroup, grpcServer *grpc.Server, port string, logger utilities.Logger) {
	wg.Add(1)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		logger.LogMessagef("Starting GRPC Server: %s\n", name)
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("GrpcServer %s failed to start: %v", name, err)
		}
		logger.LogMessagef("Stopping GRPC Server: %s\n", name)
		wg.Done()
	}()
}

func tlsGrpcServerOptions(cfg *config.Config, logger utilities.Logger) grpc.ServerOption {
	if cfg.EnableTls {
		tlsCredentials, err := tls.LoadTLSCredentials(cfg)
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}
		logger.LogMessageln("using TLS")
		return grpc.Creds(tlsCredentials)
	}
	return nil
}
