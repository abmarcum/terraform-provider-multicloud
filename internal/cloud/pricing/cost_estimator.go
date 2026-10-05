package pricing

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AzureRetailPriceResponse models Azure Retail Prices API JSON payload
type AzureRetailPriceResponse struct {
	Items []struct {
		UnitPrice     float64 `json:"unitPrice"`
		UnitOfMeasure string  `json:"unitOfMeasure"`
		CurrencyCode  string  `json:"currencyCode"`
	} `json:"Items"`
}

// AWSPriceResponse models AWS Price List bulk index / offer JSON payload
type AWSPriceResponse struct {
	PriceList []string `json:"PriceList"`
	Terms     struct {
		OnDemand map[string]map[string]struct {
			PriceDimensions map[string]struct {
				PricePerUnit map[string]string `json:"pricePerUnit"`
			} `json:"priceDimensions"`
		} `json:"OnDemand"`
	} `json:"terms"`
}

// GCPCatalogResponse models GCP Cloud Billing Catalog API SKU payload
type GCPCatalogResponse struct {
	SKUs []struct {
		Description string `json:"description"`
		PricingInfo []struct {
			PricingExpression struct {
				TieredRates []struct {
					UnitPrice struct {
						Nanos int64  `json:"nanos"`
						Units string `json:"units"`
					} `json:"unitPrice"`
				} `json:"tieredRates"`
			} `json:"pricingExpression"`
		} `json:"pricingInfo"`
	} `json:"skus"`
}

var liveHTTPClient = func() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	return &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
	}
}()

type cachedPrice struct {
	price     float64
	expiresAt time.Time
}

var (
	priceCacheMu sync.RWMutex
	priceCache   = make(map[string]cachedPrice)
)

func getCachedPrice(key string) (float64, bool) {
	now := time.Now()
	priceCacheMu.RLock()
	defer priceCacheMu.RUnlock()
	if entry, ok := priceCache[key]; ok && now.Before(entry.expiresAt) {
		return entry.price, true
	}
	return 0, false
}

func setCachedPrice(key string, price float64) {
	priceCacheMu.Lock()
	priceCache[key] = cachedPrice{
		price:     price,
		expiresAt: time.Now().Add(15 * time.Minute),
	}
	priceCacheMu.Unlock()
}

// FetchLiveAzurePrice fetches live pricing from Azure Retail Prices API
func FetchLiveAzurePrice(skuName string) (float64, error) {
	return FetchLiveAzurePriceWithContext(context.Background(), skuName)
}

// FetchLiveAzurePriceWithContext fetches live pricing from Azure Retail Prices API with context and TTL caching
func FetchLiveAzurePriceWithContext(ctx context.Context, skuName string) (float64, error) {
	cacheKey := "azure:" + skuName
	if cached, ok := getCachedPrice(cacheKey); ok {
		return cached, nil
	}

	filter := fmt.Sprintf("armSkuName eq '%s' and priceType eq 'Consumption'", skuName)
	endpoint := fmt.Sprintf("https://prices.azure.com/api/retail/arm/prices?$filter=%s", url.QueryEscape(filter))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := liveHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("azure pricing API HTTP %d", resp.StatusCode)
	}

	var data AzureRetailPriceResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Items) > 0 && data.Items[0].UnitPrice > 0 {
		monthly := data.Items[0].UnitPrice * 730.0
		setCachedPrice(cacheKey, monthly)
		return monthly, nil
	}

	return 0, fmt.Errorf("no pricing found for SKU %s", skuName)
}

