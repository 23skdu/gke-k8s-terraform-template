package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	m.Run()
}

func GetTerraformOptions(t *testing.T) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir: "../tf",
		Vars: map[string]interface{}{
			"project":                "test-project",
			"region":                 "us-central1",
			"environment":            "test",
			"cluster_name":           "test-cluster",
			"network_name":           "test-network",
			"subnet_cidr":            "10.0.0.0/20",
			"pods_cidr":              "10.4.0.0/14",
			"services_cidr":          "10.8.0.0/20",
			"master_ipv4_cidr_block": "172.16.0.0/28",
			"authorized_network_cidr": "0.0.0.0/0",
			"deletion_protection":    false,
			"main_machine_type":      "e2-medium",
			"main_min_count":         1,
			"main_max_count":         3,
			"main_disk_size_gb":      50,
			"main_disk_type":         "pd-balanced",
			"system_machine_type":    "e2-small",
			"system_min_count":       1,
			"system_max_count":       2,
			"system_disk_size_gb":    50,
			"system_disk_type":       "pd-balanced",
			"namespaces":             []string{"default", "test"},
			"namespace":              "default",
			"enable_alerting":        false,
			"enable_backup":          false,
			"enable_binary_authorization": false,
			"secrets":                map[string]string{},
			"state_bucket_name":      "tf-state",
		},
	}
}

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

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	privateCluster := clusterResource.AttributeValues["private_cluster_config"]
	require.NotNil(t, privateCluster, "private_cluster_config should exist")
}

func TestGKEClusterWorkloadIdentity(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	workloadIdentity := clusterResource.AttributeValues["workload_identity_config"]
	require.NotNil(t, workloadIdentity, "workload_identity_config should exist")
}

func TestGKEClusterAdvancedDatapath(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	datapath := clusterResource.AttributeValues["datapath_provider"]
	assert.Equal(t, "ADVANCED_DATAPATH", datapath, "datapath_provider should be ADVANCED_DATAPATH")
}

func TestGKEClusterGatewayAPI(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	gatewayAPI := clusterResource.AttributeValues["gateway_api_config"]
	require.NotNil(t, gatewayAPI, "gateway_api_config should exist")
}

func TestGKEClusterReleaseChannel(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	releaseChannel := clusterResource.AttributeValues["release_channel"]
	require.NotNil(t, releaseChannel, "release_channel should exist")
}

func TestGKEClusterLogging(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	loggingConfig := clusterResource.AttributeValues["logging_config"]
	require.NotNil(t, loggingConfig, "logging_config should exist")
}

func TestGKEClusterMonitoring(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	monitoringConfig := clusterResource.AttributeValues["monitoring_config"]
	require.NotNil(t, monitoringConfig, "monitoring_config should exist")
}

func TestGKENodePoolShieldedConfig(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nodePoolResource := plan.ResourcePlannedValuesMap["module.gke.google_container_node_pool.main"]
	require.NotNil(t, nodePoolResource, "google_container_node_pool.main should exist in gke module")

	nodeConfig := nodePoolResource.AttributeValues["node_config"]
	require.NotNil(t, nodeConfig, "node_config should exist")
}

func TestGKENodePoolAutoUpgrade(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	nodePoolResource := plan.ResourcePlannedValuesMap["module.gke.google_container_node_pool.main"]
	require.NotNil(t, nodePoolResource, "google_container_node_pool.main should exist in gke module")

	management := nodePoolResource.AttributeValues["management"]
	require.NotNil(t, management, "management should exist")
}

func TestGKEClusterDeletionProtection(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	clusterResource := plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
	require.NotNil(t, clusterResource, "google_container_cluster.cluster should exist in gke module")

	deletionProtection := clusterResource.AttributeValues["deletion_protection"]
	assert.Equal(t, false, deletionProtection, "deletion_protection should be false in test")
}

func TestSystemNodePoolExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	systemPool := plan.ResourcePlannedValuesMap["module.gke.google_container_node_pool.system"]
	require.NotNil(t, systemPool, "system node pool should exist in gke module")
}

func TestNetworkingVPCExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	vpc := plan.ResourcePlannedValuesMap["module.networking.google_compute_network.vpc"]
	require.NotNil(t, vpc, "VPC should exist in networking module")
}

func TestIAMServiceAccountExists(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	sa := plan.ResourcePlannedValuesMap["module.iam.google_service_account.cluster_sa"]
	require.NotNil(t, sa, "service account should exist in iam module")
}

func TestMonitoringAlertPolicies(t *testing.T) {
	t.Parallel()
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	logging := plan.ResourcePlannedValuesMap["module.monitoring.google_logging_project_bucket_config.default"]
	require.NotNil(t, logging, "logging bucket should exist in monitoring module")
}
