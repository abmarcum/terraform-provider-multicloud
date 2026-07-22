package telemetry

import (
	"encoding/json"
	"fmt"
	"time"
)

// TelemetryEvent models structured OpenTelemetry and audit logging events
type TelemetryEvent struct {
	Timestamp   string                 `json:"timestamp"`
	EventType   string                 `json:"event_type"` // "AUDIT", "COST", "RETRY", "POLICY"
	Provider    string                 `json:"provider"`
	Resource    string                 `json:"resource"`
	DurationMs  int64                  `json:"duration_ms"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TelemetryExporter outputs structured OpenTelemetry JSON logs
type TelemetryExporter struct{}

// NewTelemetryExporter returns a new TelemetryExporter instance
func NewTelemetryExporter() *TelemetryExporter {
	return &TelemetryExporter{}
}

// RecordEvent formats and emits a structured telemetry event
func (t *TelemetryExporter) RecordEvent(eventType string, provider string, resource string, duration time.Duration, meta map[string]interface{}) (string, error) {
	event := TelemetryEvent{
		Timestamp:  time.Now().Format(time.RFC3339),
		EventType:  eventType,
		Provider:   provider,
		Resource:   resource,
		DurationMs: duration.Milliseconds(),
		Metadata:   meta,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("[TelemetryExporter] %s", string(bytes))
	return out, nil
}
