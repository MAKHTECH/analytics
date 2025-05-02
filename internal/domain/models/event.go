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

//// test data
//
//{
//		"id": "1",
//		"timestamp": "2025-05-02T10:00:00Z",
//		"eventType": "create",
//		"userId": "1",
//		"properties": {
//			"test": "test"
//		}
//}
