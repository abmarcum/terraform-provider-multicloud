package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSAdapter struct{}

var awsConfigCache sync.Map

func loadAWSConfig(ctx context.Context, region string, attrs map[string]interface{}) (aws.Config, error) {
	if attrs != nil {
		ak, _ := attrs["aws_access_key"].(string)
		sk, _ := attrs["aws_secret_key"].(string)
		profile, _ := attrs["aws_profile"].(string)
		if ak != "" && sk != "" {
			return config.LoadDefaultConfig(ctx,
				config.WithRegion(region),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ak, sk, "")),
			)
		}
		if profile != "" {
			return config.LoadDefaultConfig(ctx,
				config.WithRegion(region),
				config.WithSharedConfigProfile(profile),
			)
		}
	}

	if cached, ok := awsConfigCache.Load(region); ok {
		return cached.(aws.Config), nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err == nil {
		awsConfigCache.Store(region, cfg)
	}
	return cfg, err
}

func awsServiceNameForResource(resType string) string {
	switch resType {
	case "storage_bucket", "storage_inventory_report":
		return "s3"
	case "db_instance":
		return "rds"
	case "secret", "secret_rotator":
		return "secretsmanager"
	case "serverless_function", "edge_function":
		return "lambda"
	case "kubernetes_cluster":
		return "eks"
	case "container_registry":
		return "ecr"
	case "container_app":
		return "apprunner"
	case "nosql_table":
		return "dynamodb"
	case "kms_key", "kms_policy":
		return "kms"
	case "iam_role", "identity_federation", "workload_identity_pool":
		return "iam"
	case "dns_zone", "dns_record", "dns_health_check", "dns_zone_link", "dns_resolver", "dnssec", "failover_policy":
		return "route53"
	case "pubsub_topic":
		return "sns"
	case "message_queue":
		return "sqs"
	case "event_bridge":
		return "events"
	case "cdn_distribution":
		return "cloudfront"
	case "cache_cluster":
		return "elasticache"
	case "api_gateway":
		return "apigateway"
	case "graphql_api":
		return "appsync"
	case "data_warehouse":
		return "redshift"
	case "search_index":
		return "es"
	case "auto_scaling_group":
		return "autoscaling"
	case "monitoring_dashboard", "metric_alert":
		return "monitoring"
	case "log_workspace":
		return "logs"
	case "data_sync", "storage_transfer_job":
		return "datasync"
	case "waf_policy":
		return "wafv2"
	case "app_config":
		return "appconfig"
	case "ai_endpoint", "feature_store":
		return "sagemaker"
	case "streaming_cluster":
		return "kafka"
	case "security_center":
		return "securityhub"
	case "data_pipeline":
		return "glue"
	case "global_anycast_ip":
		return "globalaccelerator"
	case "load_balancer":
		return "elasticloadbalancing"
	case "shared_filesystem":
		return "elasticfilesystem"
	case "tls_certificate":
		return "acm"
	case "workflow":
		return "states"
	case "batch_compute":
		return "batch"
	case "backup_vault":
		return "backup"
	case "distributed_tracing":
		return "xray"
	case "budget_alert":
		return "budgets"
	case "vector_index":
		return "aoss"
	case "service_mesh":
		return "appmesh"
	default:
		return "ec2"
	}
}

func signAWSRequest(ctx context.Context, cfg aws.Config, httpReq *http.Request, payload []byte, region string, resType string) {
	if cfg.Credentials == nil {
		return
	}
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil || !creds.HasKeys() {
		return
	}
	sum := sha256.Sum256(payload)
	payloadHash := hex.EncodeToString(sum[:])
	signer := v4.NewSigner()
	_ = signer.SignHTTP(ctx, creds, httpReq, payloadHash, awsServiceNameForResource(resType), region, time.Now())
}

func getAWSAccountID() string {
	acc := os.Getenv("AWS_ACCOUNT_ID")
	if acc == "" {
		acc = "unknown-account"
	}
	return url.PathEscape(acc)
}

