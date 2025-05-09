package kafka

import (
	"analytics/internal/domain/models"
	"analytics/internal/services"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer представляет собой потребителя сообщений из Kafka
type Consumer struct {
	reader  *kafka.Reader
	service *services.AnalyticsService
}

// NewConsumer создает нового потребителя Kafka
func NewConsumer(brokers []string, topic, groupID string, service *services.AnalyticsService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  time.Second,
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

// Start запускает обработку сообщений
func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stopping Kafka consumer...")
				if err := c.reader.Close(); err != nil {
					log.Printf("Error closing Kafka reader: %v", err)
				}
				return
			default:
				c.processMessages(ctx)
			}
		}
	}()
}

// Stop для graceful shutdown
func (c *Consumer) Stop() error {
	return c.reader.Close()
}

// processMessages обрабатывает сообщения из Kafka
func (c *Consumer) processMessages(ctx context.Context) {
	// Установим таймаут для чтения сообщений
	readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	msg, err := c.reader.ReadMessage(readCtx)
	if err != nil {
		// Проверим, не вызвана ли ошибка контекстом
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return
		}
		log.Printf("Error reading message: %v", err)
		return
	}

	// Десериализуем сообщение в событие аналитики
	var event models.AnalyticsEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Обработаем событие
	if err = c.service.ProcessEvent(ctx, event); err != nil {
		log.Printf("Error processing event: %v", err)
		return
	}

	slog.Debug("Processed message",
		slog.String("topic", msg.Topic),
		slog.Int("partition", msg.Partition),
		slog.Int64("offset", msg.Offset),
	)
}
