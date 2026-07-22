package sanitizer

import (
	"regexp"
	"strings"
)

var (
	azureSanitizeRegex = regexp.MustCompile(`[^a-z0-9]`)
	awsSanitizeRegex   = regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
	gcpSanitizeRegex   = regexp.MustCompile(`[^a-z0-9\-_]`)
)

// SanitizeResourceName applies cloud-specific naming constraints to raw resource names
func SanitizeResourceName(rawName string, providerType string, resourceType string) string {
	if rawName == "" {
		return "multicloud-resource"
	}

	p := strings.ToLower(providerType)
	name := strings.TrimSpace(rawName)

	// Hardened Path Traversal Protection: Strip parent directory references
	name = strings.ReplaceAll(name, "../", "")
	name = strings.ReplaceAll(name, "..\\", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")

	switch p {
	case "azure":
		// Azure Storage Account: 3-24 chars, lowercase alphanumeric only
		if strings.Contains(resourceType, "storage") {
			name = strings.ToLower(name)
			name = azureSanitizeRegex.ReplaceAllString(name, "")
			if len(name) > 24 {
				name = name[:24]
			}
			if len(name) < 3 {
				name = name + "stg"
			}
			return name
		}

	case "aws":
		// AWS S3 Bucket: 3-63 chars, lowercase alphanumeric, dots, hyphens
		if strings.Contains(resourceType, "storage") {
			name = strings.ToLower(name)
			name = awsSanitizeRegex.ReplaceAllString(name, "")
			if len(name) > 63 {
				name = name[:63]
			}
			if len(name) < 3 {
				name = name + "-s3"
			}
			return name
		}

	case "gcp":
		// GCP GCS Bucket: 3-63 chars, lowercase alphanumeric, underscores, hyphens
		if strings.Contains(resourceType, "storage") {
			name = strings.ToLower(name)
			name = gcpSanitizeRegex.ReplaceAllString(name, "")
			if len(name) > 63 {
				name = name[:63]
			}
			if len(name) < 3 {
				name = name + "-gcs"
			}
			return name
		}
	}

	return name
}
