package resiliency

import (
	"context"
	"errors"
	"testing"
)

func TestExecuteWithRetryExhaustion(t *testing.T) {
	attempts := 0
	op := func() (string, error) {
		attempts++
		return "", errors.New("ThrottlingException: Rate exceeded")
	}

	ctx := context.Background()
	_, err := ExecuteWithRetry(ctx, op)

	if err == nil {
		t.Fatalf("expected error after retry exhaustion, got nil")
	}
	if attempts != 5 {
		t.Errorf("expected 5 retry attempts, got %d", attempts)
	}
}

func TestExecuteWithRetryContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	attempts := 0
	op := func() (string, error) {
		attempts++
		if attempts == 1 {
			cancel() // Cancel context after first attempt
		}
		return "", errors.New("429 Too Many Requests")
	}

	_, err := ExecuteWithRetry(ctx, op)
	if err == nil {
		t.Fatalf("expected context error, got nil")
	}
}

func TestExecuteWithRetryNonRetryableFastFail(t *testing.T) {
	attempts := 0
	op := func() (string, error) {
		attempts++
		return "", errors.New("400 Bad Request: Invalid Parameter")
	}

	ctx := context.Background()
	_, err := ExecuteWithRetry(ctx, op)

	if err == nil {
		t.Fatalf("expected error for 400 Bad Request, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected fast-fail (1 attempt), got %d attempts", attempts)
	}
}
