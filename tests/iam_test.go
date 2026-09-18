package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAccountExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	saResource := plan.ResourcePlannedValuesMap["google_service_account.cluster_sa"]
	require.NotNil(t, saResource, "google_service_account.cluster_sa should exist")
}

func TestServiceAccountEmail(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	email := terraform.Output(t, terraformOptions, "service_account_email")
	require.NotEmpty(t, email, "service_account_email output should not be empty")
	assert.Contains(t, email, "@", "email should contain @")
}

func TestIAMRolesAssigned(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	iamRolesResource := plan.ResourcePlannedValuesMap["google_project_iam_member.cluster_sa_roles[\"roles/logging.logWriter\"]"]
	require.NotNil(t, iamRolesResource, "IAM role for logging.logWriter should exist")

	iamRolesResource2 := plan.ResourcePlannedValuesMap["google_project_iam_member.cluster_sa_roles[\"roles/monitoring.metricWriter\"]"]
	require.NotNil(t, iamRolesResource2, "IAM role for monitoring.metricWriter should exist")

	iamRolesResource3 := plan.ResourcePlannedValuesMap["google_project_iam_member.cluster_sa_roles[\"roles/monitoring.viewer\"]"]
	require.NotNil(t, iamRolesResource3, "IAM role for monitoring.viewer should exist")

	iamRolesResource4 := plan.ResourcePlannedValuesMap["google_project_iam_member.cluster_sa_roles[\"roles/stackdriver.resourceMetadata.writer\"]"]
	require.NotNil(t, iamRolesResource4, "IAM role for stackdriver.resourceMetadata.writer should exist")
}

func TestWorkloadIdentityBinding(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	wiResource := plan.ResourcePlannedValuesMap["google_service_account_iam_member.workload_identity"]
	require.NotNil(t, wiResource, "workload_identity IAM binding should exist")
}
