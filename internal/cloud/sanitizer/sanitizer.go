package sanitizer

import (
	"strings"
)

func filterASCII(s string, lower bool, allowDot, allowHyphen, allowUnderscore bool) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if lower && c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (!lower && c >= 'A' && c <= 'Z') ||
			(allowDot && c == '.') || (allowHyphen && c == '-') || (allowUnderscore && c == '_') {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// SanitizeResourceName applies cloud-specific naming constraints to raw resource names
func SanitizeResourceName(rawName string, providerType string, resourceType string) string {
	name := strings.TrimSpace(rawName)
	if name == "" {
		return "multicloud-resource"
	}

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

	isStorage := strings.Contains(resourceType, "storage")
	switch {
	case strings.EqualFold(providerType, "azure") && isStorage:
		// Azure Storage Account: 3-24 chars, lowercase alphanumeric only
		name = filterASCII(name, true, false, false, false)
		if len(name) > 24 {
			name = name[:24]
		}
		if len(name) < 3 {
			name = name + "stg"
		}
		return name

	case strings.EqualFold(providerType, "aws") && isStorage:
		// AWS S3 Bucket: 3-63 chars, lowercase alphanumeric, dots, hyphens
		name = filterASCII(name, true, true, true, true)
		if len(name) > 63 {
			name = name[:63]
		}
		if len(name) < 3 {
			name = name + "-s3"
		}
		return name

	case strings.EqualFold(providerType, "gcp") && isStorage:
		// GCP GCS Bucket: 3-63 chars, lowercase alphanumeric, underscores, hyphens
		name = filterASCII(name, true, false, true, true)
		if len(name) > 63 {
			name = name[:63]
		}
		if len(name) < 3 {
			name = name + "-gcs"
		}
		return name
	}

	name = filterASCII(name, false, true, true, true)
	name = strings.Trim(name, ".-_")
	if name == "" {
		return "multicloud-resource"
	}
	if len(name) > 128 {
		name = name[:128]
	}

	return name
}

func isCleanIdentifier(s string) bool {
	if len(s) == 0 || len(s) > 128 {
		return false
	}
	if s[0] == '-' || s[0] == '_' || s[len(s)-1] == '-' || s[len(s)-1] == '_' {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

// SanitizeCloudIdentifier restricts regions, project IDs, subscription IDs, and resource groups
// to safe alphanumeric, hyphen, and underscore characters to prevent URL host/authority injection.
func SanitizeCloudIdentifier(raw string, fallback string) string {
	cleaned := strings.TrimSpace(raw)
	if isCleanIdentifier(cleaned) {
		return cleaned
	}
	for strings.Contains(cleaned, "..") {
		cleaned = strings.ReplaceAll(cleaned, "..", "")
	}
	cleaned = filterASCII(cleaned, false, false, true, true)
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
	hasControl := false
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] == 0x7f {
			hasControl = true
			break
		}
	}
	if !hasControl {
		return s
	}
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
}

