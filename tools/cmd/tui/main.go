package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/pricing"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Cyan      = "\033[36m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Red       = "\033[31m"
	White     = "\033[97m"
	BgCyan    = "\033[46m"
	BgBlue    = "\033[44m"
	FgBlack   = "\033[30m"
)

type ResourceSample struct {
	Name      string
	Provider  string
	Type      string
	Tier      string
	Public    bool
	Encrypted bool
}

type TFState struct {
	Resources []struct {
		Mode      string `json:"mode"`
		Type      string `json:"type"`
		Name      string `json:"name"`
		Instances []struct {
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

var activeResources []ResourceSample

var sampleResources = []ResourceSample{
	{"prod-s3-bucket", "aws", "storage_bucket", "", false, true},
	{"unencrypted-gcs-bucket", "gcp", "storage_bucket", "", false, false},
	{"app-vm-aws", "aws", "virtual_machine", "medium", false, true},
	{"app-vm-azure", "azure", "virtual_machine", "large", true, true},
}

func loadResources(statePath string) []ResourceSample {
	if statePath == "" {
		if _, err := os.Stat("terraform.tfstate"); err == nil {
			statePath = "terraform.tfstate"
		}
	}
	if statePath != "" {
		/* #nosec G304 */
		data, err := os.ReadFile(filepath.Clean(statePath))
		if err == nil {
			var state TFState
			if err := json.Unmarshal(data, &state); err == nil && len(state.Resources) > 0 {
				var parsed []ResourceSample
				for _, r := range state.Resources {
					if r.Mode != "managed" || !strings.HasPrefix(r.Type, "multicloud_") {
						continue
					}
					cleanType := strings.TrimPrefix(r.Type, "multicloud_")
					for _, inst := range r.Instances {
						pType, _ := inst.Attributes["provider_type"].(string)
						if pType == "" {
							pType = "gcp"
						}
						tier, _ := inst.Attributes["size_tier"].(string)
						resName := r.Name
						if nameAttr, ok := inst.Attributes["bucket_name"].(string); ok && nameAttr != "" {
							resName = nameAttr
						} else if nameAttr, ok := inst.Attributes["instance_name"].(string); ok && nameAttr != "" {
							resName = nameAttr
						} else if nameAttr, ok := inst.Attributes["function_name"].(string); ok && nameAttr != "" {
							resName = nameAttr
						} else if nameAttr, ok := inst.Attributes["balancer_name"].(string); ok && nameAttr != "" {
							resName = nameAttr
						}
						isPublic := false
						if pub, ok := inst.Attributes["associate_public_ip"].(bool); ok {
							isPublic = pub
						} else if pubBlock, ok := inst.Attributes["public_access_block"].(bool); ok {
							isPublic = !pubBlock
						}
						isEncrypted := true
						if enc, ok := inst.Attributes["encryption_enabled"].(bool); ok {
							isEncrypted = enc
						}
						parsed = append(parsed, ResourceSample{
							Name:      resName,
							Provider:  pType,
							Type:      cleanType,
							Tier:      tier,
							Public:    isPublic,
							Encrypted: isEncrypted,
						})
					}
				}
				if len(parsed) > 0 {
					return parsed
				}
			}
		}
	}
	return sampleResources
}

func formatProviderBadge(provider string) string {
	switch strings.ToLower(provider) {
	case "gcp":
		return fmt.Sprintf("%s%s GCP %s", Cyan, Bold, Reset)
	case "aws":
		return fmt.Sprintf("%s%s AWS %s", Yellow, Bold, Reset)
	case "azure":
		return fmt.Sprintf("%s%s AZURE %s", Blue, Bold, Reset)
	default:
		return strings.ToUpper(provider)
	}
}

func main() {
	stateFlag := flag.String("state", "", "Path to terraform.tfstate file")
	flag.Parse()

	activeResources = loadResources(*stateFlag)

	// If non-interactive terminal (e.g. piped or automated unit test), print summary once and exit
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		renderSummary()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		printMenu()
		fmt.Printf(" %sSelect Option [1-5]:%s ", Bold, Reset)
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
			fmt.Printf("%sExiting Interactive Infrastructure Inspector. Goodbye!%s\n", Dim, Reset)
			return
		default:
			fmt.Printf("%sInvalid choice. Please select an option between 1 and 5.%s\n", Red, Reset)
		}

		fmt.Println()
		fmt.Printf("%sPress ENTER to return to menu...%s", Dim, Reset)
		scanner.Scan()
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("╭──────────────────────────────────────────────────────────────────────────────╮")
	fmt.Printf("│  %s%s 🌐 MULTI-CLOUD TERRAFORM PROVIDER — INFRASTRUCTURE INSPECTOR%s         │\n", Bold, Cyan, Reset)
	fmt.Println("├──────────────────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  %s[1]%s 💰  Inspect Monthly Cost Estimates per Cloud Provider               │\n", Green, Reset)
	fmt.Printf("│  %s[2]%s 🛡️   Run CIS & SOC 2 Security Audit Scan                            │\n", Yellow, Reset)
	fmt.Printf("│  %s[3]%s ⚡  View Cross-Cloud Arm Architecture Cost Optimizations            │\n", Cyan, Reset)
	fmt.Printf("│  %s[4]%s 📊  View Full Infrastructure Dashboard Summary                      │\n", Magenta, Reset)
	fmt.Printf("│  %s[5]%s 🚪  Exit Inspector                                                  │\n", Dim, Reset)
	fmt.Println("╰──────────────────────────────────────────────────────────────────────────────╯")
}

func renderCostTable() {
	fmt.Printf("%s%s┌──────────────────────────────────┬──────────────────┬──────────────────────┬───────────────────────────────┐%s\n", Cyan, Bold, Reset)
	fmt.Printf("%s%s│  💰 MONTHLY COST ESTIMATES PER RESOURCE                                                                      │%s\n", Cyan, Bold, Reset)
	fmt.Printf("%s%s├──────────────────────────────────┼──────────────────┼──────────────────────┼───────────────────────────────┤%s\n", Cyan, Bold, Reset)
	fmt.Printf("%s│ %-32s │ %-8s │ %-20s │ %-29s │%s\n", Bold, "RESOURCE NAME", "CLOUD", "UNIFIED RESOURCE TYPE", "MONTHLY COST (USD)", Reset)
	fmt.Printf("%s├──────────────────────────────────┼──────────────────┼──────────────────────┼───────────────────────────────┤%s\n", Dim, Reset)

	totalCost := 0.0
	for _, r := range activeResources {
		cleanName := sanitizer.SanitizeResourceName(r.Name, r.Provider, r.Type)
		cost := pricing.EstimateMonthlyCost(r.Provider, r.Type, r.Tier)
		totalCost += cost

		var rawCostStr, colorCode string
		if cost == 0.0 {
			rawCostStr = "$0.00 (Free Tier / Included)"
			colorCode = Dim
		} else {
			rawCostStr = fmt.Sprintf("$%.2f USD/mo", cost)
			colorCode = Green
		}

		pBadge := formatProviderBadge(r.Provider)
		fmt.Printf("│ %s%-34s%s │ %-17s │ %s%-20s%s │ %s%-29s%s │\n", White, cleanName, Reset, pBadge, Dim, r.Type, Reset, colorCode, rawCostStr, Reset)
	}

	fmt.Printf("%s├──────────────────────────────────┴──────────────────┴──────────────────────┴───────────────────────────────┤%s\n", Dim, Reset)
	fmt.Printf("│ %s%s TOTAL ESTIMATED MONTHLY SPEND: %s$%.2f USD/month%s                                                   │\n", Bold, White, Green, totalCost, Reset)
	fmt.Printf("%s%s└─────────────────────────────────────────────────────────────────────────────────────────────────────────────┘%s\n", Cyan, Bold, Reset)
}

func renderSecurityAudit() {
	fmt.Printf("%s%s🛡️  SECURITY AUDIT & COMPLIANCE FINDINGS%s\n", Yellow, Bold, Reset)
	fmt.Println("──────────────────────────────────────────────────────────────────────────────")

	findingsCount := 0
	for _, r := range activeResources {
		findings := security.AuditResource(r.Provider, r.Type, r.Name, r.Public, r.Encrypted)
		for _, f := range findings {
			findingsCount++
			sevColor := Yellow
			if f.Severity == "CRITICAL" || f.Severity == "HIGH" {
				sevColor = Red
			}
			fmt.Printf("  %s[%s]%s %s%s%s: %s\n", sevColor, f.Severity, Reset, Bold, f.RuleID, Reset, f.Message)
		}
	}

	if findingsCount == 0 {
		fmt.Printf("  %s✔ 100%% COMPLIANT: 0 CIS / SOC 2 security violations detected.%s\n", Green, Reset)
	}
}

func renderCostOptimizations() {
	fmt.Printf("%s%s⚡ CROSS-CLOUD COST OPTIMIZATION RECOMMENDATIONS%s\n", Cyan, Bold, Reset)
	fmt.Println("──────────────────────────────────────────────────────────────────────────────")

	recsCount := 0
	for _, r := range activeResources {
		if rec := pricing.RecommendCostOptimizations(r.Provider, r.Name, r.Tier); rec != nil {
			recsCount++
			fmt.Printf("  %s💡%s %s\n", Green, Reset, rec.Message)
		}
	}

	if recsCount == 0 {
		fmt.Printf("  %s✔ OPTIMIZED: All provisioned workloads are running on optimal compute architectures.%s\n", Green, Reset)
	}
}

func renderSummary() {
	renderCostTable()
	fmt.Println()
	renderSecurityAudit()
	fmt.Println()
	renderCostOptimizations()
}
