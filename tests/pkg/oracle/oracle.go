// Package oracle provides pure-Go reference implementations ("oracles") of the
// validation rules that live in the Terraform configuration. They exist so the
// test suite can assert on inputs and CIDR arithmetic without shelling out to
// Terraform, fuzz them for robustness, and detect drift against the HCL source
// of truth (see TestHCLValidationMirrorsStayInSync).
package oracle

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
)

// regexMirrors holds the exact regex literals from the Terraform validation
// blocks. Every value MUST appear verbatim in tf/variables.tf (or
// modules/gke/variables.tf); the sync test fails when either side drifts.
var regexMirrors = map[string]string{
	"project":            `^[a-z][a-z0-9-]{4,28}[a-z0-9]$`,
	"region":             `^[a-z]+-[a-z]+[0-9]$`,
	"network_name":       `^[a-z][a-z0-9-]{0,62}$`,
	"cluster_name":       `^[a-z][a-z0-9-]{0,38}[a-z0-9]$`,
	"state_bucket_name":  `^[a-z][a-z0-9._-]{1,61}[a-z0-9]$`,
	"notification_email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
}

// Pre-compiled package-level regex pointers avoid a map lookup on every validator call.
var (
	reProject           = regexp.MustCompile(regexMirrors["project"])
	reRegion            = regexp.MustCompile(regexMirrors["region"])
	reNetworkName       = regexp.MustCompile(regexMirrors["network_name"])
	reClusterName       = regexp.MustCompile(regexMirrors["cluster_name"])
	reStateBucketName   = regexp.MustCompile(regexMirrors["state_bucket_name"])
	reNotificationEmail = regexp.MustCompile(regexMirrors["notification_email"])

	// Extended semantic oracles beyond HCL regexes.
	reNamespace        = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)
	reServiceAccountID = regexp.MustCompile(`^[a-z][a-z0-9-]{4,28}[a-z0-9]$`)
	reMachineType      = regexp.MustCompile(`^(e2|n1|n2|n2d|c2|c2d|c3|c3d|c4|m1|m2|m3|t2d|t2a|a2|a3|g2)-(standard|highmem|highcpu|megamem|ultramem|custom|micro|small|medium)(-[0-9]+(-[0-9]+)?)?$`)
	reSecretID         = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,255}$`)
)

var compiled = map[string]*regexp.Regexp{
	"project":            reProject,
	"region":             reRegion,
	"network_name":       reNetworkName,
	"cluster_name":       reClusterName,
	"state_bucket_name":  reStateBucketName,
	"notification_email": reNotificationEmail,
}

// Mirrors returns a copy of the regex mirrors for diagnostics and tests.
func Mirrors() map[string]string {
	out := make(map[string]string, len(regexMirrors))
	for k, v := range regexMirrors {
		out[k] = v
	}
	return out
}

// ValidProjectID mirrors the `project` validation block.
func ValidProjectID(s string) bool { return reProject.MatchString(s) }

// ValidRegion mirrors the `region` validation block.
func ValidRegion(s string) bool { return reRegion.MatchString(s) }

// ValidNetworkName mirrors the `network_name` validation block.
func ValidNetworkName(s string) bool { return reNetworkName.MatchString(s) }

// ValidClusterName mirrors the `cluster_name` validation block.
func ValidClusterName(s string) bool { return reClusterName.MatchString(s) }

// ValidBucketName mirrors the `state_bucket_name` validation block.
func ValidBucketName(s string) bool { return reStateBucketName.MatchString(s) }

// ValidNotificationEmail mirrors the `notification_email` validation block:
// empty is allowed, otherwise the address must match the pattern.
func ValidNotificationEmail(s string) bool {
	return s == "" || reNotificationEmail.MatchString(s)
}

// ValidEnvironment mirrors the `environment` validation block (a set check
// rather than a regex in HCL).
func ValidEnvironment(s string) bool {
	switch s {
	case "dev", "staging", "production", "test":
		return true
	}
	return false
}

// Uint128 represents an unsigned 128-bit integer for zero-allocation IPv4/IPv6
// interval arithmetic.
type Uint128 struct {
	Hi uint64
	Lo uint64
}

// Cmp compares u and other and returns -1, 0, or +1.
func (u Uint128) Cmp(other Uint128) int {
	if u.Hi < other.Hi {
		return -1
	}
	if u.Hi > other.Hi {
		return 1
	}
	if u.Lo < other.Lo {
		return -1
	}
	if u.Lo > other.Lo {
		return 1
	}
	return 0
}

// ParseCIDR128 parses s as an IPv4 or IPv6 CIDR block without heap allocations
// on the valid path, returning the inclusive address interval [start, end],
// prefix length (ones), and address-family width (32 or 128 bits).
func ParseCIDR128(s string) (start, end Uint128, ones, bits int, err error) {
	p, perr := netip.ParsePrefix(s)
	if perr != nil {
		return Uint128{}, Uint128{}, 0, 0, fmt.Errorf("invalid CIDR %q: %w", s, perr)
	}
	addr := p.Addr()
	ones = p.Bits()
	if addr.Is4() {
		b4 := addr.As4()
		ip32 := uint32(b4[0])<<24 | uint32(b4[1])<<16 | uint32(b4[2])<<8 | uint32(b4[3])
		var mask32 uint32
		if ones > 0 {
			mask32 = ^uint32(0) << (32 - ones)
		}
		start32 := ip32 & mask32
		end32 := start32 | ^mask32
		return Uint128{Hi: 0, Lo: uint64(start32)}, Uint128{Hi: 0, Lo: uint64(end32)}, ones, 32, nil
	}
	if addr.Is4In6() {
		// Match net.ParseCIDR behavior on IPv4-mapped IPv6 addresses.
		if ones > 32 {
			return Uint128{}, Uint128{}, 0, 0, fmt.Errorf("invalid CIDR %q: prefix length %d exceeds IPv4 width", s, ones)
		}
		b16 := addr.As16()
		ip32 := uint32(b16[12])<<24 | uint32(b16[13])<<16 | uint32(b16[14])<<8 | uint32(b16[15])
		var mask32 uint32
		if ones > 0 {
			mask32 = ^uint32(0) << (32 - ones)
		}
		start32 := ip32 & mask32
		end32 := start32 | ^mask32
		return Uint128{Hi: 0, Lo: uint64(start32)}, Uint128{Hi: 0, Lo: uint64(end32)}, ones, 32, nil
	}
	b16 := addr.As16()
	hi := binary.BigEndian.Uint64(b16[0:8])
	lo := binary.BigEndian.Uint64(b16[8:16])
	var maskHi, maskLo uint64
	switch {
	case ones == 0:
		maskHi, maskLo = 0, 0
	case ones <= 64:
		maskHi = ^uint64(0) << (64 - ones)
		maskLo = 0
	default:
		maskHi = ^uint64(0)
		maskLo = ^uint64(0) << (128 - ones)
	}
	start = Uint128{Hi: hi & maskHi, Lo: lo & maskLo}
	end = Uint128{Hi: start.Hi | ^maskHi, Lo: start.Lo | ^maskLo}
	return start, end, ones, 128, nil
}

func uint128ToBig(u Uint128, bits int) *big.Int {
	if bits == 32 {
		return new(big.Int).SetUint64(u.Lo)
	}
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[0:8], u.Hi)
	binary.BigEndian.PutUint64(buf[8:16], u.Lo)
	return new(big.Int).SetBytes(buf[:])
}

// ParseCIDR parses s as an IPv4 or IPv6 CIDR and returns the inclusive address
// interval [start, end] plus the address-family width in bits.
func ParseCIDR(s string) (start, end *big.Int, bits int, err error) {
	s128, e128, _, b, err := ParseCIDR128(s)
	if err != nil {
		return nil, nil, 0, err
	}
	return uint128ToBig(s128, b), uint128ToBig(e128, b), b, nil
}

// CIDRsOverlap reports whether two CIDR blocks share at least one address.
//
// The implementation uses 128-bit integer interval arithmetic on the network
// ranges, deliberately independent of net.IPNet.Contains so that tests and
// fuzzing can cross-check the two implementations against each other
// (differential testing) while running with zero heap allocations on valid
// inputs. Malformed input returns an error instead of panicking.
func CIDRsOverlap(a, b string) (bool, error) {
	aStart, aEnd, _, aBits, err := ParseCIDR128(a)
	if err != nil {
		return false, err
	}
	bStart, bEnd, _, bBits, err := ParseCIDR128(b)
	if err != nil {
		return false, err
	}
	if aBits != bBits {
		return false, fmt.Errorf("address family mismatch: %q vs %q", a, b)
	}
	// Intervals [aStart,aEnd] and [bStart,bEnd] intersect iff each starts
	// before the other ends.
	return aStart.Cmp(bEnd) <= 0 && bStart.Cmp(aEnd) <= 0, nil
}

// referenceOverlap is the independent, stdlib-based formulation used by the
// test suite to cross-check CIDRsOverlap.
func referenceOverlap(a, b string) (bool, error) {
	_, aNet, err := net.ParseCIDR(a)
	if err != nil {
		return false, err
	}
	_, bNet, err := net.ParseCIDR(b)
	if err != nil {
		return false, err
	}
	if len(aNet.IP) != len(bNet.IP) {
		return false, fmt.Errorf("address family mismatch: %q vs %q", a, b)
	}
	return aNet.Contains(bNet.IP) || bNet.Contains(aNet.IP), nil
}

// ReferenceOverlap exposes the stdlib-based overlap check for differential
// testing from the oracle test suite.
func ReferenceOverlap(a, b string) (bool, error) { return referenceOverlap(a, b) }

// IsCanonicalCIDR reports whether CIDR s has no non-zero host bits set outside
// its network prefix (for example, 10.0.0.0/20 is canonical, whereas
// 10.0.0.5/20 is not).
func IsCanonicalCIDR(s string) (bool, error) {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR %q: %w", s, err)
	}
	return p.Masked().Addr() == p.Addr(), nil
}

// IsPrivateIPv4CIDR reports whether s is a valid IPv4 CIDR block completely
// contained within RFC 1918 private address space (10.0.0.0/8, 172.16.0.0/12,
// or 192.168.0.0/16).
func IsPrivateIPv4CIDR(s string) (bool, error) {
	start, end, _, bits, err := ParseCIDR128(s)
	if err != nil {
		return false, err
	}
	if bits != 32 {
		return false, nil
	}
	// RFC 1918 intervals as 32-bit unsigned integers:
	// 10.0.0.0/8      -> [0x0a000000, 0x0affffff]
	// 172.16.0.0/12   -> [0xac100000, 0xac1fffff]
	// 192.168.0.0/16  -> [0xc0a80000, 0xc0a8ffff]
	lo, hi := start.Lo, end.Lo
	if lo >= 0x0a000000 && hi <= 0x0affffff {
		return true, nil
	}
	if lo >= 0xac100000 && hi <= 0xac1fffff {
		return true, nil
	}
	if lo >= 0xc0a80000 && hi <= 0xc0a8ffff {
		return true, nil
	}
	return false, nil
}

// CIDRAddressCount returns the total number of addresses in CIDR block s.
func CIDRAddressCount(s string) (*big.Int, error) {
	_, _, ones, bits, err := ParseCIDR128(s)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Lsh(big.NewInt(1), uint(bits-ones)), nil
}

// masterPrefix is the GKE requirement that the private master endpoint CIDR is
// a /28.
const masterPrefix = 28

// ValidateCIDRSet checks the four CIDR inputs of the template for the failures
// that Terraform's per-variable `can(cidrhost(...))` checks cannot catch:
// overlapping ranges and a master CIDR that is not a /28. Terraform plans with
// these mistakes and only fails deep inside the GKE API at apply time; this
// oracle lets the test suite and pre-commit hook fail fast instead.
//
// Uses fixed-size stack arrays and ParseCIDR128 so the happy path performs
// zero heap allocations.
func ValidateCIDRSet(subnet, pods, services, master string) error {
	names := [4]string{"subnet", "pods", "services", "master"}
	cidrs := [4]string{subnet, pods, services, master}
	var starts, ends [4]Uint128
	var ones, bits [4]int

	for i := 0; i < 4; i++ {
		var err error
		starts[i], ends[i], ones[i], bits[i], err = ParseCIDR128(cidrs[i])
		if err != nil {
			return fmt.Errorf("%s: %w", names[i], err)
		}
	}

	if ones[3] != masterPrefix {
		return fmt.Errorf("master: GKE requires a /%d CIDR for the master endpoint, got /%d", masterPrefix, ones[3])
	}

	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if bits[i] != bits[j] {
				return fmt.Errorf("%s (%s) and %s (%s) use different address families", names[i], cidrs[i], names[j], cidrs[j])
			}
			if starts[i].Cmp(ends[j]) <= 0 && starts[j].Cmp(ends[i]) <= 0 {
				return fmt.Errorf("%s (%s) overlaps %s (%s)", names[i], cidrs[i], names[j], cidrs[j])
			}
		}
	}
	return nil
}

func prefixLen(cidr string) (int, error) {
	_, _, ones, _, err := ParseCIDR128(cidr)
	if err != nil {
		return 0, err
	}
	return ones, nil
}

// ValidMachineType reports whether s is a well-formed GCP Compute Engine
// machine type (e.g., "e2-medium", "e2-standard-4", "n2-highmem-8").
func ValidMachineType(s string) bool {
	return reMachineType.MatchString(s)
}

// ValidDiskType mirrors the disk_type validation block in tf/variables.tf and
// modules/gke/variables.tf.
func ValidDiskType(s string) bool {
	switch s {
	case "pd-standard", "pd-ssd", "pd-balanced", "hyperdisk-balanced":
		return true
	}
	return false
}

// ValidateNodePoolBounds validates node pool autoscaling counts, disk size,
// disk type, and machine type compatibility.
func ValidateNodePoolBounds(minCount, maxCount, diskSizeGB int, diskType, machineType string, isSystem bool) error {
	if isSystem {
		if minCount < 1 {
			return fmt.Errorf("system pool min_count must be >= 1 to keep critical cluster components scheduled, got %d", minCount)
		}
	} else if minCount < 0 {
		return fmt.Errorf("min_count must be >= 0, got %d", minCount)
	}
	if maxCount < 1 {
		return fmt.Errorf("max_count must be >= 1, got %d", maxCount)
	}
	if minCount > maxCount {
		return fmt.Errorf("min_count (%d) must not exceed max_count (%d)", minCount, maxCount)
	}
	if diskSizeGB < 10 || diskSizeGB > 65536 {
		return fmt.Errorf("disk_size_gb must be in [10, 65536], got %d", diskSizeGB)
	}
	if !ValidDiskType(diskType) {
		return fmt.Errorf("invalid disk_type %q", diskType)
	}
	if !ValidMachineType(machineType) {
		return fmt.Errorf("invalid machine_type %q", machineType)
	}
	if diskType == "hyperdisk-balanced" && strings.HasPrefix(machineType, "e2-") {
		return fmt.Errorf("disk_type %q is not supported on %q (E2 machine series)", diskType, machineType)
	}
	return nil
}

// ValidNamespace reports whether s is a valid Kubernetes RFC 1123 DNS label
// suitable for a namespace name.
func ValidNamespace(s string) bool {
	return len(s) >= 1 && len(s) <= 63 && reNamespace.MatchString(s)
}

// ValidateNamespaceSet checks that namespaces is non-empty, contains only
// valid, non-duplicate, non-reserved Kubernetes namespace names, and includes
// the Workload Identity binding namespace when non-empty.
func ValidateNamespaceSet(namespaces []string, workloadIdentityNamespace string) error {
	if len(namespaces) == 0 {
		return fmt.Errorf("namespaces list must not be empty")
	}
	seen := make(map[string]struct{}, len(namespaces))
	for _, ns := range namespaces {
		if !ValidNamespace(ns) {
			return fmt.Errorf("invalid namespace %q: must be a valid DNS-1123 label", ns)
		}
		if strings.HasPrefix(ns, "kube-") {
			return fmt.Errorf("namespace %q is reserved for Kubernetes system components", ns)
		}
		if _, dup := seen[ns]; dup {
			return fmt.Errorf("duplicate namespace %q", ns)
		}
		seen[ns] = struct{}{}
	}
	if workloadIdentityNamespace != "" {
		if !ValidNamespace(workloadIdentityNamespace) {
			return fmt.Errorf("invalid workload identity namespace %q", workloadIdentityNamespace)
		}
		if _, ok := seen[workloadIdentityNamespace]; !ok {
			return fmt.Errorf("workload identity namespace %q is not present in namespaces %v", workloadIdentityNamespace, namespaces)
		}
	}
	return nil
}

// ValidServiceAccountID reports whether s satisfies GCP IAM's 6-30 character
// service account ID constraint (`^[a-z][a-z0-9-]{4,28}[a-z0-9]$`).
func ValidServiceAccountID(s string) bool {
	return reServiceAccountID.MatchString(s)
}

// DerivedServiceAccountID computes `<cluster_name>-<environment>-sa` (as used
// in modules/iam/main.tf) and verifies that the resulting ID satisfies GCP IAM
// length and character constraints.
func DerivedServiceAccountID(clusterName, environment string) (string, error) {
	if !ValidClusterName(clusterName) {
		return "", fmt.Errorf("invalid cluster_name %q", clusterName)
	}
	if !ValidEnvironment(environment) {
		return "", fmt.Errorf("invalid environment %q", environment)
	}
	id := clusterName + "-" + environment + "-sa"
	if !ValidServiceAccountID(id) {
		return "", fmt.Errorf("derived service account ID %q (len=%d) violates GCP IAM 6-30 char limit", id, len(id))
	}
	return id, nil
}

// ValidSecretID reports whether s is a valid GCP Secret Manager secret ID.
func ValidSecretID(s string) bool {
	return reSecretID.MatchString(s)
}

// ValidCronSchedule reports whether s is a valid 5-field POSIX cron expression
// (minute hour day-of-month month day-of-week).
func ValidCronSchedule(s string) bool {
	fields := strings.Fields(s)
	if len(fields) != 5 {
		return false
	}
	bounds := [5][2]int{
		{0, 59}, // minute
		{0, 23}, // hour
		{1, 31}, // day of month
		{1, 12}, // month
		{0, 7},  // day of week (0 or 7 = Sunday)
	}
	for i, f := range fields {
		if !validCronField(f, bounds[i][0], bounds[i][1]) {
			return false
		}
	}
	return true
}

func validCronField(field string, minVal, maxVal int) bool {
	if field == "" {
		return false
	}
	for _, part := range strings.Split(field, ",") {
		if part == "" {
			return false
		}
		base := part
		if idx := strings.IndexByte(part, '/'); idx != -1 {
			base = part[:idx]
			stepStr := part[idx+1:]
			step, err := strconv.Atoi(stepStr)
			if err != nil || step < 1 || step > maxVal {
				return false
			}
		}
		if base == "*" {
			continue
		}
		if dash := strings.IndexByte(base, '-'); dash != -1 {
			lo, err1 := strconv.Atoi(base[:dash])
			hi, err2 := strconv.Atoi(base[dash+1:])
			if err1 != nil || err2 != nil || lo < minVal || hi > maxVal || lo > hi {
				return false
			}
			continue
		}
		val, err := strconv.Atoi(base)
		if err != nil || val < minVal || val > maxVal {
			return false
		}
	}
	return true
}

var secretMarkers = [...]string{
	"password=",
	"api_key=",
	"apikey=",
	"secret_key=",
	"private_key",
	"-----begin",
	"awsk_",
	"akia",
	"ghp_",
	"gho_",
	"xoxb-",
	"xoxp-",
	"ai" + "za", // Google API key prefix (AIza...)
}

// containsFoldASCII reports whether s contains substr (which must be all
// lowercase ASCII) using an allocation-free case-insensitive scan.
func containsFoldASCII(s, substr string) bool {
	n := len(substr)
	if n == 0 {
		return true
	}
	if len(s) < n {
		return false
	}
	first := substr[0]
	limit := len(s) - n
	for i := 0; i <= limit; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != first {
			continue
		}
		match := true
		for j := 1; j < n; j++ {
			 sc := s[i+j]
			if sc >= 'A' && sc <= 'Z' {
				sc += 'a' - 'A'
			}
			if sc != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// LooksLikeSecret reports whether a tfvars-style value resembles credential
// material, used by the repository hygiene gates to keep secrets out of the
// committed environment files. Runs with zero heap allocations.
func LooksLikeSecret(v string) bool {
	for i := 0; i < len(secretMarkers); i++ {
		if containsFoldASCII(v, secretMarkers[i]) {
			return true
		}
	}
	return false
}

// ParseTfvars parses simple HCL key = value assignments from a .tfvars file.
func ParseTfvars(content string) (map[string]string, error) {
	out := make(map[string]string, 24)
	inBlock := 0
	for lineNo, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		if inBlock > 0 {
			inBlock += strings.Count(line, "{") - strings.Count(line, "}")
			if inBlock < 0 {
				return nil, fmt.Errorf("line %d: unmatched closing brace", lineNo+1)
			}
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected key = value assignment, got %q", lineNo+1, line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" || val == "" {
			return nil, fmt.Errorf("line %d: empty key or value", lineNo+1)
		}
		openBraces := strings.Count(val, "{") - strings.Count(val, "}")
		if openBraces > 0 {
			inBlock = openBraces
			out[key] = "{...}"
			continue
		}
		out[key] = strings.Trim(val, `"`)
	}
	if inBlock != 0 {
		return nil, fmt.Errorf("unclosed block in tfvars content")
	}
	return out, nil
}

