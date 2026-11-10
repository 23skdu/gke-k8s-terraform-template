package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Plan-only coverage for modules that previously had no tests at all:
// networking (subnet, firewalls, NAT), kubernetes, iam (Workload Identity),
// backup, security (Binary Authorization), secrets, statebucket and the
// monitoring alert policies.

func TestNetworkingSubnetPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	subnet := requireAttr(t, plan, "module.networking.google_compute_subnetwork.gke", "ip_cidr_range")
	assert.Equal(t, "10.0.0.0/20", subnet, "subnet CIDR must match input")

	privateAccess := requireAttr(t, plan, "module.networking.google_compute_subnetwork.gke", "private_ip_google_access")
	assert.Equal(t, true, privateAccess, "subnet must have private Google access enabled")

	ranges := plan.ResourcePlannedValuesMap["module.networking.google_compute_subnetwork.gke"].AttributeValues["secondary_ip_range"].([]interface{})
	require.Len(t, ranges, 2, "subnet must define pods and services secondary ranges")
	byName := map[string]string{}
	for _, r := range ranges {
		block := r.(map[string]interface{})
		byName[block["range_name"].(string)] = block["ip_cidr_range"].(string)
	}
	assert.Equal(t, "10.4.0.0/14", byName["pods"], "pods secondary range must match input")
	assert.Equal(t, "10.8.0.0/20", byName["services"], "services secondary range must match input")
}

func TestNetworkingFirewallPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	for _, address := range []string{
		"module.networking.google_compute_firewall.allow_internal",
		"module.networking.google_compute_firewall.allow_health_checks",
		"module.networking.google_compute_firewall.allow_master_to_nodes",
		"module.networking.google_compute_firewall.allow_nodes_to_master",
	} {
		require.NotNil(t, plan.ResourcePlannedValuesMap[address], "%s must exist", address)
	}

	internalRanges := requireAttr(t, plan, "module.networking.google_compute_firewall.allow_internal", "source_ranges").([]interface{})
	require.Len(t, internalRanges, 3, "internal firewall must allow subnet, pods and services CIDRs")
	assert.True(t, containsAll(internalRanges, "10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20"),
		"internal firewall source ranges must match the three cluster CIDRs, got %v", internalRanges)

	healthRanges := requireAttr(t, plan, "module.networking.google_compute_firewall.allow_health_checks", "source_ranges").([]interface{})
	require.Len(t, healthRanges, 2, "health-check firewall must allow Google health check ranges")
	assert.True(t, containsAll(healthRanges, "130.211.0.0/22", "35.191.0.0/16"),
		"health-check firewall must allow 130.211.0.0/22 and 35.191.0.0/16, got %v", healthRanges)

	masterAllow := requireBlock(t, plan, "module.networking.google_compute_firewall.allow_master_to_nodes", "allow")
	masterPorts := masterAllow["ports"].([]interface{})
	assert.True(t, containsAll(masterPorts, "443", "10250"),
		"master-to-nodes firewall must allow 443 and 10250, got %v", masterPorts)

	masterRanges := requireAttr(t, plan, "module.networking.google_compute_firewall.allow_master_to_nodes", "source_ranges").([]interface{})
	assert.True(t, containsAll(masterRanges, "172.16.0.0/28"),
		"master-to-nodes firewall must be restricted to the master CIDR, got %v", masterRanges)
}

func TestNetworkingNATPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	nat := plan.ResourcePlannedValuesMap["module.networking.google_compute_router_nat.nat[0]"]
	require.NotNil(t, nat, "Cloud NAT must be planned (enable_cloud_nat defaults to true)")
	assert.Equal(t, "AUTO_ONLY", nat.AttributeValues["nat_ip_allocate_option"], "NAT IPs must be auto allocated")
	assert.Equal(t, "ALL_SUBNETWORKS_ALL_IP_RANGES", nat.AttributeValues["source_subnetwork_ip_ranges_to_nat"],
		"NAT must cover all subnetwork IP ranges")

	logConfig := nat.AttributeValues["log_config"].([]interface{})
	require.Len(t, logConfig, 1, "NAT logging must be configured")
	assert.Equal(t, true, logConfig[0].(map[string]interface{})["enable"], "NAT logging must be enabled")

	require.NotNil(t, plan.ResourcePlannedValuesMap["module.networking.google_compute_router.router[0]"],
		"Cloud Router must be planned")
}

func TestKubernetesNamespacesPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	for _, ns := range []string{"default", "test"} {
		address := `module.kubernetes.kubernetes_namespace_v1.namespaces["` + ns + `"]`
		res := plan.ResourcePlannedValuesMap[address]
		require.NotNil(t, res, "namespace %q must be planned", ns)

		metadata := res.AttributeValues["metadata"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, ns, metadata["name"], "namespace metadata.name must match")

		labels := metadata["labels"].(map[string]interface{})
		assert.Equal(t, "test", labels["environment"], "namespace must be labelled with the environment")
	}
}

func TestIAMWorkloadIdentityBindingPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	binding := plan.ResourcePlannedValuesMap["module.iam.google_service_account_iam_member.workload_identity"]
	require.NotNil(t, binding, "Workload Identity binding must be planned")
	assert.Equal(t, "roles/iam.workloadIdentityUser", binding.AttributeValues["role"],
		"Workload Identity binding must grant roles/iam.workloadIdentityUser")
}

func TestBackupPlanEnabled(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	backup := plan.ResourcePlannedValuesMap["module.backup.google_gke_backup_backup_plan.cluster_backup[0]"]
	require.NotNil(t, backup, "backup plan must be planned when enable_backup=true")
	assert.Equal(t, "test-cluster-test-backup", backup.AttributeValues["name"],
		"backup plan name must be '<cluster_name>-<environment>-backup'")

	retention := backup.AttributeValues["retention_policy"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, float64(14), requireFloat(t, retention["backup_retain_days"]), "retention must match backup_retention_days")
	assert.Equal(t, float64(14), requireFloat(t, retention["backup_delete_lock_days"]), "delete lock must match backup_retention_days")

	backupConfig := backup.AttributeValues["backup_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, true, backupConfig["include_volume_data"], "backups must include volume data")
	assert.Equal(t, true, backupConfig["include_secrets"], "backups must include secrets")
	assert.Equal(t, true, backupConfig["all_namespaces"], "backups must cover all namespaces")
}

func TestBinaryAuthorizationEnabled(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	policy := plan.ResourcePlannedValuesMap["module.security.google_binary_authorization_policy.policy[0]"]
	require.NotNil(t, policy, "Binary Authorization policy must be planned when enabled")
	assert.Equal(t, "ENABLE", policy.AttributeValues["global_policy_evaluation_mode"],
		"global policy evaluation must be enabled")

	whitelist := policy.AttributeValues["admission_whitelist_patterns"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "gcr.io/test-project/*", whitelist["name_pattern"],
		"images from the project's own registry must be whitelisted")

	rule := policy.AttributeValues["default_admission_rule"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "REQUIRE_ATTESTATION", rule["evaluation_mode"],
		"default admission rule must require attestation")
	assert.Equal(t, "ENFORCED_BLOCK_AND_AUDIT_LOG", rule["enforcement_mode"],
		"default admission rule must block and audit")

	attestor := plan.ResourcePlannedValuesMap["module.security.google_binary_authorization_attestor.attestor[0]"]
	require.NotNil(t, attestor, "Binary Authorization attestor must be planned")
	assert.Equal(t, "test-attestor", attestor.AttributeValues["name"], "attestor name must be '<environment>-attestor'")

	// The attestor must reference a Container Analysis note; the configuration
	// previously used a non-existent google_binary_authorization_note resource.
	note := plan.ResourcePlannedValuesMap["module.security.google_container_analysis_note.note[0]"]
	require.NotNil(t, note, "Container Analysis note must be planned")
	assert.Equal(t, "test-note", note.AttributeValues["name"], "note name must be '<environment>-note'")
}

func TestSecretsManagerPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	secret := plan.ResourcePlannedValuesMap[`module.secrets.google_secret_manager_secret.secrets["db-password"]`]
	require.NotNil(t, secret, "secret must be planned when present in the secrets map")
	assert.Equal(t, "db-password-test", secret.AttributeValues["secret_id"],
		"secret id must be '<key>-<environment>'")

	member := plan.ResourcePlannedValuesMap[`module.secrets.google_secret_manager_secret_iam_member.secret_access["db-password"]`]
	require.NotNil(t, member, "secret accessor IAM binding must be planned")
	assert.Equal(t, "roles/secretmanager.secretAccessor", member.AttributeValues["role"],
		"cluster service account must be granted secretAccessor")
}

