package resiliency

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	bearerTokenRegex = regexp.MustCompile(`(?i)(Bearer\s+)[^\s"',;]+`)
	awsSigV4Regex    = regexp.MustCompile(`(AWS4-HMAC-SHA256)\s+[^\r\n"']+`)
	secretParamRegex = regexp.MustCompile(`(?i)((?:client_secret|access_token|refresh_token|secret_key|aws_secret_access_key|api_key|password|private_key|x-amz-security-token)["']?\s*[=:]\s*["']?)[^\s"',;&}]+`)

	retryableKeywords = [...]string{
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
)

// IsRetryableError returns true if the error indicates a transient cloud API error
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())

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
	const (
		maxRetries = 5
		baseDelay  = 100 * time.Millisecond
		maxDelay   = 3 * time.Second
	)

	// Enforce 10-minute maximum context timeout bound if no shorter deadline is already present
	ctxWithTimeout := ctx
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > 10*time.Minute {
		var cancel context.CancelFunc
		ctxWithTimeout, cancel = context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
	}

	var (
		lastErr error
		timer   *time.Timer
	)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

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
		sleepDuration := time.Duration(rand.Float64() * backoff)
		if timer == nil {
			timer = time.NewTimer(sleepDuration)
		} else {
			timer.Reset(sleepDuration)
		}

		select {
		case <-ctxWithTimeout.Done():
			return zero, ctxWithTimeout.Err()
		case <-timer.C:
		}
	}

	return zero, RedactSensitiveLogInfo(lastErr)
}

// RedactSensitiveString strips authorization headers, signatures, and secret parameters from raw strings
func RedactSensitiveString(msg string) string {
	if msg == "" {
		return ""
	}
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "bearer") {
		msg = bearerTokenRegex.ReplaceAllString(msg, "${1}[REDACTED]")
	}
	if strings.Contains(msg, "AWS4-HMAC-SHA256") {
		msg = awsSigV4Regex.ReplaceAllString(msg, "${1} [REDACTED]")
	}
	if strings.Contains(lower, "secret") || strings.Contains(lower, "token") || strings.Contains(lower, "key") || strings.Contains(lower, "password") {
		msg = secretParamRegex.ReplaceAllString(msg, "${1}[REDACTED]")
	}
	return msg
}

// RedactSensitiveLogInfo strips authorization headers and tokens from log messages
func RedactSensitiveLogInfo(err error) error {
	if err == nil {
		return nil
	}
	orig := err.Error()
	redacted := RedactSensitiveString(orig)
	if redacted == orig {
		return err
	}
	return errors.New(redacted)
}

