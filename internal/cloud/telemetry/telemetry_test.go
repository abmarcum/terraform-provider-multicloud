package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestTelemetryExporterRecordEventAndBuffer(t *testing.T) {
	exporter := NewTelemetryExporter()
	meta := map[string]interface{}{
		"rule_id": "CIS-STORAGE-ENCRYPTION-1.1",
		"passed":  true,
	}

	eventStr, err := exporter.RecordEvent("AUDIT", "aws", "prod-bucket", 45*time.Millisecond, meta)
	if err != nil {
		t.Fatalf("unexpected error emitting telemetry event: %v", err)
	}

	if !strings.Contains(eventStr, "AUDIT") || !strings.Contains(eventStr, "prod-bucket") {
		t.Errorf("telemetry string missing expected event fields: %s", eventStr)
	}

	events := exporter.Events()
	if len(events) != 1 || events[0].Resource != "prod-bucket" {
		t.Errorf("expected 1 buffered telemetry event for prod-bucket, got %+v", events)
	}

	for i := 0; i < maxBufferedEvents+10; i++ {
		_, _ = exporter.RecordEvent("READ", "aws", "ring-res", time.Millisecond, nil)
	}
	if got := len(exporter.Events()); got != maxBufferedEvents {
		t.Errorf("expected ring buffer to cap at %d events, got %d", maxBufferedEvents, got)
	}
}

func TestTelemetryExporterOTLPCollectorExport(t *testing.T) {
	var received int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&received, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("MULTICLOUD_MOCK_MODE", "false")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", srv.URL)

	exporter := NewTelemetryExporter()
	_, err := exporter.RecordEvent("PROVISION", "gcp", "test-vm", 12*time.Millisecond, nil)
	if err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}
	if atomic.LoadInt32(&received) != 1 {
		t.Errorf("expected OTLP endpoint to receive 1 POST request, got %d", received)
	}
}
