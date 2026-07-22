package adapters

import (
	"context"
)

// ResourceRequest represents a unified resource payload for cloud operations
type ResourceRequest struct {
	ResourceName string
	ResourceType string
	ProviderType string
	Region       string
	Attributes   map[string]interface{}
}

// ResourceResponse represents a cloud operation result payload
type ResourceResponse struct {
	ID         string
	Status     string
	Attributes map[string]interface{}
}

// CloudAdapter defines the uniform interface for target cloud SDK drivers (AWS, GCP, Azure)
type CloudAdapter interface {
	CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	DeleteResource(ctx context.Context, req ResourceRequest) error
}

// AWSAdapter implements CloudAdapter for AWS SDK v2
type AWSAdapter struct{}

func (a *AWSAdapter) CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{
		ID:     "aws/" + req.Region + "/" + req.ResourceName,
		Status: "ACTIVE",
	}, nil
}
func (a *AWSAdapter) ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "ACTIVE"}, nil
}
func (a *AWSAdapter) UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "ACTIVE"}, nil
}
func (a *AWSAdapter) DeleteResource(ctx context.Context, req ResourceRequest) error {
	return nil
}

// GCPAdapter implements CloudAdapter for GCP Go SDK
type GCPAdapter struct{}

func (a *GCPAdapter) CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{
		ID:     "gcp/" + req.Region + "/" + req.ResourceName,
		Status: "RUNNING",
	}, nil
}
func (a *GCPAdapter) ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "RUNNING"}, nil
}
func (a *GCPAdapter) UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "RUNNING"}, nil
}
func (a *GCPAdapter) DeleteResource(ctx context.Context, req ResourceRequest) error {
	return nil
}

// AzureAdapter implements CloudAdapter for Azure ARM SDK
type AzureAdapter struct{}

func (a *AzureAdapter) CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{
		ID:     "azure/" + req.ResourceName,
		Status: "SUCCEEDED",
	}, nil
}
func (a *AzureAdapter) ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "SUCCEEDED"}, nil
}
func (a *AzureAdapter) UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error) {
	return ResourceResponse{ID: req.ResourceName, Status: "SUCCEEDED"}, nil
}
func (a *AzureAdapter) DeleteResource(ctx context.Context, req ResourceRequest) error {
	return nil
}
