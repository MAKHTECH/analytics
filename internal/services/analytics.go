// internal/domain/services/analytics.go
package services

import (
	"analytics/internal/domain/models"
	"context"
	"time"
)

// EventRepository интерфейс для работы с хранилищем событий
type EventRepository interface {
	StoreEvent(ctx context.Context, event models.AnalyticsEvent) error
	GetEventsByType(ctx context.Context, eventType string, from, to time.Time) ([]models.AnalyticsEvent, error)
	GetEventCounts(ctx context.Context, from, to time.Time) (map[string]int, error)
	GetUniqueUsers(ctx context.Context, from, to time.Time) (int, error)
}

// MetricsProvider интерфейс для работы с метриками
type MetricsProvider interface {
	RecordEvent(eventType string)
	RecordProcessingTime(duration time.Duration)
}

// AnalyticsService сервис для работы с аналитикой
type AnalyticsService struct {
	repo    EventRepository
	metrics MetricsProvider
}

// NewAnalyticsService создает новый экземпляр сервиса аналитики
func NewAnalyticsService(repo EventRepository, metrics MetricsProvider) *AnalyticsService {
	return &AnalyticsService{
		repo:    repo,
		metrics: metrics,
	}
}

// ProcessEvent обрабатывает событие аналитики
func (s *AnalyticsService) ProcessEvent(ctx context.Context, event models.AnalyticsEvent) error {
	start := time.Now()

	// Если timestamp не установлен, используем текущее время
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Сохраняем событие в хранилище
	if err := s.repo.StoreEvent(ctx, event); err != nil {
		return err
	}

	// Обновляем метрики
	s.metrics.RecordEvent(event.EventType)
	s.metrics.RecordProcessingTime(time.Since(start))

	return nil
}

// GetMetrics получает агрегированные метрики за период
func (s *AnalyticsService) GetMetrics(ctx context.Context, from, to time.Time, eventType string) (map[string]interface{}, error) {
	var counts map[string]int
	var err error

	// Если указан конкретный тип события, получаем счетчики только для него
	if eventType != "" {
		events, err := s.repo.GetEventsByType(ctx, eventType, from, to)
		if err != nil {
			return nil, err
		}
		counts = map[string]int{eventType: len(events)}
	} else {
		// Иначе получаем счетчики для всех типов событий
		counts, err = s.repo.GetEventCounts(ctx, from, to)
		if err != nil {
			return nil, err
		}
	}

	// Получаем количество уникальных пользователей
	users, err := s.repo.GetUniqueUsers(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Формируем результат
	result := map[string]interface{}{
		"counts": counts,
		"users":  users,
		"period": map[string]time.Time{
			"from": from,
			"to":   to,
		},
	}

	return result, nil
}