func getAWSIntelInstanceType(sizeTier string, extraAttrs map[string]interface{}) string {
	if extraAttrs != nil {
		if inst, ok := extraAttrs["instance_type"].(string); ok && inst != "" {
			return inst
		}
		if inst, ok := extraAttrs["aws_instance_type"].(string); ok && inst != "" {
			return inst
		}
		if arch, ok := extraAttrs["aws_hardware_architecture"].(string); ok && strings.ToLower(arch) == "intel" {
			switch strings.ToLower(sizeTier) {
			case "large":
				return "m6i.xlarge"
			case "medium":
				return "m6i.large"
			default:
				return "t3.medium"
			}
		}
	}
	switch strings.ToLower(sizeTier) {
	case "large":
		return "m6i.xlarge"
	case "medium":
		return "m6i.large"
	default:
		return "t3.medium"
	}
}

func getAWSServiceEndpoint(region string, resType string, name string) (string, string, []byte) {
	var endpoint string
	var method = "POST"
	var payload []byte

	escName := url.PathEscape(name)
	escQueryName := url.QueryEscape(name)
	escRegion := url.PathEscape(region)
	accID := getAWSAccountID()

	switch resType {
	case "storage_bucket", "storage_inventory_report":
		endpoint = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", escName, escRegion)
		method = "PUT"
	case "virtual_machine", "custom_machine_type", "bastion_host":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=RunInstances&ImageId=ami-0c55b159cbfafe1f0&InstanceType=m6i.large&MinCount=1&MaxCount=1&Version=2016-11-15", escRegion)
	case "block_volume":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateVolume&AvailabilityZone=%sa&Size=100&VolumeType=gp3&Version=2016-11-15", escRegion, escRegion)
	case "virtual_network", "vpc_peering":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateVpc&CidrBlock=10.0.0.0/16&Version=2016-11-15", escRegion)
	case "subnet":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateSubnet&CidrBlock=10.0.1.0/24&Version=2016-11-15", escRegion)
	case "security_group":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateSecurityGroup&GroupName=%s&Version=2016-11-15", escRegion, escQueryName)
	case "static_ip":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=AllocateAddress&Domain=vpc&Version=2016-11-15", escRegion)
	case "nat_gateway":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateNatGateway&SubnetId=subnet-default&AllocationId=eipalloc-default&Version=2016-11-15", escRegion)
	case "route_table":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateRouteTable&VpcId=vpc-default&Version=2016-11-15", escRegion)
	case "vpn_gateway":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateVpnGateway&Type=ipsec.1&Version=2016-11-15", escRegion)
	case "transit_gateway":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateTransitGateway&Version=2016-11-15", escRegion)
	case "private_endpoint":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateVpcEndpoint&VpcId=vpc-default&ServiceName=%s&Version=2016-11-15", escRegion, escQueryName)
	case "load_balancer":
		endpoint = fmt.Sprintf("https://elasticloadbalancing.%s.amazonaws.com/?Action=CreateLoadBalancer&Name=%s&Version=2015-12-01", escRegion, escQueryName)
	case "db_instance":
		endpoint = fmt.Sprintf("https://rds.%s.amazonaws.com/?Action=CreateDBInstance&DBInstanceIdentifier=%s&Version=2014-10-31", escRegion, escQueryName)
	case "secret", "secret_rotator":
		endpoint = fmt.Sprintf("https://secretsmanager.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name})
	case "serverless_function", "edge_function":
		endpoint = fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions", escRegion)
		payload, _ = json.Marshal(map[string]string{"FunctionName": name, "Runtime": "python3.11", "Role": fmt.Sprintf("arn:aws:iam::%s:role/service-role", accID)})
	case "kubernetes_cluster":
		endpoint = fmt.Sprintf("https://eks.%s.amazonaws.com/clusters", escRegion)
		payload, _ = json.Marshal(map[string]string{"name": name})
	case "container_registry":
		endpoint = fmt.Sprintf("https://api.ecr.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"repositoryName": name})
	case "container_app":
		endpoint = fmt.Sprintf("https://apprunner.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"ServiceName": name})
	case "nosql_table":
		endpoint = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"TableName": name, "AttributeDefinitions": []map[string]string{{"AttributeName": "id", "AttributeType": "S"}}})
	case "kms_key", "kms_policy":
		endpoint = fmt.Sprintf("https://kms.%s.amazonaws.com", escRegion)
		payload = []byte(`{"Description":"Multi-cloud KMS Key"}`)
	case "iam_role", "identity_federation", "workload_identity_pool":
		endpoint = fmt.Sprintf("https://iam.amazonaws.com/?Action=CreateRole&RoleName=%s&Version=2010-05-08", escQueryName)
	case "dns_zone", "dns_record", "dns_health_check", "dns_zone_link", "dns_resolver", "dnssec", "failover_policy":
		endpoint = fmt.Sprintf("https://route53.amazonaws.com/2013-04-01/hostedzone/%s", escName)
	case "pubsub_topic":
		endpoint = fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=CreateTopic&Name=%s&Version=2010-03-31", escRegion, escQueryName)
	case "message_queue":
		endpoint = fmt.Sprintf("https://sqs.%s.amazonaws.com/?Action=CreateQueue&QueueName=%s&Version=2012-11-05", escRegion, escQueryName)
	case "event_bridge":
		endpoint = fmt.Sprintf("https://events.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name})
	case "cdn_distribution":
		endpoint = "https://cloudfront.amazonaws.com/2020-05-31/distribution"
	case "cache_cluster":
		endpoint = fmt.Sprintf("https://elasticache.%s.amazonaws.com/?Action=CreateCacheCluster&CacheClusterId=%s&Version=2015-02-02", escRegion, escQueryName)
	case "api_gateway":
		endpoint = fmt.Sprintf("https://apigateway.%s.amazonaws.com/v2/apis", escRegion)
		payload, _ = json.Marshal(map[string]string{"name": name, "protocolType": "HTTP"})
	case "graphql_api":
		endpoint = fmt.Sprintf("https://appsync.%s.amazonaws.com/v1/apis", escRegion)
		payload, _ = json.Marshal(map[string]string{"name": name, "authenticationType": "API_KEY"})
	case "data_warehouse":
		endpoint = fmt.Sprintf("https://redshift.%s.amazonaws.com/?Action=CreateCluster&ClusterIdentifier=%s&NodeType=ra3.xlplus&Version=2012-12-01", escRegion, escQueryName)
	case "search_index", "vector_index":
		endpoint = fmt.Sprintf("https://es.%s.amazonaws.com/2021-01-01/opensearch/domain", escRegion)
		payload, _ = json.Marshal(map[string]string{"DomainName": name})
	case "auto_scaling_group":
		endpoint = fmt.Sprintf("https://autoscaling.%s.amazonaws.com/?Action=CreateAutoScalingGroup&AutoScalingGroupName=%s&MinSize=1&MaxSize=3&Version=2011-01-01", escRegion, escQueryName)
	case "monitoring_dashboard", "metric_alert":
		endpoint = fmt.Sprintf("https://monitoring.%s.amazonaws.com/?Action=PutDashboard&DashboardName=%s&Version=2010-08-01", escRegion, escQueryName)
	case "log_workspace":
		endpoint = fmt.Sprintf("https://logs.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"logGroupName": name})
	case "data_sync", "storage_transfer_job":
		endpoint = fmt.Sprintf("https://datasync.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name})
	case "waf_policy":
		endpoint = fmt.Sprintf("https://wafv2.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name, "Scope": "REGIONAL"})
	case "app_config":
		endpoint = fmt.Sprintf("https://appconfig.%s.amazonaws.com/applications", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name})
	case "ai_endpoint", "feature_store":
		endpoint = fmt.Sprintf("https://api.sagemaker.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"EndpointName": name})
	case "streaming_cluster":
		endpoint = fmt.Sprintf("https://kafka.%s.amazonaws.com/v1/clusters", escRegion)
		payload, _ = json.Marshal(map[string]string{"clusterName": name})
	case "security_center":
		endpoint = fmt.Sprintf("https://securityhub.%s.amazonaws.com/hub", escRegion)
		payload, _ = json.Marshal(map[string]bool{"EnableDefaultStandards": true})
	case "data_pipeline":
		endpoint = fmt.Sprintf("https://glue.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"Name": name})
	case "global_anycast_ip":
		endpoint = "https://globalaccelerator.us-west-2.amazonaws.com"
		payload, _ = json.Marshal(map[string]string{"Name": name, "IpAddressType": "IPV4"})
	case "shared_filesystem":
		endpoint = fmt.Sprintf("https://elasticfilesystem.%s.amazonaws.com/2015-02-01/file-systems", escRegion)
		payload, _ = json.Marshal(map[string]string{"CreationToken": name})
	case "tls_certificate":
		endpoint = fmt.Sprintf("https://acm.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"DomainName": name})
	case "workflow":
		endpoint = fmt.Sprintf("https://states.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"name": name})
	case "batch_compute":
		endpoint = fmt.Sprintf("https://batch.%s.amazonaws.com/v1/createcomputeenvironment", escRegion)
		payload, _ = json.Marshal(map[string]string{"computeEnvironmentName": name})
	case "backup_vault":
		endpoint = fmt.Sprintf("https://backup.%s.amazonaws.com/backup-vaults/%s", escRegion, escName)
		method = "PUT"
	case "distributed_tracing":
		endpoint = fmt.Sprintf("https://xray.%s.amazonaws.com/CreateSamplingRule", escRegion)
		payload, _ = json.Marshal(map[string]string{"RuleName": name})
	case "budget_alert":
		endpoint = fmt.Sprintf("https://budgets.amazonaws.com/accounts/%s/budgets", accID)
		payload, _ = json.Marshal(map[string]string{"BudgetName": name})
	case "service_mesh":
		endpoint = fmt.Sprintf("https://appmesh.%s.amazonaws.com/v20190125/meshes", escRegion)
		method = "PUT"
		payload, _ = json.Marshal(map[string]string{"meshName": name})
	default:
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeInstances&Version=2016-11-15", escRegion)
	}

	return endpoint, method, payload
}

func getAWSDeleteEndpoint(region string, resType string, name string) (string, string, []byte) {
	var endpoint string
	var method = "POST"
	var payload []byte

	escName := url.PathEscape(name)
	escQueryName := url.QueryEscape(name)
	escRegion := url.PathEscape(region)
	accID := getAWSAccountID()

	switch resType {
	case "storage_bucket", "storage_inventory_report":
		endpoint = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", escName, escRegion)
		method = "DELETE"
	case "virtual_network", "vpc_peering":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DeleteVpc&VpcId=%s&Version=2016-11-15", escRegion, escQueryName)
	case "subnet":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DeleteSubnet&SubnetId=%s&Version=2016-11-15", escRegion, escQueryName)
	case "security_group":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DeleteSecurityGroup&GroupName=%s&Version=2016-11-15", escRegion, escQueryName)
	case "db_instance":
		endpoint = fmt.Sprintf("https://rds.%s.amazonaws.com/?Action=DeleteDBInstance&DBInstanceIdentifier=%s&SkipFinalSnapshot=true&Version=2014-10-31", escRegion, escQueryName)
	case "secret", "secret_rotator":
		endpoint = fmt.Sprintf("https://secretsmanager.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"SecretId": name, "ForceDeleteWithoutRecovery": true})
	case "serverless_function", "edge_function":
		endpoint = fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions/%s", escRegion, escName)
		method = "DELETE"
	case "kubernetes_cluster":
		endpoint = fmt.Sprintf("https://eks.%s.amazonaws.com/clusters/%s", escRegion, escName)
		method = "DELETE"
	case "nosql_table":
		endpoint = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", escRegion)
		payload, _ = json.Marshal(map[string]string{"TableName": name})
	case "pubsub_topic":
		endpoint = fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=DeleteTopic&TopicArn=arn:aws:sns:%s:%s:%s&Version=2010-03-31", escRegion, escRegion, accID, escQueryName)
	case "message_queue":
		endpoint = fmt.Sprintf("https://sqs.%s.amazonaws.com/?Action=DeleteQueue&QueueUrl=https://sqs.%s.amazonaws.com/%s/%s&Version=2012-11-05", escRegion, escRegion, accID, escName)
	default:
		endpoint, method, payload = getAWSServiceEndpoint(region, resType, name)
		if method == "POST" && len(payload) == 0 && strings.Contains(endpoint, "ec2.") {
			endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=TerminateInstances&InstanceId.1=%s&Version=2016-11-15", escRegion, escQueryName)
		} else {
			method = "DELETE"
		}
	}

	return endpoint, method, payload
}

func getAWSReadEndpoint(region string, resType string, name string) (string, string) {
	escName := url.PathEscape(name)
	escQueryName := url.QueryEscape(name)
	escRegion := url.PathEscape(region)

	switch resType {
	case "storage_bucket", "storage_inventory_report":
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", escName, escRegion), "HEAD"
	case "virtual_network", "vpc_peering":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeVpcs&VpcId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "subnet":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeSubnets&SubnetId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "security_group":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeSecurityGroups&GroupName.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "db_instance":
		return fmt.Sprintf("https://rds.%s.amazonaws.com/?Action=DescribeDBInstances&DBInstanceIdentifier=%s&Version=2014-10-31", escRegion, escQueryName), "POST"
	case "serverless_function", "edge_function":
		return fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions/%s", escRegion, escName), "GET"
	case "kubernetes_cluster":
		return fmt.Sprintf("https://eks.%s.amazonaws.com/clusters/%s", escRegion, escName), "GET"
	default:
		ep, _, _ := getAWSServiceEndpoint(region, resType, name)
		if strings.Contains(ep, "ec2.") {
			return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeInstances&InstanceId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
		}
		return ep, "GET"
	}
}

func buildMockAttributes(provider, region string, req common.ResourceRequest) map[string]interface{} {
	attrs := map[string]interface{}{
		"region":         region,
		"provider_type":  provider,
		"failoverstatus": "PRIMARY_HEALTHY",
		"syncstatus":     "REPLICATION_ACTIVE",
		"apiendpoint":    fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", req.ResourceName, region),
		"ip_address":     fmt.Sprintf("198.51.100.%d", len(req.ResourceName)*7%250+1),
		"dns_name":       fmt.Sprintf("%s.anycast.%s.net", req.ResourceName, provider),
	}
	for k, v := range req.Attributes {
		attrs[k] = v
	}
	return attrs
}

func (a *AWSAdapter) CreateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         fmt.Sprintf("aws/%s/%s/%s", req.ResourceType, region, req.ResourceName),
			Status:     "ACTIVE",
			Attributes: buildMockAttributes("aws", region, req),
		}, nil
	}

	cfg, err := loadAWSConfig(ctx, region, req.Attributes)
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
		if err == nil {
			return common.ResourceResponse{
				ID:         fmt.Sprintf("arn:aws:s3:::%s", req.ResourceName),
				Status:     "ACTIVE",
				Attributes: buildMockAttributes("aws", region, req),
			}, nil
		}
		return common.ResourceResponse{}, fmt.Errorf("AWS S3 CreateBucket API error: %w", err)
	}

	apiEndpoint, method, payload := getAWSServiceEndpoint(region, req.ResourceType, req.ResourceName)
	respAttrs := buildMockAttributes("aws", region, req)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("failed to create AWS HTTP request: %w", err)
		}
		if len(payload) > 0 && payload[0] == '{' {
			httpReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
		} else {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		signAWSRequest(ctx, cfg, httpReq, payload, region, req.ResourceType)

		/* #nosec G107 G704 */
		resp, err := common.HTTPClient.Do(httpReq)
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("AWS live API call error: %w", err)
		}
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if resp.StatusCode >= 400 {
			return common.ResourceResponse{}, fmt.Errorf("AWS live API error (status %d) for %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
		}
		_ = json.Unmarshal(bodyBytes, &respAttrs)
	}

	return common.ResourceResponse{
		ID:         fmt.Sprintf("arn:aws:%s:%s:%s:%s", req.ResourceType, region, getAWSAccountID(), req.ResourceName),
		Status:     "ACTIVE",
		Attributes: respAttrs,
	}, nil
}

func (a *AWSAdapter) ReadResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "ACTIVE",
			Attributes: buildMockAttributes("aws", region, req),
		}, nil
	}

	cfg, err := loadAWSConfig(ctx, region, req.Attributes)
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "notfound") || strings.Contains(err.Error(), "404") {
				return common.ResourceResponse{}, fmt.Errorf("%w: AWS S3 bucket %s not found: %v", common.ErrNotFound, req.ResourceName, err)
			}
			return common.ResourceResponse{}, fmt.Errorf("AWS S3 bucket %s read error: %w", req.ResourceName, err)
		}
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "ACTIVE",
			Attributes: buildMockAttributes("aws", region, req),
		}, nil
	}

	apiEndpoint, method := getAWSReadEndpoint(region, req.ResourceType, req.ResourceName)
	respAttrs := buildMockAttributes("aws", region, req)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, nil)
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("failed to create AWS Read HTTP request: %w", err)
		}
		signAWSRequest(ctx, cfg, httpReq, nil, region, req.ResourceType)

		/* #nosec G107 G704 */
		resp, err := common.HTTPClient.Do(httpReq)
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("AWS live Read API call error: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return common.ResourceResponse{}, fmt.Errorf("%w: AWS resource %s not found", common.ErrNotFound, req.ResourceName)
		}
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if resp.StatusCode >= 400 {
			return common.ResourceResponse{}, fmt.Errorf("AWS live API error (status %d) reading %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
		}
		_ = json.Unmarshal(bodyBytes, &respAttrs)
	}

	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "ACTIVE",
		Attributes: respAttrs,
	}, nil
}

