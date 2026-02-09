package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/vtapaskar/brahma/internal/config"
	"github.com/vtapaskar/brahma/internal/metrics"
	"github.com/vtapaskar/brahma/internal/server"
	"github.com/vtapaskar/brahma/internal/splunk"
	"github.com/vtapaskar/brahma/internal/storage"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s3Client, err := storage.NewS3Client(ctx, cfg.S3)
	if err != nil {
		logger.Fatal("Failed to initialize S3 client", zap.Error(err))
	}

	splunkClient := splunk.NewClient(cfg.Splunk, logger)

	metricsCollector := metrics.NewCollector(cfg.Metrics, splunkClient, s3Client, logger)

	srv := server.New(cfg.Server, metricsCollector, logger)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	logger.Info("Brahma service started",
		zap.String("address", cfg.Server.Address),
		zap.Int("port", cfg.Server.Port),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
	srv.Stop(ctx)
}
