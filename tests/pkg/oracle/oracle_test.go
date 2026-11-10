package oracle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Validator tables: these encode the expected behaviour of the HCL validation
// blocks so a regression in either implementation is caught here.
// ---------------------------------------------------------------------------

func TestValidProjectID(t *testing.T) {
	t.Parallel()
	valid := []string{"my-gcp-project", "test-project", "abcde-f", "a12345", "my-project-1"}
	invalid := []string{"", "abc", "-leading", "trailing-", "UPPER-CASE", "has_underscore", "a-b-", "toolongprojectidthatisdefinitelyoverthirtychars"}
	for _, s := range valid {
		assert.True(t, ValidProjectID(s), "%q must be a valid project id", s)
	}
	for _, s := range invalid {
		assert.False(t, ValidProjectID(s), "%q must be rejected as a project id", s)
	}
}

func TestValidRegion(t *testing.T) {
	t.Parallel()
	valid := []string{"us-central1", "europe-west1", "northamerica-northeast1", "australia-southeast1"}
	invalid := []string{"", "us_central1", "US-CENTRAL1", "us-central", "central1", "us-central1a"}
	for _, s := range valid {
		assert.True(t, ValidRegion(s), "%q must be a valid region", s)
	}
	for _, s := range invalid {
		assert.False(t, ValidRegion(s), "%q must be rejected as a region", s)
	}
}

func TestValidClusterName(t *testing.T) {
	t.Parallel()
	valid := []string{"gke-dev", "test-cluster", "a1", strings.Repeat("a", 39) + "1"}
	invalid := []string{"", "a", "-nope", "9-start", "ends-with-", "Has-Caps", strings.Repeat("a", 40) + "1"}
	for _, s := range valid {
		assert.True(t, ValidClusterName(s), "%q must be a valid cluster name", s)
	}
	for _, s := range invalid {
		assert.False(t, ValidClusterName(s), "%q must be rejected as a cluster name", s)
	}
}

func TestValidNetworkName(t *testing.T) {
	t.Parallel()
	assert.True(t, ValidNetworkName("gke-network"))
	assert.True(t, ValidNetworkName("n"+strings.Repeat("x", 61)+"t"))
	assert.False(t, ValidNetworkName("9net"), "names must start with a letter")
	assert.False(t, ValidNetworkName(""))
}

func TestValidBucketName(t *testing.T) {
	t.Parallel()
	valid := []string{"tf-state", "tf-state-test", "a.b_c-d"}
	invalid := []string{"", "ab", "UPPER", "ends-", "-starts", "has space"}
	for _, s := range valid {
		assert.True(t, ValidBucketName(s), "%q must be a valid bucket name", s)
	}
	for _, s := range invalid {
		assert.False(t, ValidBucketName(s), "%q must be rejected as a bucket name", s)
	}
}

func TestValidNotificationEmail(t *testing.T) {
	t.Parallel()
	valid := []string{"", "ops@example.com", "first.last+tag@sub.example.co.uk"}
	invalid := []string{"not-an-email", "a@b", "@example.com", "ops@example", "two @ spaces.com"}
	for _, s := range valid {
		assert.True(t, ValidNotificationEmail(s), "%q must be an acceptable notification email", s)
	}
	for _, s := range invalid {
		assert.False(t, ValidNotificationEmail(s), "%q must be rejected as a notification email", s)
	}
}

func TestValidEnvironment(t *testing.T) {
	t.Parallel()
	for _, s := range []string{"dev", "staging", "production", "test"} {
		assert.True(t, ValidEnvironment(s), "%q must be an allowed environment", s)
	}
	assert.False(t, ValidEnvironment("prod"), `"prod" is not an allowed environment label`)
}

// ---------------------------------------------------------------------------
// CIDR arithmetic, cross-checked against the independent stdlib formulation.
// ---------------------------------------------------------------------------

func TestCIDRsOverlap(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b    string
		overlap bool
		wantErr bool
	}{
		// CIDRs used by environments/dev/terraform.tfvars are disjoint.
		{"10.0.0.0/20", "10.4.0.0/14", false, false},
		{"10.0.0.0/20", "10.8.0.0/20", false, false},
		{"10.4.0.0/14", "10.8.0.0/20", false, false},
		// Nested and identical blocks overlap.
		{"10.0.0.0/20", "10.0.0.0/24", true, false},
		{"10.0.0.0/20", "10.0.0.0/20", true, false},
		{"10.0.0.0/16", "10.0.1.0/24", true, false},
		// Adjacent blocks do not overlap.
		{"10.0.0.0/24", "10.0.1.0/24", false, false},
		// The master CIDR must stay clear of the pod/service ranges.
		{"172.16.0.0/28", "10.4.0.0/14", false, false},
		// IPv6.
		{"2001:db8::/32", "2001:db8:1::/48", true, false},
		{"2001:db8::/32", "2001:db8::/48", true, false},
		// Malformed input errors instead of panicking.
		{"not-a-cidr", "10.0.0.0/20", false, true},
		{"10.0.0.0/20", "", false, true},
		{"999.1.1.1/32", "10.0.0.0/20", false, true},
	}
	for _, tc := range cases {
		got, err := CIDRsOverlap(tc.a, tc.b)
		if tc.wantErr {
			require.Error(t, err, "CIDRsOverlap(%q, %q) must error", tc.a, tc.b)
			continue
		}
		require.NoError(t, err, "CIDRsOverlap(%q, %q)", tc.a, tc.b)
		assert.Equal(t, tc.overlap, got, "CIDRsOverlap(%q, %q)", tc.a, tc.b)

		// Differential check against the independent stdlib implementation.
		ref, refErr := ReferenceOverlap(tc.a, tc.b)
		require.NoError(t, refErr, "reference overlap for %q/%q", tc.a, tc.b)
		assert.Equal(t, ref, got, "interval arithmetic must agree with net.IPNet for %q/%q", tc.a, tc.b)

		// Symmetry.
		reverse, err := CIDRsOverlap(tc.b, tc.a)
		require.NoError(t, err)
		assert.Equal(t, got, reverse, "overlap must be symmetric for %q/%q", tc.a, tc.b)
	}
}

