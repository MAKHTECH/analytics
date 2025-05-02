package models

import (
	"time"
)

// AnalyticsEvent представляет собой событие аналитики
type AnalyticsEvent struct {
	ID         string            `json:"id,omitempty"`
	EventType  string            `json:"eventType"`
	UserID     string            `json:"userId"`
	Properties map[string]string `json:"properties,omitempty"`
	Duration   int64             `json:"duration"`
	Timestamp  time.Time         `json:"timestamp"`
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

//{
//"id": "event_12345",
//"eventType": "user_signup",
//"userId": "user_98765",
//"properties": {
//"source": "landing_page",
//"campaign": "spring_sale",
//"referrer": "google_ads"
//},
//"duration": 248,
//"timestamp": "2025-05-02T12:34:56Z"
//}
