package common

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"time"
)

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

var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	},
}

func GetGCPProject(req ResourceRequest) (string, error) {
	project := os.Getenv("GCP_PROJECT")
	if project == "" {
		if p, ok := req.Attributes["gcp_project"].(string); ok && p != "" {
			project = p
		}
	}
	if project == "" {
		project = "default-gcp-project"
	}
	return url.PathEscape(project), nil
}

func GetRegion(r string, fallback string) string {
	if r != "" {
		return r
	}
	return fallback
}

func SanitizeErrorBody(body []byte) string {
	str := string(body)
	if len(str) > 500 {
		return str[:500] + "... (truncated)"
	}
	return str
}