func TestMonitoringAlertPolicyThresholds(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	channel := plan.ResourcePlannedValuesMap["module.monitoring.google_monitoring_notification_channel.email[0]"]
	require.NotNil(t, channel, "email notification channel must be planned when notification_email is set")
	assert.Equal(t, "email", channel.AttributeValues["type"], "notification channel must be of type email")
	channelLabels := channel.AttributeValues["labels"].(map[string]interface{})
	assert.Equal(t, "alerts@example.com", channelLabels["email_address"], "channel must target notification_email")

	cases := []struct {
		address    string
		display    string
		comparison string
		threshold  float64
	}{
		{"module.monitoring.google_monitoring_alert_policy.cluster_health[0]", "GKE Cluster Health (test)", "COMPARISON_LT", 1},
		{"module.monitoring.google_monitoring_alert_policy.high_cpu[0]", "GKE High CPU Usage (test)", "COMPARISON_GT", 0.8},
		{"module.monitoring.google_monitoring_alert_policy.high_memory[0]", "GKE High Memory Usage (test)", "COMPARISON_GT", 0.85},
		{"module.monitoring.google_monitoring_alert_policy.auto_repair_failure[0]", "GKE Node Auto-Repair Failure (test)", "COMPARISON_GT", 0},
	}
	for _, tc := range cases {
		policy := plan.ResourcePlannedValuesMap[tc.address]
		require.NotNil(t, policy, "%s must be planned when enable_alerting=true", tc.address)
		assert.Equal(t, tc.display, policy.AttributeValues["display_name"], "alert display name must match")

		condition := policy.AttributeValues["conditions"].([]interface{})[0].(map[string]interface{})
		threshold := condition["condition_threshold"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, tc.comparison, threshold["comparison"], "alert comparison must match")
		assert.Equal(t, tc.threshold, requireFloat(t, threshold["threshold_value"]), "alert threshold must match")
	}

	// The cluster-health alert must query node readiness with the k8s_node
	// resource type; the previous k8s_cluster filter matched no time series.
	health := plan.ResourcePlannedValuesMap["module.monitoring.google_monitoring_alert_policy.cluster_health[0]"]
	healthThreshold := health.AttributeValues["conditions"].([]interface{})[0].(map[string]interface{})["condition_threshold"].([]interface{})[0].(map[string]interface{})
	filter := requireString(t, healthThreshold["filter"])
	assert.Contains(t, filter, `resource.type = "k8s_node"`, "cluster health filter must use the k8s_node resource type")
	assert.Contains(t, filter, "kubernetes.io/node/condition/ready", "cluster health filter must query node readiness")
}

func TestStateBucketPlan(t *testing.T) {
	t.Parallel()
	plan := SharedPlan(t)

	bucket := plan.ResourcePlannedValuesMap["module.statebucket.google_storage_bucket.tf_state"]
	require.NotNil(t, bucket, "state bucket must be planned")
	assert.Equal(t, "tf-state-test", bucket.AttributeValues["name"], "bucket name must be '<bucket_name>-<environment>'")
	assert.Equal(t, "US-CENTRAL1", bucket.AttributeValues["location"], "bucket location must match the region")
	assert.Equal(t, false, bucket.AttributeValues["force_destroy"], "state bucket must not be force-destroyable")
	assert.Equal(t, true, bucket.AttributeValues["uniform_bucket_level_access"], "uniform bucket-level access must be enforced")
	assert.Equal(t, "enforced", bucket.AttributeValues["public_access_prevention"], "public access prevention must be enforced")

	versioning := bucket.AttributeValues["versioning"].([]interface{})
	require.Len(t, versioning, 1, "versioning must be configured")
	assert.Equal(t, true, versioning[0].(map[string]interface{})["enabled"], "versioning must be enabled")

	lifecycle := bucket.AttributeValues["lifecycle_rule"].([]interface{})
	assert.Len(t, lifecycle, 2, "state bucket must keep both lifecycle rules")
}
