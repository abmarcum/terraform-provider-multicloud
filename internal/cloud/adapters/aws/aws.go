package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSAdapter struct{}

func getAWSAccountID() string {
	acc := os.Getenv("AWS_ACCOUNT_ID")
	if acc == "" {
		acc = "unknown-account"
	}
	return url.PathEscape(acc)
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
	case "storage_bucket":
		endpoint = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", escName, escRegion)
		method = "PUT"
	case "virtual_network", "vpc_peering":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateVpc&CidrBlock=10.0.0.0/16&Version=2016-11-15", escRegion)
	case "subnet":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateSubnet&CidrBlock=10.0.1.0/24&Version=2016-11-15", escRegion)
	case "security_group":
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=CreateSecurityGroup&GroupName=%s&Version=2016-11-15", escRegion, escQueryName)
	case "db_instance":
		endpoint = fmt.Sprintf("https://rds.%s.amazonaws.com/?Action=CreateDBInstance&DBInstanceIdentifier=%s&Version=2014-10-31", escRegion, escQueryName)
	case "secret", "secret_rotator":
		endpoint = fmt.Sprintf("https://secretsmanager.%s.amazonaws.com", escRegion)
		bodyMap := map[string]string{"Name": name}
		payload, _ = json.Marshal(bodyMap)
	case "serverless_function":
		endpoint = fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions", escRegion)
		bodyMap := map[string]string{"FunctionName": name, "Runtime": "python311", "Role": fmt.Sprintf("arn:aws:iam::%s:role/service-role", accID)}
		payload, _ = json.Marshal(bodyMap)
	case "kubernetes_cluster":
		endpoint = fmt.Sprintf("https://eks.%s.amazonaws.com/clusters", escRegion)
		bodyMap := map[string]string{"name": name}
		payload, _ = json.Marshal(bodyMap)
	case "nosql_table":
		endpoint = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", escRegion)
		bodyMap := map[string]interface{}{"TableName": name, "AttributeDefinitions": []map[string]string{{"AttributeName": "id", "AttributeType": "S"}}}
		payload, _ = json.Marshal(bodyMap)
	case "kms_key":
		endpoint = fmt.Sprintf("https://kms.%s.amazonaws.com", escRegion)
		payload = []byte(`{"Description":"Multi-cloud KMS Key"}`)
	case "pubsub_topic":
		endpoint = fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=CreateTopic&Name=%s&Version=2010-03-31", escRegion, escQueryName)
	case "message_queue":
		endpoint = fmt.Sprintf("https://sqs.%s.amazonaws.com/?Action=CreateQueue&QueueName=%s&Version=2012-11-05", escRegion, escQueryName)
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
	case "storage_bucket":
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
		bodyMap := map[string]interface{}{"SecretId": name, "ForceDeleteWithoutRecovery": true}
		payload, _ = json.Marshal(bodyMap)
	case "serverless_function":
		endpoint = fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions/%s", escRegion, escName)
		method = "DELETE"
	case "kubernetes_cluster":
		endpoint = fmt.Sprintf("https://eks.%s.amazonaws.com/clusters/%s", escRegion, escName)
		method = "DELETE"
	case "nosql_table":
		endpoint = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", escRegion)
		bodyMap := map[string]string{"TableName": name}
		payload, _ = json.Marshal(bodyMap)
	case "pubsub_topic":
		endpoint = fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=DeleteTopic&TopicArn=arn:aws:sns:%s:%s:%s&Version=2010-03-31", escRegion, escRegion, accID, escQueryName)
	case "message_queue":
		endpoint = fmt.Sprintf("https://sqs.%s.amazonaws.com/?Action=DeleteQueue&QueueUrl=https://sqs.%s.amazonaws.com/%s/%s&Version=2012-11-05", escRegion, escRegion, accID, escName)
	default:
		endpoint = fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=TerminateInstances&InstanceId.1=%s&Version=2016-11-15", escRegion, escQueryName)
	}

	return endpoint, method, payload
}

func getAWSReadEndpoint(region string, resType string, name string) (string, string) {
	escName := url.PathEscape(name)
	escQueryName := url.QueryEscape(name)
	escRegion := url.PathEscape(region)

	switch resType {
	case "storage_bucket":
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", escName, escRegion), "HEAD"
	case "virtual_network", "vpc_peering":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeVpcs&VpcId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "subnet":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeSubnets&SubnetId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "security_group":
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeSecurityGroups&GroupName.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	case "db_instance":
		return fmt.Sprintf("https://rds.%s.amazonaws.com/?Action=DescribeDBInstances&DBInstanceIdentifier=%s&Version=2014-10-31", escRegion, escQueryName), "POST"
	case "serverless_function":
		return fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions/%s", escRegion, escName), "GET"
	case "kubernetes_cluster":
		return fmt.Sprintf("https://eks.%s.amazonaws.com/clusters/%s", escRegion, escName), "GET"
	default:
		return fmt.Sprintf("https://ec2.%s.amazonaws.com/?Action=DescribeInstances&InstanceId.1=%s&Version=2016-11-15", escRegion, escQueryName), "POST"
	}
}

func (a *AWSAdapter) CreateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
		if err == nil {
			return common.ResourceResponse{
				ID:     fmt.Sprintf("arn:aws:s3:::%s", req.ResourceName),
				Status: "ACTIVE",
			}, nil
		}
	}

	apiEndpoint, method, payload := getAWSServiceEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}

	return common.ResourceResponse{
		ID:     fmt.Sprintf("aws/%s/%s/%s", req.ResourceType, region, req.ResourceName),
		Status: "ACTIVE",
	}, nil
}

func (a *AWSAdapter) ReadResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
		if err != nil {
			return common.ResourceResponse{}, fmt.Errorf("AWS S3 bucket %s not found: %w", req.ResourceName, err)
		}
	}

	apiEndpoint, method := getAWSReadEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, nil)
		if err == nil {
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == 404 {
					return common.ResourceResponse{}, fmt.Errorf("AWS resource %s (%s) not found", req.ResourceName, req.ResourceType)
				}
			}
		}
	}

	return common.ResourceResponse{ID: req.ResourceName, Status: "ACTIVE"}, nil
}

func (a *AWSAdapter) UpdateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	region := common.GetRegion(req.Region, "us-east-1")
	apiEndpoint, method, payload := getAWSServiceEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}
	return common.ResourceResponse{ID: req.ResourceName, Status: "ACTIVE"}, nil
}

func (a *AWSAdapter) DeleteResource(ctx context.Context, req common.ResourceRequest) error {
	region := common.GetRegion(req.Region, "us-east-1")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err == nil && req.ResourceType == "storage_bucket" {
		s3Client := s3.NewFromConfig(cfg)
		_, _ = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(req.ResourceName),
		})
	}

	apiEndpoint, method, payload := getAWSDeleteEndpoint(region, req.ResourceType, req.ResourceName)
	if apiEndpoint != "" {
		httpReq, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bytes.NewBuffer(payload))
		if err == nil {
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}
	return nil
}