func TestValidateCIDRSet(t *testing.T) {
	t.Parallel()

	// The exact CIDR sets committed for each environment must validate.
	environments := map[string][4]string{
		"dev":     {"10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/28"},
		"staging": {"10.16.0.0/20", "10.20.0.0/14", "10.24.0.0/20", "172.16.0.16/28"},
		"prod":    {"10.32.0.0/20", "10.36.0.0/14", "10.40.0.0/20", "172.16.0.32/28"},
	}
	for name, c := range environments {
		assert.NoError(t, ValidateCIDRSet(c[0], c[1], c[2], c[3]), "environment %q CIDRs must be disjoint", name)
	}

	// Overlaps that per-variable cidrhost() validation cannot detect.
	err := ValidateCIDRSet("10.0.0.0/16", "10.0.1.0/24", "10.8.0.0/20", "172.16.0.0/28")
	require.Error(t, err, "subnet/pods overlap must be rejected")
	assert.Contains(t, err.Error(), "overlaps")

	err = ValidateCIDRSet("10.0.0.0/20", "10.4.0.0/14", "10.4.0.0/20", "172.16.0.0/28")
	require.Error(t, err, "pods/services overlap must be rejected")

	err = ValidateCIDRSet("10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20", "10.0.1.0/28")
	require.Error(t, err, "master inside the subnet must be rejected")

	// Master CIDR must be a /28.
	err = ValidateCIDRSet("10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/30")
	require.Error(t, err, "master CIDR that is not a /28 must be rejected")
	assert.Contains(t, err.Error(), "/28")

	// Malformed input reports which variable is at fault.
	err = ValidateCIDRSet("garbage", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/28")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subnet")
}

// ---------------------------------------------------------------------------
// Drift gate: the Go oracles must stay byte-identical to the HCL regexes.
// ---------------------------------------------------------------------------

// TestHCLValidationMirrorsStayInSync fails when a validation regex in the
// Terraform configuration is added, changed or removed without updating the Go
// mirror (or vice versa). This keeps the fuzzable oracle honest: it is only a
// useful oracle while it matches the source of truth.
func TestHCLValidationMirrorsStayInSync(t *testing.T) {
	t.Parallel()

	hclFiles := []string{
		filepath.Join("..", "..", "..", "tf", "variables.tf"),
		filepath.Join("..", "..", "..", "modules", "gke", "variables.tf"),
	}
	pattern := regexp.MustCompile(`can\(regex\("([^"]*)"`)

	hclSet := map[string]string{} // regex literal -> file it came from
	for _, file := range hclFiles {
		raw, err := os.ReadFile(file)
		require.NoError(t, err, "validation source of truth must be readable")
		for _, m := range pattern.FindAllStringSubmatch(string(raw), -1) {
			literal := strings.ReplaceAll(m[1], `\\`, `\`)
			hclSet[literal] = file
		}
	}

	require.NotEmpty(t, hclSet, "expected to find regex validation blocks in the HCL")

	for name, mirror := range regexMirrors {
		found := false
		for literal := range hclSet {
			if literal == mirror {
				found = true
				break
			}
		}
		assert.True(t, found,
			"oracle %q mirror %q not found in HCL validation blocks; update regexMirrors or the HCL", name, mirror)
	}

	for literal, file := range hclSet {
		found := false
		for _, mirror := range regexMirrors {
			if mirror == literal {
				found = true
				break
			}
		}
		assert.True(t, found,
			"HCL validation regex %q in %s has no Go oracle; add it to regexMirrors", literal, file)
	}
}

// TestEnvironmentTfvarsContainNoSecrets keeps committed environment files free
// of credential-shaped values (they are tracked by git on purpose, so they must
// never hold secret material).
func TestEnvironmentTfvarsContainNoSecrets(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob(filepath.Join("..", "..", "..", "environments", "*", "*.tfvars"))
	require.NoError(t, err)
	require.NotEmpty(t, files, "expected committed environment tfvars files")

	for _, file := range files {
		raw, err := os.ReadFile(file)
		require.NoError(t, err)
		for i, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			value := strings.Trim(parts[1], ` "`)
			assert.False(t, LooksLikeSecret(value),
				"%s:%d looks like credential material; keep secrets out of committed tfvars", file, i+1)
		}
	}
}
