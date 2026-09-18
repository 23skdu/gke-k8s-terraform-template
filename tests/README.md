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

## Running Tests

From the project root:

```bash
# Plan-only tests (no GCP resources created)
make test-unit

# Full integration tests (creates real resources, costs money)
make test-integration

# All tests
make test
```

Or from the `tests/` directory:

```bash
# Plan-only tests
go test -v -timeout 10m -run "Private|Workload|Datapath|Gateway|Release|Logging|Monitoring|Shielded|AutoUpgrade|Deletion|System|Networking|IAM|Alert" ./...

# All tests
go test -v -timeout 30m ./...
```

### Useful Flags

| Flag | Description |
|------|-------------|
| `-v` | Verbose output |
| `-run <pattern>` | Run only tests matching regex |
| `-timeout <duration>` | Test timeout (default: 10m, recommended: 30m+) |
| `-count=1` | Disable test caching |

## Test Suite

### Cluster Tests (`TestGKECluster*`)

| Test | Type | What it validates |
|------|------|-------------------|
| `TestGKEClusterExists` | Apply | Creates cluster, validates name output |
| `TestGKEClusterEndpoint` | Apply | Endpoint uses HTTPS |
| `TestGKEClusterLocation` | Apply | Location output is non-empty |
| `TestGKEClusterCA` | Apply | CA certificate output is non-empty |
| `TestGKEClusterID` | Apply | Cluster ID output is non-empty |
| `TestGKEClusterPrivateNodes` | Plan | `private_cluster_config` exists |
| `TestGKEClusterWorkloadIdentity` | Plan | `workload_identity_config` exists |
| `TestGKEClusterAdvancedDatapath` | Plan | `datapath_provider` is `ADVANCED_DATAPATH` |
| `TestGKEClusterGatewayAPI` | Plan | `gateway_api_config` exists |
| `TestGKEClusterReleaseChannel` | Plan | `release_channel` exists |
| `TestGKEClusterLogging` | Plan | `logging_config` exists |
| `TestGKEClusterMonitoring` | Plan | `monitoring_config` exists |
| `TestGKEClusterDeletionProtection` | Plan | `deletion_protection` is `false` in test |

### Node Pool Tests (`TestGKE*NodePool*`)

| Test | Type | What it validates |
|------|------|-------------------|
| `TestGKENodePoolShieldedConfig` | Plan | `node_config` exists in main pool |
| `TestGKENodePoolAutoUpgrade` | Plan | `management` exists in main pool |
| `TestSystemNodePoolExists` | Plan | System node pool exists |

### Module Tests

| Test | Type | What it validates |
|------|------|-------------------|
| `TestNetworkingVPCExists` | Plan | VPC exists in networking module |
| `TestIAMServiceAccountExists` | Plan | Service account exists in IAM module |
| `TestMonitoringAlertPolicies` | Plan | Logging bucket exists in monitoring module |

## Test Configuration

Tests use `GetTerraformOptions()` in `gke_test.go` which sets all required variables for the root module including:
- Project, region, environment
- Networking CIDRs
- Node pool machine types, counts, disk config
- Feature flags (alerting, backup, binary auth)
- State bucket name

## Module-Aware Testing

Tests use `module.<name>.<resource>` prefixed resource addresses since resources are inside modules:
```go
plan.ResourcePlannedValuesMap["module.gke.google_container_cluster.cluster"]
plan.ResourcePlannedValuesMap["module.gke.google_container_node_pool.main"]
plan.ResourcePlannedValuesMap["module.networking.google_compute_network.vpc"]
plan.ResourcePlannedValuesMap["module.iam.google_service_account.cluster_sa"]
plan.ResourcePlannedValuesMap["module.monitoring.google_logging_project_bucket_config.default"]
```

## Cleanup

Tests use `defer terraform.Destroy()` for automatic cleanup. If interrupted:

```bash
make destroy ENV=dev
```

## Cost Considerations

- **Plan-only tests**: Free, no GCP resources created
- **Apply tests**: Create real GKE clusters and associated resources. Expect GCP charges.
- Run `make test-unit` for quick validation during development
