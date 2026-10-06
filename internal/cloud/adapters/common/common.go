package common

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/resiliency"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
)

// ErrNotFound indicates the target cloud resource does not exist (HTTP 404)
var ErrNotFound = errors.New("cloud resource not found")

// ClientConfigProvider exposes provider-level settings and credentials to resources without import cycles
type ClientConfigProvider interface {
	IsMockMode() bool
	DefaultAttrs(providerType string) map[string]interface{}
}

type ResourceRequest struct {
	ResourceName string
	ResourceType string
	ProviderType string
	Region       string
	Attributes   map[string]interface{}
}

type ResourceResponse struct {
	ID         string
	Status     string
	Attributes map[string]interface{}
}

// NewHardenedHTTPClient creates an HTTP client with TLS 1.2+ and connection pooling tuned for concurrent Terraform operations.
func NewHardenedHTTPClient(timeout time.Duration) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 32
	transport.IdleConnTimeout = 90 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

var HTTPClient = NewHardenedHTTPClient(10 * time.Second)

func GetGCPProject(req ResourceRequest) (string, error) {
	project := os.Getenv("GCP_PROJECT")
	if project == "" && req.Attributes != nil {
		if p, ok := req.Attributes["gcp_project"].(string); ok && p != "" {
			project = p
		}
	}
	return sanitizer.SanitizeCloudIdentifier(project, "default-gcp-project"), nil
}

func GetRegion(r string, fallback string) string {
	if r != "" {
		return sanitizer.SanitizeCloudIdentifier(r, fallback)
	}
	return fallback
}

// IsSensitiveOrInternalKey returns true for internal provider keys and cloud credentials
// that must not be overridden via resource extra_config or echoed into ResourceResponse.Attributes.
func IsSensitiveOrInternalKey(k string) bool {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case "_client_manager",
		"aws_access_key", "aws_secret_key", "aws_profile",
		"gcp_credentials",
		"azure_client_secret", "azure_bearer_token", "azure_tenant_id", "azure_client_id",
		"mock_mode", "provider_default_region":
		return true
	default:
		return false
	}
}

// IsSafeOutboundURL validates that an endpoint URL uses https (or localhost http for testing)
// and blocks cloud metadata link-local addresses to prevent SSRF.
func IsSafeOutboundURL(rawURL string) bool {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "169.254.169.254" || host == "fd00:ec2::254" || host == "metadata.google.internal" || strings.HasPrefix(host, "169.254.") {
		return false
	}
	if u.Scheme == "https" {
		return true
	}
	if u.Scheme == "http" && (host == "127.0.0.1" || host == "localhost" || host == "::1") {
		return true
	}
	return false
}

func SanitizeErrorBody(body []byte) string {
	str := resiliency.RedactSensitiveString(string(body))
	if len(str) > 500 {
		return str[:500] + "... (truncated)"
	}
	return str
}

func IsMockMode() bool {
	return os.Getenv("MULTICLOUD_MOCK_MODE") == "true"
}

func IsRequestMockMode(req ResourceRequest) bool {
	if IsMockMode() {
		return true
	}
	if req.Attributes != nil {
		if m, ok := req.Attributes["mock_mode"].(bool); ok && m {
			return true
		}
	}
	return false
}
