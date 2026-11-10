package oracle

import "testing"

// Benchmarks give `make bench` / `make profile` something real to measure so
// performance regressions in the oracles (and the pre-commit hook that runs
// them) show up as numbers instead of vibes.

func BenchmarkCIDRsOverlap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := CIDRsOverlap("10.0.0.0/20", "10.4.0.0/14"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCIDRsOverlapDisjoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		overlap, err := CIDRsOverlap("10.0.0.0/16", "10.255.0.0/16")
		if err != nil {
			b.Fatal(err)
		}
		if overlap {
			b.Fatal("expected disjoint ranges")
		}
	}
}

func BenchmarkValidateCIDRSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := ValidateCIDRSet("10.0.0.0/20", "10.4.0.0/14", "10.8.0.0/20", "172.16.0.0/28"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidProjectID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !ValidProjectID("my-gcp-project") {
			b.Fatal("expected valid project id")
		}
	}
}

func BenchmarkValidNotificationEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !ValidNotificationEmail("ops@example.com") {
			b.Fatal("expected valid email")
		}
	}
}
