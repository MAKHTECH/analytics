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
	"errors"
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

func main() {
	cfg := config.MustLoad("local.json")
	logging.SetupLogger()
	slog.Info("starting application", slog.Any("config", cfg))

	// Инициализация метрик
	metricsClient := metrics.NewPrometheusClient()

	// Инициализация хранилища clickhouse
	chClient, err := clickhouse.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	// Создание репозитория
	repository := clickhouse.NewClickHouseRepository(chClient)

	// Создание сервиса аналитики
	analyticsService := services.NewAnalyticsService(repository, metricsClient)

	// Инициализация потребителя Kafka
	consumer := kafka.NewConsumer(
		[]string{"brokers"},
		"topic",
		analyticsService,
	)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск потребителя Kafka
	consumer.Start(ctx)

	// Инициализация HTTP сервера для API и метрик
	handlers := handlers.NewServerAPI(analyticsService)
	httpServer, err := api.NewServer(handlers)
	server := &http.Server{
		Addr:         cfg.App.Address + ":" + strconv.Itoa(cfg.App.Port),
		Handler:      httpServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Добавление обработчика для метрик Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
