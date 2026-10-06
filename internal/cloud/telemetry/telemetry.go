package telemetry

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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
	transport.MaxIdleConns = 32
	transport.MaxIdleConnsPerHost = 16
	transport.IdleConnTimeout = 90 * time.Second
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

// TelemetryExporter outputs structured OpenTelemetry JSON logs and buffers events in a circular ring buffer
type TelemetryExporter struct {
	mu     sync.RWMutex
	events []TelemetryEvent
	head   int
	count  int
}

// NewTelemetryExporter returns a new TelemetryExporter instance
func NewTelemetryExporter() *TelemetryExporter {
	return &TelemetryExporter{
		events: make([]TelemetryEvent, maxBufferedEvents),
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
	if t.count < maxBufferedEvents {
		idx := (t.head + t.count) % maxBufferedEvents
		t.events[idx] = event
		t.count++
	} else {
		t.events[t.head] = event
		t.head = (t.head + 1) % maxBufferedEvents
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

	return "[TelemetryExporter] " + string(payload), nil
}

// Events returns a chronological snapshot of recorded telemetry events
func (t *TelemetryExporter) Events() []TelemetryEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	copied := make([]TelemetryEvent, t.count)
	for i := 0; i < t.count; i++ {
		copied[i] = t.events[(t.head+i)%maxBufferedEvents]
	}
	return copied
}