// FetchLiveAWSPrice fetches pricing rate for AWS EC2 instance types, querying AWS Price List endpoint when configured
func FetchLiveAWSPrice(instanceType string) (float64, error) {
	cacheKey := "aws:" + instanceType
	if cached, ok := getCachedPrice(cacheKey); ok {
		return cached, nil
	}

	if offerURL := os.Getenv("AWS_PRICING_OFFER_URL"); offerURL != "" && os.Getenv("MULTICLOUD_MOCK_MODE") != "true" {
		/* #nosec G107 G704 */
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, offerURL, nil)
		if err == nil {
			/* #nosec G107 G704 */
			resp, err := liveHTTPClient.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var awsData AWSPriceResponse
					if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&awsData); err == nil {
						for _, skuMap := range awsData.Terms.OnDemand {
							for _, term := range skuMap {
								for _, dim := range term.PriceDimensions {
									if usdStr, ok := dim.PricePerUnit["USD"]; ok {
										if hourly, err := strconv.ParseFloat(usdStr, 64); err == nil && hourly > 0 {
											monthly := hourly * 730.0
											setCachedPrice(cacheKey, monthly)
											return monthly, nil
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	hourlyRates := map[string]float64{
		"t3.nano":     0.0052,
		"t3.micro":    0.0208,
		"t3.small":    0.0208,
		"t3.medium":   0.0416,
		"t3.large":    0.0832,
		"t3.xlarge":   0.1664,
		"t4g.small":   0.0168,
		"t4g.medium":  0.0336,
		"t4g.large":   0.0672,
		"m6i.large":   0.0960,
		"m6i.xlarge":  0.1920,
		"m6i.2xlarge": 0.3840,
		"c6i.large":   0.0850,
		"r6i.large":   0.1260,
	}
	if hourly, ok := hourlyRates[strings.ToLower(instanceType)]; ok {
		monthly := hourly * 730.0
		setCachedPrice(cacheKey, monthly)
		return monthly, nil
	}
	return 0, fmt.Errorf("unsupported AWS instance type %s", instanceType)
}

// FetchLiveGCPPrice fetches pricing rate for GCP Compute Engine machine types, querying Cloud Billing Catalog API when configured
func FetchLiveGCPPrice(machineType string) (float64, error) {
	cacheKey := "gcp:" + machineType
	if cached, ok := getCachedPrice(cacheKey); ok {
		return cached, nil
	}

	if apiKey := os.Getenv("GCP_BILLING_API_KEY"); apiKey != "" && os.Getenv("MULTICLOUD_MOCK_MODE") != "true" {
		endpoint := fmt.Sprintf("https://cloudbilling.googleapis.com/v1/services/6F81-5844-456A/skus?key=%s&pageSize=20", url.QueryEscape(apiKey))
		/* #nosec G107 G704 */
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
		if err == nil {
			/* #nosec G107 G704 */
			resp, err := liveHTTPClient.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var catalog GCPCatalogResponse
					if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&catalog); err == nil {
						for _, sku := range catalog.SKUs {
							if strings.Contains(strings.ToLower(sku.Description), strings.ToLower(machineType)) && len(sku.PricingInfo) > 0 {
								rates := sku.PricingInfo[0].PricingExpression.TieredRates
								if len(rates) > 0 {
									units, _ := strconv.ParseFloat(rates[0].UnitPrice.Units, 64)
									hourly := units + float64(rates[0].UnitPrice.Nanos)/1e9
									if hourly > 0 {
										monthly := hourly * 730.0
										setCachedPrice(cacheKey, monthly)
										return monthly, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}

	hourlyRates := map[string]float64{
		"e2-micro":       0.0198,
		"e2-small":       0.0198,
		"e2-medium":      0.0408,
		"e2-standard-2":  0.0816,
		"e2-standard-4":  0.1632,
		"n2-standard-2":  0.0971,
		"n2-standard-4":  0.1942,
		"n2-standard-8":  0.3884,
		"t2a-standard-1": 0.0385,
		"t2a-standard-2": 0.0770,
		"t2a-standard-4": 0.1540,
	}
	if hourly, ok := hourlyRates[strings.ToLower(machineType)]; ok {
		monthly := hourly * 730.0
		setCachedPrice(cacheKey, monthly)
		return monthly, nil
	}
	return 0, fmt.Errorf("unsupported GCP machine type %s", machineType)
}

// EstimateMonthlyCost calculates estimated monthly cost using live pricing feeds with resilient offline fallback
func EstimateMonthlyCost(providerType string, resourceType string, sizeTier string) float64 {
	p := strings.ToLower(providerType)
	r := strings.ToLower(resourceType)
	tier := strings.ToLower(sizeTier)

	// 1. Live Azure Retail Prices API Feed (skipped in mock mode)
	if p == "azure" && r == "virtual_machine" && os.Getenv("MULTICLOUD_MOCK_MODE") != "true" {
		sku := "Standard_B2s"
		if tier == "small" {
			sku = "Standard_B1s"
		} else if tier == "large" {
			sku = "Standard_B2ms"
		}

		if livePrice, err := FetchLiveAzurePrice(sku); err == nil && livePrice > 0 {
			return livePrice
		}
	}

	// 2. Live AWS Price List Feed
	if p == "aws" && r == "virtual_machine" {
		inst := "t3.medium"
		if tier == "small" {
			inst = "t3.micro"
		} else if tier == "large" {
			inst = "t3.large"
		}

		if livePrice, err := FetchLiveAWSPrice(inst); err == nil && livePrice > 0 {
			return livePrice
		}
	}

	// 3. Live GCP Billing Catalog Feed
	if p == "gcp" && r == "virtual_machine" {
		machine := "e2-medium"
		if tier == "small" {
			machine = "e2-micro"
		} else if tier == "large" {
			machine = "e2-standard-2"
		}

		if livePrice, err := FetchLiveGCPPrice(machine); err == nil && livePrice > 0 {
			return livePrice
		}
	}

	// Resilient Offline Fallback Matrix (Used if any API network call fails or times out)
	return EstimateOfflineMonthlyCost(p, r, tier)
}

// EstimateOfflineMonthlyCost provides fallback baseline prices across unified multi-cloud resource types
func EstimateOfflineMonthlyCost(p, r, tier string) float64 {
	switch r {
	case "virtual_machine", "custom_machine_type", "bastion_host":
		switch tier {
		case "small":
			switch p {
			case "aws":
				return 15.20 // t3.micro
			case "gcp":
				return 14.50 // e2-micro
			case "azure":
				return 14.80 // Standard_B1s
			}
		case "medium", "":
			switch p {
			case "aws":
				return 30.40 // t3.medium
			case "gcp":
				return 29.80 // e2-medium
			case "azure":
				return 30.10 // Standard_B2s
			}
		case "large":
			switch p {
			case "aws":
				return 60.80 // t3.large
			case "gcp":
				return 59.60 // e2-standard-2
			case "azure":
				return 60.20 // Standard_B2ms
			}
		}

	case "db_instance":
		switch p {
		case "aws":
			return 45.00 // db.t3.medium
		case "gcp":
			return 43.50 // db-custom-2-7680
		case "azure":
			return 44.00 // GP_Gen5_2
		}

	case "kubernetes_cluster":
		switch p {
		case "aws":
			return 73.00 // EKS control plane
		case "gcp":
			return 73.00 // GKE control plane
		case "azure":
			return 0.00 // AKS free tier control plane
		}

	case "storage_bucket", "storage_inventory_report", "storage_transfer_job", "data_sync":
		switch p {
		case "aws":
			return 2.30 // S3 Standard 100GB baseline
		case "gcp":
			return 2.00 // GCS Standard 100GB baseline
		case "azure":
			return 1.84 // Azure Blob Hot 100GB baseline
		}

	case "load_balancer", "nat_gateway", "vpn_gateway", "transit_gateway", "global_anycast_ip":
		switch p {
		case "aws":
			return 22.50
		case "gcp":
			return 18.00
		case "azure":
			return 21.90
		}

	case "cache_cluster", "search_index", "streaming_cluster":
		switch p {
		case "aws":
			return 35.00
		case "gcp":
			return 33.50
		case "azure":
			return 34.20
		}

	case "data_warehouse", "ai_endpoint", "feature_store", "data_pipeline":
		switch p {
		case "aws":
			return 180.00
		case "gcp":
			return 165.00
		case "azure":
			return 175.00
		}

	case "virtual_network", "subnet", "security_group", "route_table", "iam_role", "kms_policy", "workload_identity_pool", "identity_federation":
		return 0.00
	}

	return 10.00
}

// FormatCostReport formats estimated cost summary for Terraform plan output
func FormatCostReport(providerType string, resourceName string, cost float64) string {
	return fmt.Sprintf("[CostEstimator] Resource '%s' on %s estimated at $%.2f USD/month", resourceName, strings.ToUpper(providerType), cost)
}
