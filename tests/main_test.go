package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestMain(m *testing.M) {
	m.Run()
}

func GetTerraformOptions(t *testing.T) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir: "../tf",
		Vars: map[string]interface{}{
			"project":                    "test-project",
			"region":                     "us-central1",
			"environment":                "test",
			"cluster_name":               "test-cluster",
			"network_name":               "test-network",
			"subnet_cidr":                "10.0.0.0/20",
			"pods_cidr":                  "10.4.0.0/14",
			"services_cidr":              "10.8.0.0/20",
			"master_ipv4_cidr_block":     "172.16.0.0/28",
			"authorized_network_cidr":    "0.0.0.0/0",
			"deletion_protection":        false,
			"main_machine_type":          "e2-medium",
			"main_min_count":             1,
			"main_max_count":             3,
			"main_disk_size_gb":          50,
			"main_disk_type":             "pd-balanced",
			"system_machine_type":        "e2-small",
			"system_min_count":           1,
			"system_max_count":           2,
			"system_disk_size_gb":        50,
			"system_disk_type":           "pd-balanced",
			"namespaces":                 []string{"default", "test"},
			"namespace":                  "default",
			"notification_email":         "",
			"enable_alerting":            false,
			"enable_backup":              false,
			"backup_retention_days":      7,
			"enable_binary_authorization": false,
			"secrets":                    map[string]string{},
			"state_bucket_name":          "tf-state",
		},
	}
}
