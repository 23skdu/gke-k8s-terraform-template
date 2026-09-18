package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGKEClusterExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	clusterName := terraform.Output(t, terraformOptions, "cluster_name")
	require.NotEmpty(t, clusterName, "cluster_name output should not be empty")
}

func TestGKEClusterEndpoint(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	endpoint := terraform.Output(t, terraformOptions, "cluster_endpoint")
	require.NotEmpty(t, endpoint, "cluster_endpoint output should not be empty")
	assert.Contains(t, endpoint, "https://", "endpoint should use HTTPS")
}

func TestGKEClusterLocation(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	location := terraform.Output(t, terraformOptions, "cluster_location")
	require.NotEmpty(t, location, "cluster_location output should not be empty")
}

func TestGKEClusterCA(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	ca := terraform.Output(t, terraformOptions, "cluster_ca_certificate")
	require.NotEmpty(t, ca, "cluster_ca_certificate output should not be empty")
}

func TestGKEClusterID(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	id := terraform.Output(t, terraformOptions, "cluster_id")
	require.NotEmpty(t, id, "cluster_id output should not be empty")
}

func TestGKEClusterPrivateNodes(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	privateCluster := clusterResource.AttributeValues["private_cluster_config"]
	require.NotNil(t, privateCluster, "private_cluster_config should exist")
}

func TestGKEClusterWorkloadIdentity(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	workloadIdentity := clusterResource.AttributeValues["workload_identity_config"]
	require.NotNil(t, workloadIdentity, "workload_identity_config should exist")
}

func TestGKEClusterAdvancedDatapath(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	datapath := clusterResource.AttributeValues["datapath_provider"]
	assert.Equal(t, "ADVANCED_DATAPATH", datapath, "datapath_provider should be ADVANCED_DATAPATH")
}

func TestGKEClusterGatewayAPI(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	gatewayAPI := clusterResource.AttributeValues["gateway_api_config"]
	require.NotNil(t, gatewayAPI, "gateway_api_config should exist")
}

func TestGKEClusterReleaseChannel(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	releaseChannel := clusterResource.AttributeValues["release_channel"]
	require.NotNil(t, releaseChannel, "release_channel should exist")
}

func TestGKEClusterLogging(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	loggingConfig := clusterResource.AttributeValues["logging_config"]
	require.NotNil(t, loggingConfig, "logging_config should exist")
}

func TestGKEClusterMonitoring(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	monitoringConfig := clusterResource.AttributeValues["monitoring_config"]
	require.NotNil(t, monitoringConfig, "monitoring_config should exist")
}

func TestGKENodePoolShieldedConfig(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nodePoolResource := plan.ResourcePlannedValuesMap["google_container_node_pool.main"]
	require.NotNil(t, nodePoolResource, "google_container_node_pool.main should exist")

	nodeConfig := nodePoolResource.AttributeValues["node_config"]
	require.NotNil(t, nodeConfig, "node_config should exist")
}

func TestGKENodePoolAutoUpgrade(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nodePoolResource := plan.ResourcePlannedValuesMap["google_container_node_pool.main"]
	require.NotNil(t, nodePoolResource, "google_container_node_pool.main should exist")

	management := nodePoolResource.AttributeValues["management"]
	require.NotNil(t, management, "management should exist")
}

func TestGKEClusterDeletionProtection(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist")

	deletionProtection := clusterResource.AttributeValues["deletion_protection"]
	assert.Equal(t, false, deletionProtection, "deletion_protection should be false")
}
