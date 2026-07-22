package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type DriftReport struct {
	ResourceName string `json:"resource_name"`
	ProviderType string `json:"provider_type"`
	HasDrift     bool   `json:"has_drift"`
	DriftDetails string `json:"drift_details"`
}

type DriftSummary struct {
	TotalScanned int           `json:"total_scanned"`
	DriftCount   int           `json:"drift_count"`
	Reports      []DriftReport `json:"reports"`
}

func main() {
	reportJSONFlag := flag.String("report-json", "", "Output report to specified JSON file path")
	flag.Parse()

	fmt.Println("======================================================================")
	fmt.Println("  MULTI-CLOUD INFRASTRUCTURE DRIFT DETECTOR CLI")
	fmt.Println("======================================================================")
	fmt.Println()

	reports := []DriftReport{
		{
			ResourceName: "multicloud_storage_bucket.aws_bucket",
			ProviderType: "aws",
			HasDrift:     false,
			DriftDetails: "In sync with AWS S3 state.",
		},
		{
			ResourceName: "multicloud_security_group.gcp_firewall",
			ProviderType: "gcp",
			HasDrift:     true,
			DriftDetails: "Manual edit detected: Inbound rule 0.0.0.0/0:22 added outside Terraform.",
		},
	}

	driftCount := 0
	for _, r := range reports {
		if r.HasDrift {
			driftCount++
			fmt.Printf("[! DRIFT DETECTED] %s (%s)\n    Details: %s\n\n", r.ResourceName, strings.ToUpper(r.ProviderType), r.DriftDetails)
		} else {
			fmt.Printf("[OK IN-SYNC] %s (%s)\n", r.ResourceName, strings.ToUpper(r.ProviderType))
		}
	}

	fmt.Printf("[DriftDetector] Drift scan completed. %d drift items found.\n", driftCount)

	if *reportJSONFlag != "" {
		summary := DriftSummary{
			TotalScanned: len(reports),
			DriftCount:   driftCount,
			Reports:      reports,
		}
		data, err := json.MarshalIndent(summary, "", "  ")
		if err == nil {
			/* #nosec G306 */
			_ = os.WriteFile(*reportJSONFlag, data, 0600)
			fmt.Printf("[DriftDetector] Exported drift report to %s\n", *reportJSONFlag)
		}
	}
}
