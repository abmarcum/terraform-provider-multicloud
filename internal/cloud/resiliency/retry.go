package resiliency

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"
)

// IsRetryableError returns true if the error indicates a transient cloud API error
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())

	retryableKeywords := []string{
		"throttlingexception",
		"rate exceeded",
		"too many requests",
		"429",
		"503 service unavailable",
		"502 bad gateway",
		"500 internal server error",
		"requestlimitexceeded",
		"resourceexhausted",
		"serviceunavailable",
		"connection reset",
		"timeout",
	}

	for _, kw := range retryableKeywords {
		if strings.Contains(msg, kw) {
			return true
		}
	}

	return false
}

// ExecuteWithRetry runs an operation with exponential backoff, jitter, context timeouts, and sensitive error redaction
func ExecuteWithRetry[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	var zero T
	maxRetries := 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 3 * time.Second

	// Enforce 10-minute maximum context timeout bound if none provided
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctxWithTimeout.Done():
			return zero, ctxWithTimeout.Err()
		default:
		}

		result, err := operation()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !IsRetryableError(err) {
			return zero, RedactSensitiveLogInfo(err)
		}

		// Calculate exponential backoff with full jitter (non-cryptographic PRNG acceptable for backoff)
		backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
		if backoff > float64(maxDelay) {
			backoff = float64(maxDelay)
		}
		/* #nosec G404 */
		jitter := rand.Float64() * backoff
		sleepDuration := time.Duration(jitter)

		select {
		case <-ctxWithTimeout.Done():
			return zero, ctxWithTimeout.Err()
		case <-time.After(sleepDuration):
		}
	}

	return zero, RedactSensitiveLogInfo(lastErr)
}

// RedactSensitiveLogInfo strips authorization headers and tokens from log messages
func RedactSensitiveLogInfo(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()

	// Redact bearer tokens or access keys if present
	if strings.Contains(msg, "Bearer ") {
		msg = strings.Split(msg, "Bearer ")[0] + "Bearer [REDACTED]"
	}
	if strings.Contains(msg, "AWS4-HMAC-SHA256") {
		msg = strings.Split(msg, "AWS4-HMAC-SHA256")[0] + "AWS4-HMAC-SHA256 [REDACTED]"
	}

	return errors.New(msg)
}
