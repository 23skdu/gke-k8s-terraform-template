package oracle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fuzz targets for the oracle helpers. They run as ordinary tests over their
// seed corpus (so `go test` and the pre-commit hook always exercise them) and
// become true fuzzing under `make test-fuzz` (go test -fuzz).

// FuzzCIDRsOverlap asserts the overlap oracle never panics, is symmetric, and
// stays in agreement with the independent net.IPNet-based formulation.
func FuzzCIDRsOverlap(f *testing.F) {
	seeds := [][2]string{
		{"10.0.0.0/20", "10.4.0.0/14"},
		{"10.0.0.0/20", "10.0.0.0/24"},
		{"10.0.0.0/24", "10.0.1.0/24"},
		{"172.16.0.0/28", "10.8.0.0/20"},
		{"0.0.0.0/0", "10.0.0.0/32"},
		{"2001:db8::/32", "2001:db8:abcd::/48"},
		{"not-a-cidr", "10.0.0.0/20"},
		{"10.0.0.0/20", ""},
		{"999.999.999.999/32", "10.0.0.0/20"},
	}
	for _, s := range seeds {
		f.Add(s[0], s[1])
	}

	f.Fuzz(func(t *testing.T, a, b string) {
		got, err := CIDRsOverlap(a, b)
		if err != nil {
			// Invalid input must be reported, never panics (checked by the
			// absence of a panic: the fuzzer would surface it).
			return
		}

		// Symmetry: overlap must not depend on argument order.
		reverse, err := CIDRsOverlap(b, a)
		require.NoError(t, err)
		assert.Equal(t, got, reverse, "overlap must be symmetric for %q / %q", a, b)

		// Differential check against the stdlib formulation.
		ref, refErr := ReferenceOverlap(a, b)
		require.NoError(t, refErr)
		assert.Equal(t, ref, got, "interval arithmetic must agree with net.IPNet for %q / %q", a, b)

		// A block always overlaps itself.
		self, err := CIDRsOverlap(a, a)
		require.NoError(t, err)
		assert.True(t, self, "%q must overlap itself", a)
	})
}

// FuzzParseCIDR asserts parsing is panic-free and canonical: re-parsing the
// printed form of a parsed network must reproduce the same interval.
func FuzzParseCIDR(f *testing.F) {
	for _, s := range []string{
		"10.0.0.0/20", "172.16.0.0/28", "0.0.0.0/0", "192.168.1.1/32",
		"2001:db8::/32", "::1/128", "not-a-cidr", "10.0.0.0/33", "", "10.0.0.0/24 ",
	} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		start, end, bits, err := ParseCIDR(s)
		if err != nil {
			return
		}
		require.NotNil(t, start)
		require.NotNil(t, end)
		assert.LessOrEqual(t, start.Cmp(end), 0, "range start must not exceed end")
		assert.Contains(t, []int{32, 128}, bits, "address width must be 32 or 128 bits")
	})
}

// FuzzValidProjectID asserts the project-id oracle enforces the documented
// GCP bounds whenever it accepts an input.
func FuzzValidProjectID(f *testing.F) {
	for _, s := range []string{"my-gcp-project", "test-project", "abc", "-bad", "", "a123456789012345678901234567890123456789"} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		if !ValidProjectID(s) {
			return
		}
		// GCP project ids: 6-30 chars, start with a letter, end alphanumeric.
		assert.GreaterOrEqual(t, len(s), 6, "accepted project id %q must be at least 6 chars", s)
		assert.LessOrEqual(t, len(s), 30, "accepted project id %q must be at most 30 chars", s)
		assert.Regexp(t, `^[a-z]`, s, "accepted project id %q must start with a lowercase letter", s)
		assert.Regexp(t, `[a-z0-9]$`, s, "accepted project id %q must end alphanumeric", s)
	})
}

// FuzzValidClusterName asserts accepted cluster names stay within the GKE
// 1-40 character, lowercase/hyphen charset.
func FuzzValidClusterName(f *testing.F) {
	for _, s := range []string{"gke-dev", "test-cluster", "a1", "-bad", "CAPS", ""} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		if !ValidClusterName(s) {
			return
		}
		assert.GreaterOrEqual(t, len(s), 2, "accepted cluster name %q is too short", s)
		assert.LessOrEqual(t, len(s), 40, "accepted cluster name %q is too long", s)
		assert.Regexp(t, `^[a-z][a-z0-9-]*[a-z0-9]$`, s, "accepted cluster name %q has an illegal charset", s)
	})
}

// FuzzValidateCIDRSet asserts the combined CIDR gate is panic-free and that
// every set it accepts is genuinely pairwise-disjoint under the independent
// reference implementation, with a /28 master block.
func FuzzValidateCIDRSet(f *testing.F) {
	seeds := [][4]string{
		{"10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/28"},
		{"10.16.0.0/20", "10.20.0.0/14", "10.24.0.0/20", "172.16.0.16/28"},
		{"10.0.0.0/16", "10.0.1.0/24", "10.8.0.0/20", "172.16.0.0/28"},
		{"garbage", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/28"},
	}
	for _, s := range seeds {
		f.Add(s[0], s[1], s[2], s[3])
	}

	f.Fuzz(func(t *testing.T, subnet, pods, services, master string) {
		err := ValidateCIDRSet(subnet, pods, services, master)
		if err != nil {
			return
		}
		// Everything accepted must parse...
		set := map[string]string{"subnet": subnet, "pods": pods, "services": services, "master": master}
		for name, cidr := range set {
			_, _, _, perr := ParseCIDR(cidr)
			require.NoError(t, perr, "accepted %s CIDR %q must parse", name, cidr)
		}
		// ...be pairwise disjoint under the reference implementation...
		pairs := [][2]string{
			{"subnet", "pods"}, {"subnet", "services"}, {"subnet", "master"},
			{"pods", "services"}, {"pods", "master"}, {"services", "master"},
		}
		for _, p := range pairs {
			ref, rerr := ReferenceOverlap(set[p[0]], set[p[1]])
			require.NoError(t, rerr)
			assert.False(t, ref, "accepted set has %s/%s overlapping", p[0], p[1])
		}
		// ...and carry a /28 master block.
		ones, perr := prefixLen(master)
		require.NoError(t, perr)
		assert.Equal(t, 28, ones, "accepted master CIDR must be a /28")
	})
}
