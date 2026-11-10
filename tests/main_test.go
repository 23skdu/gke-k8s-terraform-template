package test

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestMain tears down the shared hermetic plan fixture after the whole suite
// finishes and propagates the real exit code to the caller.
func TestMain(m *testing.M) {
	code := m.Run()
	teardownFixture()
	os.Exit(code)
}

// ---------------------------------------------------------------------------
// Hermetic plan fixture
//
// Plan-only tests must be free of GCP credentials, remote state and network
// side effects. They therefore never plan inside ../tf (which declares a GCS
// backend and would require real credentials); instead the first caller builds
// a backend-free copy of the root module plus the reusable modules tree in a
// temporary directory, and every plan-only test shares a single plan result.
// ---------------------------------------------------------------------------

var (
	fixtureOnce sync.Once
	fixturePath string
	fixtureErr  error

	planOnce   sync.Once
	sharedPlan *terraform.PlanStruct
)

// fixtureSkip lists files that must not leak into the backend-free fixture.
var fixtureSkip = map[string]bool{
	"backend.tf":          true, // GCS backend declaration
	"terraform.tfvars":    true, // local overrides
	"tfplan":              true, // stale binary plan
	".terraform":          true, // plugin/cache state
	".terraform.lock.hcl": true,
}

func fixtureRoot() (string, error) {
	fixtureOnce.Do(func() {
		repoRoot, err := filepath.Abs("..")
		if err != nil {
			fixtureErr = err
			return
		}
		tmp, err := os.MkdirTemp("", "gke-tf-fixture-")
		if err != nil {
			fixtureErr = err
			return
		}
		fixturePath = tmp
		for _, dir := range []string{"tf", "modules"} {
			if err := copyTree(filepath.Join(repoRoot, dir), filepath.Join(tmp, dir)); err != nil {
				fixtureErr = err
				return
			}
		}
	})
	return fixturePath, fixtureErr
}

func teardownFixture() {
	if fixturePath != "" {
		_ = os.RemoveAll(fixturePath)
	}
}

// copyTree copies src into dst, skipping fixtureSkip entries and any .terraform
// directory, preserving the relative layout that module sources rely on
// (tf/ references ../modules).
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		base := d.Name()
		if d.IsDir() {
			if base == ".terraform" || fixtureSkip[base] {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if fixtureSkip[base] {
			return nil
		}
		return copyFile(path, filepath.Join(dst, rel))
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// GetPlanTerraformOptions returns options for hermetic plan-only tests: the
// backend-free fixture, a static (fake) access token so provider configuration
// never needs real GCP credentials, and an enriched feature-flag set so that
// ONE plan exercises every conditional resource in the configuration
// (backup, binary authorization, alerting, notification channel, secrets).
func GetPlanTerraformOptions(t *testing.T) *terraform.Options {
	t.Helper()
	fixture, err := fixtureRoot()
	require.NoError(t, err, "plan fixture must be built")

	return &terraform.Options{
		TerraformDir: filepath.Join(fixture, "tf"),
		PlanFilePath: filepath.Join(fixture, "tfplan.out"),
		NoColor:      true,
		EnvVars: map[string]string{
			// Plan-only tests must never read ambient GCP credentials: a
			// static token satisfies `provider "google"` configuration and
			// `data.google_client_config` without performing any API call.
			"GOOGLE_OAUTH_ACCESS_TOKEN": "hermetic-plan-only-token",
		},
		Vars: map[string]interface{}{
			"project":                 "test-project",
			"region":                  "us-central1",
			"environment":             "test",
			"cluster_name":            "test-cluster",
			"network_name":            "test-network",
			"subnet_cidr":             "10.0.0.0/20",
			"pods_cidr":               "10.4.0.0/14",
			"services_cidr":           "10.8.0.0/20",
			"master_ipv4_cidr_block":  "172.16.0.0/28",
			"authorized_network_cidr": "0.0.0.0/0",
			"deletion_protection":     false,
			"main_machine_type":       "e2-medium",
			"main_min_count":          1,
			"main_max_count":          3,
			"main_disk_size_gb":       50,
			"main_disk_type":          "pd-balanced",
			"system_machine_type":     "e2-small",
			"system_min_count":        1,
			"system_max_count":        2,
			"system_disk_size_gb":     50,
			"system_disk_type":        "pd-balanced",
			"namespaces":              []string{"default", "test"},
			"namespace":               "default",
			// Enable every optional feature so a single shared plan covers the
			// backup, security, monitoring and secrets modules as well.
			"notification_email":          "alerts@example.com",
			"enable_alerting":             true,
			"enable_backup":               true,
			"backup_retention_days":       14,
			"enable_binary_authorization": true,
			"secrets":                     map[string]string{"db-password": "Database connection password"},
			"state_bucket_name":           "tf-state",
		},
	}
}

// SharedPlan runs terraform init + plan + show exactly once for the whole test
// binary and returns the result for read-only assertions. Sharing a single plan
// turns N plan invocations into 1, which keeps the plan-only suite fast enough
// to run from a pre-commit hook.
func SharedPlan(t *testing.T) *terraform.PlanStruct {
	t.Helper()
	planOnce.Do(func() {
		sharedPlan = terraform.InitAndPlanAndShowWithStruct(t, GetPlanTerraformOptions(t))
	})
	require.NotNil(t, sharedPlan, "shared plan fixture must be available to every plan-only test")
	return sharedPlan
}

// GetTerraformOptions returns options for integration (apply) tests. These run
// against the real ../tf directory with its GCS backend and create real
// resources, so they are skipped entirely under `go test -short`.
func GetTerraformOptions(t *testing.T) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir: "../tf",
		Vars: map[string]interface{}{
			"project":                     "test-project",
			"region":                      "us-central1",
			"environment":                 "test",
			"cluster_name":                "test-cluster",
			"network_name":                "test-network",
			"subnet_cidr":                 "10.0.0.0/20",
			"pods_cidr":                   "10.4.0.0/14",
			"services_cidr":               "10.8.0.0/20",
			"master_ipv4_cidr_block":      "172.16.0.0/28",
			"authorized_network_cidr":     "0.0.0.0/0",
			"deletion_protection":         false,
			"main_machine_type":           "e2-medium",
			"main_min_count":              1,
			"main_max_count":              3,
			"main_disk_size_gb":           50,
			"main_disk_type":              "pd-balanced",
			"system_machine_type":         "e2-small",
			"system_min_count":            1,
			"system_max_count":            2,
			"system_disk_size_gb":         50,
			"system_disk_type":            "pd-balanced",
			"namespaces":                  []string{"default", "test"},
			"namespace":                   "default",
			"notification_email":          "",
			"enable_alerting":             false,
			"enable_backup":               false,
			"backup_retention_days":       7,
			"enable_binary_authorization": false,
			"secrets":                     map[string]string{},
			"state_bucket_name":           "tf-state",
		},
	}
}

