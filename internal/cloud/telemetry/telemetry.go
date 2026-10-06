package telemetry

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func isSafeOTLPEndpoint(rawURL string) bool {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "169.254.169.254" || host == "fd00:ec2::254" || host == "metadata.google.internal" || strings.HasPrefix(host, "169.254.") {
		return false
	}
	return u.Scheme == "https" || (u.Scheme == "http" && (host == "127.0.0.1" || host == "localhost" || host == "::1"))
}

// TelemetryEvent models structured OpenTelemetry and audit logging events
type TelemetryEvent struct {
	Timestamp  string                 `json:"timestamp"`
	EventType  string                 `json:"event_type"` // "AUDIT", "COST", "RETRY", "POLICY", "PROVISION", "READ", "UPDATE", "DELETE"
	Provider   string                 `json:"provider"`
	Resource   string                 `json:"resource"`
	DurationMs int64                  `json:"duration_ms"`
	Metadata   map[string]interface{} `json:"metadata"`
}

const maxBufferedEvents = 512

var otlpHTTPClient = func() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}
}()

// DefaultExporter is the shared telemetry exporter used across cloud adapter lifecycles
var DefaultExporter = NewTelemetryExporter()

// TelemetryExporter outputs structured OpenTelemetry JSON logs and buffers events for inspection
type TelemetryExporter struct {
	mu     sync.RWMutex
	events []TelemetryEvent
}

// NewTelemetryExporter returns a new TelemetryExporter instance
func NewTelemetryExporter() *TelemetryExporter {
	return &TelemetryExporter{
		events: make([]TelemetryEvent, 0, 64),
	}
}

// RecordEvent formats, buffers, and optionally exports a structured telemetry event to an OTLP HTTP collector
func (t *TelemetryExporter) RecordEvent(eventType string, provider string, resource string, duration time.Duration, meta map[string]interface{}) (string, error) {
	event := TelemetryEvent{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		EventType:  eventType,
		Provider:   provider,
		Resource:   resource,
		DurationMs: duration.Milliseconds(),
		Metadata:   meta,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	t.mu.Lock()
	if len(t.events) >= maxBufferedEvents {
		t.events = append(t.events[1:], event)
	} else {
		t.events = append(t.events, event)
	}
	t.mu.Unlock()

	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" && os.Getenv("MULTICLOUD_MOCK_MODE") != "true" && isSafeOTLPEndpoint(endpoint) {
		/* #nosec G107 G704 */
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(payload))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			/* #nosec G107 G704 */
			if resp, err := otlpHTTPClient.Do(req); err == nil {
				_ = resp.Body.Close()
			}
		}
	}

	out := fmt.Sprintf("[TelemetryExporter] %s", string(payload))
	return out, nil
}

// Events returns a snapshot of recorded telemetry events
func (t *TelemetryExporter) Events() []TelemetryEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	copied := make([]TelemetryEvent, len(t.events))
	copy(copied, t.events)
	return copied
}
