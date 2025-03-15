package models

import (
	"time"
)

// AnalyticsEvent представляет собой событие аналитики
type AnalyticsEvent struct {
	ID         string            `json:"id,omitempty"`
	Timestamp  time.Time         `json:"timestamp"`
	EventType  string            `json:"eventType"`
	UserID     string            `json:"userId"`
	Properties map[string]string `json:"properties,omitempty"`
}
