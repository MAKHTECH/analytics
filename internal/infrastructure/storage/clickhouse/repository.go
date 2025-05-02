package clickhouse

import (
	"analytics/internal/domain/models"
	"context"
	"time"
)

// Repository реализация репозитория для ClickHouse
type Repository struct {
	client *Client
}

// NewClickHouseRepository создает новый репозиторий для ClickHouse
func NewClickHouseRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

// StoreEvent сохраняет событие в ClickHouse
func (r *Repository) StoreEvent(ctx context.Context, event models.AnalyticsEvent) error {
	panic("implement me")
}

// GetEventsByType получает события по типу за период
func (r *Repository) GetEventsByType(ctx context.Context, eventType string, from, to time.Time) ([]models.AnalyticsEvent, error) {
	panic("implement me")
}

// GetEventCounts получает количество событий по типам за период
func (r *Repository) GetEventCounts(ctx context.Context, from, to time.Time) (map[string]int, error) {
	panic("implement me")
}

// GetUniqueUsers получает количество уникальных пользователей за период
func (r *Repository) GetUniqueUsers(ctx context.Context, from, to time.Time) (int, error) {
	panic("implement me")
}