// ---------------------------------------------------------------------------
// Plan assertion helpers: value-level oracles instead of existence checks.
// ---------------------------------------------------------------------------

// requireAttr returns the value of a top-level attribute of a planned resource
// and fails the test when the resource or attribute is absent.
func requireAttr(t *testing.T, plan *terraform.PlanStruct, address, key string) interface{} {
	t.Helper()
	res := plan.ResourcePlannedValuesMap[address]
	require.NotNil(t, res, "planned resource %s must exist", address)
	val, ok := res.AttributeValues[key]
	require.True(t, ok, "resource %s must have attribute %q", address, key)
	return val
}

// requireBlock returns the first element of a nested block attribute (Terraform
// JSON encodes blocks as a list of objects) and fails otherwise.
func requireBlock(t *testing.T, plan *terraform.PlanStruct, address, key string) map[string]interface{} {
	t.Helper()
	raw := requireAttr(t, plan, address, key)
	list, ok := raw.([]interface{})
	require.True(t, ok, "%s.%s must be a list of blocks, got %T", address, key, raw)
	require.NotEmpty(t, list, "%s.%s must contain at least one block", address, key)
	block, ok := list[0].(map[string]interface{})
	require.True(t, ok, "%s.%s[0] must be an object, got %T", address, key, list[0])
	return block
}

// requireFloat normalizes JSON numbers (always float64 after plan parsing) so
// assertions read as exact values instead of type accidents.
func requireFloat(t *testing.T, v interface{}) float64 {
	t.Helper()
	f, ok := v.(float64)
	require.True(t, ok, "expected numeric value, got %T", v)
	return f
}

// requireString normalizes a string attribute.
func requireString(t *testing.T, v interface{}) string {
	t.Helper()
	s, ok := v.(string)
	require.True(t, ok, "expected string value, got %T", v)
	return s
}

// containsAll reports whether haystack contains every needle.
func containsAll(haystack []interface{}, needles ...string) bool {
	seen := map[string]bool{}
	for _, h := range haystack {
		if s, ok := h.(string); ok {
			seen[s] = true
		}
	}
	for _, n := range needles {
		if !seen[n] {
			return false
		}
	}
	return true
}

// integrationSkip skips tests that create real GCP resources unless the full
// integration suite was requested (`go test` without -short).
func integrationSkip(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test: creates real GCP resources; skipped with -short (plan-only mode)")
	}
}
