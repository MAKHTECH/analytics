package main

import (
	api "analytics/internal/api/gen"
	"analytics/internal/api/handlers"
	"analytics/internal/config"
	"analytics/internal/infrastructure/kafka"
	"analytics/internal/infrastructure/storage/clickhouse"
	"analytics/internal/infrastructure/storage/metrics"
	"analytics/internal/services"
	"analytics/pkg/logging"
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const kafkaUI string = "localhost:8001"

func main() {
	cfg := config.MustLoad("local.json")
	logging.SetupLogger()
	slog.Info("starting application", slog.Any("config", cfg), slog.String("kafka-ui", kafkaUI))

	// Инициализация метрик
	metricsClient := metrics.NewPrometheusClient()

	// Инициализация хранилища clickhouse
	chClient, err := clickhouse.NewClient(
		cfg.Database.Database, cfg.Database.Username,
		cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port,
	)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	// Создание репозитория
	repository := clickhouse.NewClickHouseRepository(chClient)

	// Создание сервиса аналитики
	analyticsService := services.NewAnalyticsService(repository, metricsClient)

	// Инициализация потребителя Kafka
	consumer := kafka.NewConsumer(
		[]string{"kafka:9092"},
		"topic",
		analyticsService,
	)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск потребителя Kafka
	consumer.Start(ctx)

	// Инициализация HTTP сервера для API и метрик
	serverApi := handlers.NewServerAPI(analyticsService)
	httpServer, err := api.NewServer(serverApi)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", httpServer)
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         cfg.App.Address + ":" + strconv.Itoa(cfg.App.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.App.Port)
		if err = server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down service...")
	cancel()

	// Graceful shutdown HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}

	log.Println("Service stopped")
}
