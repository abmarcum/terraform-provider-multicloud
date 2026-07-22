package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/pricing"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
)

type ResourceSample struct {
	Name      string
	Provider  string
	Type      string
	Tier      string
	Public    bool
	Encrypted bool
}

var sampleResources = []ResourceSample{
	{"prod-s3-bucket", "aws", "storage_bucket", "", false, true},
	{"unencrypted-gcs-bucket", "gcp", "storage_bucket", "", false, false},
	{"app-vm-aws", "aws", "virtual_machine", "medium", false, true},
	{"app-vm-azure", "azure", "virtual_machine", "large", true, true},
}

func main() {
	// If non-interactive terminal (e.g. piped or automated unit test), print summary once and exit
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		renderSummary()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		printMenu()
		fmt.Print("Enter option [1-5]: ")
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())
		fmt.Println()

		switch choice {
		case "1":
			renderCostTable()
		case "2":
			renderSecurityAudit()
		case "3":
			renderCostOptimizations()
		case "4":
			renderSummary()
		case "5", "q", "quit", "exit":
			fmt.Println("Exiting Interactive Infrastructure Inspector. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select an option between 1 and 5.")
		}

		fmt.Println()
		fmt.Print("Press ENTER to return to menu...")
		scanner.Scan()
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("======================================================================")
	fmt.Println("  MULTI-CLOUD TERRAFORM PROVIDER - INTERACTIVE INFRASTRUCTURE INSPECTOR")
	fmt.Println("======================================================================")
	fmt.Println("1. Inspect Monthly Cost Estimates per Cloud Provider")
	fmt.Println("2. Run CIS & SOC 2 Security Audit Scan")
	fmt.Println("3. View Cross-Cloud Arm Architecture Cost Optimizations")
	fmt.Println("4. View Full Infrastructure Dashboard Summary")
	fmt.Println("5. Exit")
	fmt.Println("----------------------------------------------------------------------")
}

func renderCostTable() {
	fmt.Println("MONTHLY COST ESTIMATES PER RESOURCE:")
	fmt.Printf("%-24s %-10s %-20s %-12s\n", "RESOURCE NAME", "CLOUD", "TYPE", "ESTIMATED COST")
	fmt.Println("----------------------------------------------------------------------")

	totalCost := 0.0
	for _, r := range sampleResources {
		cleanName := sanitizer.SanitizeResourceName(r.Name, r.Provider, r.Type)
		cost := pricing.EstimateMonthlyCost(r.Provider, r.Type, r.Tier)
		totalCost += cost
		fmt.Printf("%-24s %-10s %-20s $%-11.2f\n", cleanName, r.Provider, r.Type, cost)
	}

	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("TOTAL ESTIMATED MONTHLY SPEND: $%.2f USD/month\n", totalCost)
}

func renderSecurityAudit() {
	fmt.Println("SECURITY AUDIT COMPLIANCE FINDINGS:")
	fmt.Println("----------------------------------------------------------------------")
	for _, r := range sampleResources {
		findings := security.AuditResource(r.Provider, r.Type, r.Name, r.Public, r.Encrypted)
		for _, f := range findings {
			fmt.Printf("  [%s] %s: %s\n", f.Severity, f.RuleID, f.Message)
		}
	}
}

func renderCostOptimizations() {
	fmt.Println("COST OPTIMIZATION RECOMMENDATIONS:")
	fmt.Println("----------------------------------------------------------------------")
	for _, r := range sampleResources {
		if rec := pricing.RecommendCostOptimizations(r.Provider, r.Name, r.Tier); rec != nil {
			fmt.Printf("  %s\n", rec.Message)
		}
	}
}

func renderSummary() {
	renderCostTable()
	fmt.Println()
	renderSecurityAudit()
	fmt.Println()
	renderCostOptimizations()
}