func (a *AWSAdapter) UpdateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "ACTIVE",
			Attributes: buildMockAttributes("aws", region, req),
		}, nil
	}

	cfg, _ := loadAWSConfig(ctx, region, req.Attributes)
	apiEndpoint, method, payload := getAWSServiceEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("failed to create AWS Update HTTP request: %w", err)
		}
		if len(payload) > 0 && payload[0] == '{' {
			httpReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
		} else {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		signAWSRequest(ctx, cfg, httpReq, payload, region, req.ResourceType)

		/* #nosec G107 G704 */
		resp, err := common.HTTPClient.Do(httpReq)
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("AWS live Update API call error: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return common.ResourceResponse{}, fmt.Errorf("AWS live API error (status %d) updating %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
		}
	}
	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "ACTIVE",
		Attributes: buildMockAttributes("aws", region, req),
	}, nil
}

func (a *AWSAdapter) DeleteResource(ctx context.Context, req common.ResourceRequest) error {
	region := common.GetRegion(req.Region, "us-east-1")
	if common.IsRequestMockMode(req) {
		return nil
	}

	cfg, err := loadAWSConfig(ctx, region, req.Attributes)
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, err := s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
		if err != nil && !strings.Contains(strings.ToLower(err.Error()), "nosuchbucket") && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("AWS S3 DeleteBucket API error: %w", err)
		}
		return nil
	}

	apiEndpoint, method, payload := getAWSDeleteEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to create AWS Delete HTTP request: %w", err)
		}
		if len(payload) > 0 && payload[0] == '{' {
			httpReq.Header.Set("Content-Type", "application/x-amz-json-1.1")
		} else {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		signAWSRequest(ctx, cfg, httpReq, payload, region, req.ResourceType)

		/* #nosec G107 G704 */
		resp, err := common.HTTPClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("AWS live Delete API call error: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return fmt.Errorf("AWS live API error (status %d) deleting %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
		}
	}
	return nil
}
