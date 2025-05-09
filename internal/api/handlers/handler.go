package handlers

import (
	api "analytics/internal/api/gen"
	"analytics/internal/domain/models"
	"context"
	"time"
)

type AnalyticService interface {
	ProcessEvent(ctx context.Context, event models.AnalyticsEvent) error
	GetMetrics(ctx context.Context, from, to time.Time, eventType string) (map[string]interface{}, error)
}

type ServerAPI struct {
	Service AnalyticService
}

func NewServerAPI(service AnalyticService) *ServerAPI {
	return &ServerAPI{
		Service: service,
	}
}

func (s *ServerAPI) GetMetrics(ctx context.Context, params api.GetMetricsParams) (api.GetMetricsRes, error) {
	panic("implement me")
}
