package common

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/resiliency"
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

var HTTPClient = func() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}()

func GetGCPProject(req ResourceRequest) (string, error) {
	project := os.Getenv("GCP_PROJECT")
	if project == "" && req.Attributes != nil {
		if p, ok := req.Attributes["gcp_project"].(string); ok && p != "" {
			project = p
		}
	}
	if project == "" {
		project = "default-gcp-project"
	}
	return project, nil
}

func GetRegion(r string, fallback string) string {
	if r != "" {
		return r
	}
	return fallback
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
