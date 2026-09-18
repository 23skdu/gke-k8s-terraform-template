package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCSBucketExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	bucketResource := plan.ResourcePlannedValuesMap["google_storage_bucket.tf_state"]
	require.NotNil(t, bucketResource, "google_storage_bucket.tf_state should exist")
}

func TestGCSBucketVersioning(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	bucketResource := plan.ResourcePlannedValuesMap["google_storage_bucket.tf_state"]
	require.NotNil(t, bucketResource, "google_storage_bucket.tf_state should exist")

	versioning := bucketResource.AttributeValues["versioning"]
	require.NotNil(t, versioning, "versioning should exist")
}

func TestGCSBucketUniformAccess(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	bucketResource := plan.ResourcePlannedValuesMap["google_storage_bucket.tf_state"]
	require.NotNil(t, bucketResource, "google_storage_bucket.tf_state should exist")

	uniformAccess := bucketResource.AttributeValues["uniform_bucket_level_access"]
	assert.Equal(t, true, uniformAccess, "uniform_bucket_level_access should be true")
}

func TestGCSBucketPublicAccessPrevention(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	bucketResource := plan.ResourcePlannedValuesMap["google_storage_bucket.tf_state"]
	require.NotNil(t, bucketResource, "google_storage_bucket.tf_state should exist")

	publicAccess := bucketResource.AttributeValues["public_access_prevention"]
	assert.Equal(t, "enforced", publicAccess, "public_access_prevention should be enforced")
}

func TestGCSBucketLifecycleRules(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	bucketResource := plan.ResourcePlannedValuesMap["google_storage_bucket.tf_state"]
	require.NotNil(t, bucketResource, "google_storage_bucket.tf_state should exist")

	lifecycleRules := bucketResource.AttributeValues["lifecycle_rule"]
	require.NotNil(t, lifecycleRules, "lifecycle_rule should exist")
}
