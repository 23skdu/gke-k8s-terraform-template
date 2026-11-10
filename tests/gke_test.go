package test

import (
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Integration tests (apply real GCP resources). Skipped with -short.
// ---------------------------------------------------------------------------

func TestGKEClusterExists(t *testing.T) {
	t.Parallel()
	integrationSkip(t)
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	clusterName := terraform.Output(t, terraformOptions, "cluster_name")
	assert.Equal(t, "test-cluster-test", clusterName,
		"cluster_name output must equal '<cluster_name>-<environment>'")
}

func TestGKEClusterEndpoint(t *testing.T) {
	t.Parallel()
	integrationSkip(t)
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	endpoint := terraform.Output(t, terraformOptions, "cluster_endpoint")
	// google_container_cluster.endpoint is a bare IP address, not a URL.
	assert.Regexp(t, `^\d{1,3}(\.\d{1,3}){3}$`, endpoint,
		"cluster_endpoint must be a bare IPv4 address")
}

func TestGKEClusterLocation(t *testing.T) {
	t.Parallel()
	integrationSkip(t)
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	location := terraform.Output(t, terraformOptions, "cluster_location")
	assert.Equal(t, "us-central1", location, "cluster_location must match the configured region")
}

func TestGKEClusterCA(t *testing.T) {
	t.Parallel()
	integrationSkip(t)
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	ca := terraform.Output(t, terraformOptions, "cluster_ca_certificate")
	der, err := base64.StdEncoding.DecodeString(ca)
	require.NoError(t, err, "cluster_ca_certificate must be base64")
	require.NotEmpty(t, der, "cluster_ca_certificate must decode to a certificate")
	assert.Equal(t, byte(0x30), der[0], "cluster_ca_certificate must be a DER SEQUENCE")

	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err, "cluster_ca_certificate must parse as an X.509 certificate")
	assert.True(t, cert.IsCA, "cluster CA certificate must be a CA certificate")
}

func TestGKEClusterID(t *testing.T) {
	t.Parallel()
	integrationSkip(t)
	terraformOptions := GetTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	id := terraform.Output(t, terraformOptions, "cluster_id")
	assert.Contains(t, id, "projects/test-project", "cluster_id must reference the configured project")
	assert.Contains(t, id, "test-cluster-test", "cluster_id must reference the cluster resource name")
}

// ---------------------------------------------------------------------------
// Plan-only tests: one shared, hermetic plan for the whole suite.
// ---------------------------------------------------------------------------

func TestGKEClusterPrivateNodes(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	private := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "private_cluster_config")
	assert.Equal(t, true, private["enable_private_nodes"], "private nodes must be enabled")
	assert.Equal(t, "172.16.0.0/28", private["master_ipv4_cidr_block"], "master CIDR must match input")
}

func TestGKEClusterWorkloadIdentity(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	wic := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "workload_identity_config")
	assert.Equal(t, "test-project.svc.id.goog", wic["workload_pool"],
		"workload identity pool must be '<project>.svc.id.goog'")
}

func TestGKEClusterAdvancedDatapath(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	datapath := requireAttr(t, plan, "module.gke.google_container_cluster.cluster", "datapath_provider")
	assert.Equal(t, "ADVANCED_DATAPATH", datapath, "datapath_provider must be ADVANCED_DATAPATH")
}

func TestGKEClusterGatewayAPI(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	gateway := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "gateway_api_config")
	assert.Equal(t, "CHANNEL_STANDARD", gateway["channel"], "gateway API channel must be CHANNEL_STANDARD")
}

func TestGKEClusterReleaseChannel(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	channel := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "release_channel")
	assert.Equal(t, "REGULAR", channel["channel"], "release channel must be REGULAR")
}

func TestGKEClusterLogging(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	logging := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "logging_config")
	components, ok := logging["enable_components"].([]interface{})
	require.True(t, ok, "logging_config.enable_components must be a list")
	require.Len(t, components, 2, "exactly two logging components must be enabled")
	assert.True(t, containsAll(components, "SYSTEM_COMPONENTS", "WORKLOADS"),
		"logging must cover system and workload components, got %v", components)
}

