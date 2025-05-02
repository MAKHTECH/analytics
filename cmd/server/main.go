package main

import (
	api "analytics/internal/api"
	"analytics/internal/config"
	"analytics/internal/infrastructure/kafka"
	"analytics/internal/infrastructure/storage/clickhouse"
	"analytics/internal/infrastructure/storage/metrics"
	"analytics/internal/services"
	"analytics/pkg/logging"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// get event types

const kafkaUI string = "localhost:8001"

func main() {
	cfg := config.MustLoad("local.json")
	logging.SetupLogger()
	slog.Info("starting application", slog.Any("config", cfg), slog.String("kafka-ui", kafkaUI))

	metricsClient := metrics.NewPrometheusClient()
	chClient, err := clickhouse.NewClient(
		cfg.Database.Database, cfg.Database.Username,
		cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port,
	)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	repository := clickhouse.NewClickHouseRepository(chClient)
	analyticsService := services.NewAnalyticsService(repository, metricsClient)
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID,
		analyticsService,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx)
	server := api.InitServer(cfg.App.Address, cfg.App.Port, analyticsService)

	// Запуск graceful shutdown
	shutdown(cancel, server, chClient, consumer)
}

func shutdown(cancel context.CancelFunc, server *http.Server, chClient *clickhouse.Client, consumer *kafka.Consumer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down service...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := chClient.Close(); err != nil {
		log.Printf("Failed to close ClickHouse client: %v", err)
	}

	if err := consumer.Stop(); err != nil {
		log.Printf("Failed to close kafka client: %v", err)
	}

	log.Println("Service stopped")
}
