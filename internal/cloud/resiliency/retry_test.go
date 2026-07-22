package resiliency

import (
	"context"
	"errors"
	"testing"
)

func TestIsRetryableError(t *testing.T) {
	if !IsRetryableError(errors.New("ThrottlingException: Rate exceeded")) {
		t.Errorf("expected ThrottlingException to be retryable")
	}
	if !IsRetryableError(errors.New("HTTP 429 Too Many Requests")) {
		t.Errorf("expected 429 to be retryable")
	}
	if IsRetryableError(errors.New("InvalidParameterValue: Bucket name invalid")) {
		t.Errorf("expected InvalidParameterValue to NOT be retryable")
	}
}

func TestExecuteWithRetrySuccess(t *testing.T) {
	attempts := 0
	op := func() (string, error) {
		attempts++
		if attempts < 2 {
			return "", errors.New("ThrottlingException: Rate exceeded")
		}
		return "success-bucket", nil
	}

	result, err := ExecuteWithRetry(context.Background(), op)
	if err != nil {
		t.Fatalf("expected operation to succeed after retry, got error: %v", err)
	}
	if result != "success-bucket" {
		t.Errorf("expected 'success-bucket', got '%s'", result)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}
