package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/calendar/pb"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/app"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Println(configFile)
		fmt.Println(err)
		return
	}
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage
	switch config.Storage.Type {
	case "postgres":
		storage = sqlstorage.New(config.Database.DSN())
		err = storage.Connect(ctx)
		if err != nil {
			logg.Error(err)
		}
	case "memory":
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)

	httpServer := internalhttp.NewServer(logg, calendar, storage, config.Server.DSN())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.UnaryLoggingInterceptor(logg)),
	)
	pb.RegisterEventsServer(grpcServer, internalgrpc.NewServer(logg, calendar, storage))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		startHTTPServer(ctx, httpServer, logg)
	}()

	go func() {
		defer wg.Done()
		startGRPCServer(ctx, grpcServer, logg)
	}()

	wg.Wait()
	logg.Info("calendar is running...")
}

func startHTTPServer(ctx context.Context, server *internalhttp.Server, logger *logger.Logger) error {
	go func() {
		<-ctx.Done()
		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
		logger.Error("server stopped")
	}()

	logger.Println("http server started: ", server.Addr)
	return server.Start(ctx)
}

func startGRPCServer(ctx context.Context, server *grpc.Server, logger *logger.Logger) error {
	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		logger.Error("server stopped")
	}()

	logger.Println("grpc server started :50051")
	return server.Serve(listener)
}
