package sanitizer

import (
	"testing"
)

func TestSanitizeResourceName(t *testing.T) {
	tests := []struct {
		name         string
		rawName      string
		providerType string
		resourceType string
		expected     string
	}{
		{"Azure Storage Invalid Chars", "My_Storage-Account!", "azure", "storage_bucket", "mystorageaccount"},
		{"AWS S3 Storage Truncation", "a-very-long-bucket-name-that-exceeds-the-maximum-allowed-length-for-aws-s3-buckets-in-terraform", "aws", "storage_bucket", "a-very-long-bucket-name-that-exceeds-the-maximum-allowed-length"},
		{"GCP GCS Bucket Valid", "my-gcp-gcs-bucket-1", "gcp", "storage_bucket", "my-gcp-gcs-bucket-1"},
		{"Empty String Fallback", "", "aws", "storage_bucket", "multicloud-resource"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeResourceName(tt.rawName, tt.providerType, tt.resourceType)
			if result != tt.expected {
				t.Errorf("SanitizeResourceName(%q, %q, %q) = %q; want %q", tt.rawName, tt.providerType, tt.resourceType, result, tt.expected)
			}
		})
	}
}
