package clickhouse

import (
	"analytics/internal/domain/models"
	"context"
	"time"
)

// ClickHouseRepository реализация репозитория для ClickHouse
type ClickHouseRepository struct {
	client *Client
}

// NewClickHouseRepository создает новый репозиторий для ClickHouse
func NewClickHouseRepository(client *Client) *ClickHouseRepository {
	return &ClickHouseRepository{
		client: client,
	}
}

// StoreEvent сохраняет событие в ClickHouse
func (r *ClickHouseRepository) StoreEvent(ctx context.Context, event models.AnalyticsEvent) error {
	panic("implement me")
}

// GetEventsByType получает события по типу за период
func (r *ClickHouseRepository) GetEventsByType(ctx context.Context, eventType string, from, to time.Time) ([]models.AnalyticsEvent, error) {
	panic("implement me")
}

// GetEventCounts получает количество событий по типам за период
func (r *ClickHouseRepository) GetEventCounts(ctx context.Context, from, to time.Time) (map[string]int, error) {
	panic("implement me")
}

// GetUniqueUsers получает количество уникальных пользователей за период
func (r *ClickHouseRepository) GetUniqueUsers(ctx context.Context, from, to time.Time) (int, error) {
	panic("implement me")
}
