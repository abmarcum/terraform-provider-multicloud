package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeepSchemaAttributeValidation(t *testing.T) {
	ctx := context.Background()

	// 1. Verify multicloud_storage_bucket sensitive & schema flags
	bucket := NewStorageBucketResource()
	bResp := &resource.SchemaResponse{}
	bucket.Schema(ctx, resource.SchemaRequest{}, bResp)

	if _, ok := bResp.Schema.Attributes["bucket_name"]; !ok {
		t.Errorf("multicloud_storage_bucket missing 'bucket_name'")
	}
	if attr, ok := bResp.Schema.Attributes["provider_type"]; !ok || !attr.IsRequired() {
		t.Errorf("multicloud_storage_bucket 'provider_type' should be required")
	}

	// 2. Verify sensitive attributes in db_instance
	db := NewDBInstanceResource()
	dbResp := &resource.SchemaResponse{}
	db.Schema(ctx, resource.SchemaRequest{}, dbResp)

	if passAttr, ok := dbResp.Schema.Attributes["password"]; ok {
		if !passAttr.IsSensitive() {
			t.Errorf("multicloud_db_instance 'password' attribute should be sensitive")
		}
	} else {
		t.Errorf("multicloud_db_instance missing 'password' attribute")
	}

	// 3. Verify sensitive attributes in secret
	sec := NewSecretResource()
	secResp := &resource.SchemaResponse{}
	sec.Schema(ctx, resource.SchemaRequest{}, secResp)

	if valAttr, ok := secResp.Schema.Attributes["secret_value"]; ok {
		if !valAttr.IsSensitive() {
			t.Errorf("multicloud_secret 'secret_value' attribute should be sensitive")
		}
	} else {
		t.Errorf("multicloud_secret missing 'secret_value' attribute")
	}

	// 4. Verify sensitive attributes in virtual_machine
	vm := NewVirtualMachineResource()
	vmResp := &resource.SchemaResponse{}
	vm.Schema(ctx, resource.SchemaRequest{}, vmResp)

	if sshAttr, ok := vmResp.Schema.Attributes["ssh_public_key"]; ok {
		if !sshAttr.IsSensitive() {
			t.Errorf("multicloud_virtual_machine 'ssh_public_key' attribute should be sensitive")
		}
	} else {
		t.Errorf("multicloud_virtual_machine missing 'ssh_public_key' attribute")
	}
}