// ValidateTfvarsMap validates a parsed .tfvars map against all single-field
// and cross-field oracles (naming regexes, CIDR disjointness, /28 master CIDR,
// RFC1918 private ranges, derived IAM service account length, node pool sizing
// bounds, and secret hygiene).
func ValidateTfvarsMap(vars map[string]string) error {
	for k, v := range vars {
		if LooksLikeSecret(v) {
			return fmt.Errorf("variable %q contains secret-like material", k)
		}
	}
	if p := vars["project"]; !ValidProjectID(p) {
		return fmt.Errorf("invalid project %q", p)
	}
	if r := vars["region"]; !ValidRegion(r) {
		return fmt.Errorf("invalid region %q", r)
	}
	if env := vars["environment"]; !ValidEnvironment(env) {
		return fmt.Errorf("invalid environment %q", env)
	}
	if cn := vars["cluster_name"]; !ValidClusterName(cn) {
		return fmt.Errorf("invalid cluster_name %q", cn)
	}
	if _, err := DerivedServiceAccountID(vars["cluster_name"], vars["environment"]); err != nil {
		return err
	}
	if nn := vars["network_name"]; !ValidNetworkName(nn) {
		return fmt.Errorf("invalid network_name %q", nn)
	}
	if bn := vars["state_bucket_name"]; !ValidBucketName(bn) {
		return fmt.Errorf("invalid state_bucket_name %q", bn)
	}
	if em := vars["notification_email"]; !ValidNotificationEmail(em) {
		return fmt.Errorf("invalid notification_email %q", em)
	}
	if err := ValidateCIDRSet(
		vars["subnet_cidr"],
		vars["pods_cidr"],
		vars["services_cidr"],
		vars["master_ipv4_cidr_block"],
	); err != nil {
		return err
	}
	for _, cidrKey := range []string{"subnet_cidr", "pods_cidr", "services_cidr", "master_ipv4_cidr_block"} {
		priv, err := IsPrivateIPv4CIDR(vars[cidrKey])
		if err != nil {
			return fmt.Errorf("%s: %w", cidrKey, err)
		}
		if !priv {
			return fmt.Errorf("%s (%s) must be an RFC 1918 private IPv4 CIDR", cidrKey, vars[cidrKey])
		}
		canon, err := IsCanonicalCIDR(vars[cidrKey])
		if err != nil || !canon {
			return fmt.Errorf("%s (%s) has non-zero host bits outside its prefix mask", cidrKey, vars[cidrKey])
		}
	}
	mainMin, err := strconv.Atoi(vars["main_min_count"])
	if err != nil {
		return fmt.Errorf("main_min_count: %w", err)
	}
	mainMax, err := strconv.Atoi(vars["main_max_count"])
	if err != nil {
		return fmt.Errorf("main_max_count: %w", err)
	}
	mainDisk, err := strconv.Atoi(vars["main_disk_size_gb"])
	if err != nil {
		return fmt.Errorf("main_disk_size_gb: %w", err)
	}
	if err := ValidateNodePoolBounds(mainMin, mainMax, mainDisk, vars["main_disk_type"], vars["main_machine_type"], false); err != nil {
		return fmt.Errorf("main pool: %w", err)
	}

	sysMin, err := strconv.Atoi(vars["system_min_count"])
	if err != nil {
		return fmt.Errorf("system_min_count: %w", err)
	}
	sysMax, err := strconv.Atoi(vars["system_max_count"])
	if err != nil {
		return fmt.Errorf("system_max_count: %w", err)
	}
	sysDisk, err := strconv.Atoi(vars["system_disk_size_gb"])
	if err != nil {
		return fmt.Errorf("system_disk_size_gb: %w", err)
	}
	if err := ValidateNodePoolBounds(sysMin, sysMax, sysDisk, vars["system_disk_type"], vars["system_machine_type"], true); err != nil {
		return fmt.Errorf("system pool: %w", err)
	}

	if retStr, ok := vars["backup_retention_days"]; ok {
		ret, err := strconv.Atoi(retStr)
		if err != nil || ret < 1 || ret > 365 {
			return fmt.Errorf("backup_retention_days must be in [1, 365], got %q", retStr)
		}
	}
	return nil
}
