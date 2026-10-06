package sanitizer

import (
	"regexp"
	"strings"
)

var (
	azureSanitizeRegex   = regexp.MustCompile(`[^a-z0-9]`)
	awsSanitizeRegex     = regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
	gcpSanitizeRegex     = regexp.MustCompile(`[^a-z0-9\-_]`)
	defaultSanitizeRegex = regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
)

// SanitizeResourceName applies cloud-specific naming constraints to raw resource names
func SanitizeResourceName(rawName string, providerType string, resourceType string) string {
	name := strings.TrimSpace(rawName)
	if name == "" {
		return "multicloud-resource"
	}

	p := strings.ToLower(providerType)

	// Hardened Path Traversal Protection: Recursively strip parent directory references and path separators
	for strings.Contains(name, "..") || strings.ContainsAny(name, "/\\") {
		name = strings.ReplaceAll(name, "../", "")
		name = strings.ReplaceAll(name, "..\\", "")
		name = strings.ReplaceAll(name, "..", "")
		name = strings.ReplaceAll(name, "/", "")
		name = strings.ReplaceAll(name, "\\", "")
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return "multicloud-resource"
	}

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

	name = defaultSanitizeRegex.ReplaceAllString(name, "")
	name = strings.Trim(name, ".-_")
	if name == "" {
		return "multicloud-resource"
	}
	if len(name) > 128 {
		name = name[:128]
	}

	return name
}

var identifierSanitizeRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

// SanitizeCloudIdentifier restricts regions, project IDs, subscription IDs, and resource groups
// to safe alphanumeric, hyphen, and underscore characters to prevent URL host/authority injection.
func SanitizeCloudIdentifier(raw string, fallback string) string {
	cleaned := strings.TrimSpace(raw)
	for strings.Contains(cleaned, "..") {
		cleaned = strings.ReplaceAll(cleaned, "..", "")
	}
	cleaned = identifierSanitizeRegex.ReplaceAllString(cleaned, "")
	cleaned = strings.Trim(cleaned, "-_")
	if cleaned == "" {
		return fallback
	}
	if len(cleaned) > 128 {
		cleaned = cleaned[:128]
	}
	return cleaned
}

// StripControlChars removes ASCII control characters and ANSI escape sequences from untrusted strings
// before rendering them in CLI or TUI output.
func StripControlChars(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
}
