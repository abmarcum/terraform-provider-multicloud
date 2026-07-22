package telemetry

import (
	"strings"
	"testing"
	"time"
)

func TestTelemetryExporterRecordEvent(t *testing.T) {
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
}