func TestGKEClusterMonitoring(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	monitoring := requireBlock(t, plan, "module.gke.google_container_cluster.cluster", "monitoring_config")
	components, ok := monitoring["enable_components"].([]interface{})
	require.True(t, ok, "monitoring_config.enable_components must be a list")
	assert.True(t, containsAll(components, "SYSTEM_COMPONENTS"),
		"system component monitoring must be enabled, got %v", components)

	prom, ok := monitoring["managed_prometheus"].([]interface{})
	require.True(t, ok, "managed_prometheus must be present")
	require.Len(t, prom, 1, "managed_prometheus must contain one block")
	promBlock := prom[0].(map[string]interface{})
	assert.Equal(t, true, promBlock["enabled"], "managed Prometheus must be enabled")
}

func TestGKENodePoolShieldedConfig(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	nodeConfig := requireBlock(t, plan, "module.gke.google_container_node_pool.main", "node_config")
	shielded := nodeConfig["shielded_instance_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, true, shielded["enable_secure_boot"], "secure boot must be enabled on main pool")
	assert.Equal(t, true, shielded["enable_integrity_monitoring"], "integrity monitoring must be enabled")

	assert.Equal(t, "e2-medium", nodeConfig["machine_type"], "main pool machine type must match input")
	assert.Equal(t, "pd-balanced", nodeConfig["disk_type"], "main pool disk type must match input")
	assert.Equal(t, float64(50), requireFloat(t, nodeConfig["disk_size_gb"]), "main pool disk size must match input")

	metadata := nodeConfig["workload_metadata_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "GKE_METADATA", metadata["mode"], "GKE metadata server must be enabled")
}

func TestGKENodePoolAutoUpgrade(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	management := requireBlock(t, plan, "module.gke.google_container_node_pool.main", "management")
	assert.Equal(t, true, management["auto_upgrade"], "auto upgrade must be enabled")
	assert.Equal(t, true, management["auto_repair"], "auto repair must be enabled")

	upgrade := requireBlock(t, plan, "module.gke.google_container_node_pool.main", "upgrade_settings")
	assert.Equal(t, float64(1), requireFloat(t, upgrade["max_surge"]), "max_surge must be 1")
	assert.Equal(t, float64(0), requireFloat(t, upgrade["max_unavailable"]), "max_unavailable must be 0")
}

func TestGKEClusterDeletionProtection(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	deletionProtection := requireAttr(t, plan, "module.gke.google_container_cluster.cluster", "deletion_protection")
	assert.Equal(t, false, deletionProtection, "deletion_protection must be false in test environment")
}

func TestSystemNodePoolExists(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	res := plan.ResourcePlannedValuesMap["module.gke.google_container_node_pool.system"]
	require.NotNil(t, res, "system node pool must exist in gke module")

	nodeConfig := res.AttributeValues["node_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "e2-small", nodeConfig["machine_type"], "system pool machine type must match input")

	taints := nodeConfig["taint"].([]interface{})
	require.Len(t, taints, 1, "system pool must carry exactly one taint")
	taint := taints[0].(map[string]interface{})
	assert.Equal(t, "node-pool", taint["key"], "taint key must be node-pool")
	assert.Equal(t, "system", taint["value"], "taint value must be system")
	assert.Equal(t, "PREFER_NO_SCHEDULE", taint["effect"], "taint effect must be PREFER_NO_SCHEDULE")

	autoscaling := requireBlock(t, plan, "module.gke.google_container_node_pool.system", "autoscaling")
	assert.Equal(t, float64(1), requireFloat(t, autoscaling["min_node_count"]), "system pool min count must match input")
	assert.Equal(t, float64(2), requireFloat(t, autoscaling["max_node_count"]), "system pool max count must match input")
}

func TestNetworkingVPCExists(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	vpc := requireAttr(t, plan, "module.networking.google_compute_network.vpc", "name")
	assert.Equal(t, "test-network-test", vpc, "VPC name must be '<network_name>-<environment>'")
}

func TestIAMServiceAccountExists(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	sa := requireAttr(t, plan, "module.iam.google_service_account.cluster_sa", "account_id")
	assert.Equal(t, "test-cluster-test-sa", sa, "service account id must be '<cluster_name>-<environment>-sa'")
}

func TestMonitoringAlertPolicies(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	bucket := requireAttr(t, plan, "module.monitoring.google_logging_project_bucket_config.default", "bucket_id")
	assert.Equal(t, "_Default", bucket, "log bucket must be the project _Default bucket")

	retention := requireAttr(t, plan, "module.monitoring.google_logging_project_bucket_config.default", "retention_days")
	assert.Equal(t, float64(30), requireFloat(t, retention), "log bucket retention must be 30 days")
}
