# Terraform Tests

Automated infrastructure tests using [Terratest](https://terratest.gruntwork.io/) for validating the GKE Terraform configuration.

## Prerequisites

- [Go](https://go.dev/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads) >= 1.15
- [gcloud CLI](https://cloud.google.com/sdk/docs/install) with authenticated credentials
- A GCP project with billing enabled and the required APIs enabled

### Required GCP APIs

```bash
gcloud services enable \
  container.googleapis.com \
  compute.googleapis.com \
  logging.googleapis.com \
  monitoring.googleapis.com \
  iam.googleapis.com
```

## Setup

1. Authenticate with GCP:
   ```bash
   gcloud auth application-default login
   gcloud auth login
   ```

2. Set your project:
   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

## Running Tests

From the `tests/` directory:

```bash
# Run all tests
go test -v -timeout 30m ./...

# Run a specific test
go test -v -run TestGKEClusterExists -timeout 30m ./...

# Run plan-only tests (no GCP resources created)
go test -v -timeout 10m -run "Plan|Private|Workload|Datapath|Gateway|Release|Logging|Monitoring|Shielded|AutoUpgrade|Deletion|ServiceAccount|IAMRoles|WorkloadIdentity" ./...
```

Or use the Makefile from the project root:

```bash
make test-unit           # Plan-only tests
make test-integration    # Full tests (creates real resources)
make test                # All tests
```

### Useful Flags

| Flag | Description |
|------|-------------|
| `-v` | Verbose output |
| `-run <pattern>` | Run only tests matching regex |
| `-timeout <duration>` | Test timeout (default: 10m, recommended: 30m+) |
| `-parallel <n>` | Concurrent test count |
| `-count=1` | Disable test caching |

## Test Suite

| File | Tests |
|------|-------|
| `gke_test.go` | Cluster creation, endpoint, location, CA cert, private nodes, workload identity, datapath, gateway API, release channel, logging, monitoring, node pool config, deletion protection |
| `iam_test.go` | Service account, email output, IAM roles, workload identity binding |
| `main_test.go` | Shared `GetTerraformOptions` helper with module-compatible variables |

## Test Types

- **Plan-only tests** (`InitAndPlanAndShowWithStruct`): Validate config without creating infrastructure. Fast, free, no credentials needed.
- **Apply tests** (`InitAndApply`): Create real GCP resources. Slower, costs money.

## Module-Aware Testing

Tests pass variables matching the new module-based root `tf/` configuration. The `GetTerraformOptions` helper in `main_test.go` sets all required variables including networking, node pools, and feature flags.

## Cleanup

Tests use `defer terraform.Destroy()` to clean up automatically. If interrupted:

```bash
make destroy ENV=dev
```

## Cost Considerations

Apply tests provision real GKE clusters. Expect GCP charges. Run plan-only tests (`make test-unit`) first.

## Troubleshooting

- **Timeout errors**: Increase with `-timeout 30m` or higher.
- **Auth errors**: Run `gcloud auth application-default login`.
- **Permission errors**: Account needs Editor/Owner on the target project.
- **Quota errors**: Check GCP quotas for compute, IPs, and GKE clusters.
