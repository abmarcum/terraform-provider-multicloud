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
	secretParamRegex = regexp.MustCompile(`(?i)((?:client_secret|access_token|secret_key|aws_secret_access_key|api_key|x-amz-security-token)\s*[=:]\s*["']?)[^\s"',;&]+`)
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

// RedactSensitiveString strips authorization headers, signatures, and secret parameters from raw strings
func RedactSensitiveString(msg string) string {
	if msg == "" {
		return ""
	}
	msg = bearerTokenRegex.ReplaceAllString(msg, "${1}[REDACTED]")
	msg = awsSigV4Regex.ReplaceAllString(msg, "${1} [REDACTED]")
	msg = secretParamRegex.ReplaceAllString(msg, "${1}[REDACTED]")
	return msg
}

// RedactSensitiveLogInfo strips authorization headers and tokens from log messages
func RedactSensitiveLogInfo(err error) error {
	if err == nil {
		return nil
	}
	redacted := RedactSensitiveString(err.Error())
	if redacted == err.Error() {
		return err
	}
	return errors.New(redacted)
}
