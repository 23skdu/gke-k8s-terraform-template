package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestLoggingBucketExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	loggingResource := plan.ResourcePlannedValuesMap["google_logging_project_bucket_config.default"]
	require.NotNil(t, loggingResource, "google_logging_project_bucket_config.default should exist")
}

func TestLoggingBucketRetention(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	loggingResource := plan.ResourcePlannedValuesMap["google_logging_project_bucket_config.default"]
	require.NotNil(t, loggingResource, "google_logging_project_bucket_config.default should exist")

	retentionDays := loggingResource.AttributeValues["retention_days"]
	require.NotNil(t, retentionDays, "retention_days should exist")
}
