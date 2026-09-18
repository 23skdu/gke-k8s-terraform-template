# Terraform Tests

Automated infrastructure tests using [Terratest](https://terratest.gruntwork.io/) for validating the GKE Terraform configuration.

## Prerequisites

- [Go](https://go.dev/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads) >= 1.15
- [gcloud CLI](https://cloud.google.com/sdk/docs/install) with authenticated credentials
- A GCP project with billing enabled and the required APIs enabled

## Running Tests

From the project root:

```bash
# Plan-only tests (no GCP resources created)
make test-unit

# Full integration tests (creates real resources)
make test-integration

# All tests
make test
```

Or from the `tests/` directory:

```bash
# Plan-only tests
go test -v -timeout 10m -run "Private|Workload|Datapath|Gateway|Release|Logging|Monitoring|Shielded|AutoUpgrade|Deletion|System|Networking|IAM|Monitoring" ./...

# All tests
go test -v -timeout 30m ./...
```

## Test Suite

| Test | Type | Description |
|------|------|-------------|
| `TestGKEClusterExists` | Apply | Creates cluster and validates name output |
| `TestGKEClusterEndpoint` | Apply | Validates endpoint uses HTTPS |
| `TestGKEClusterLocation` | Apply | Validates location output |
| `TestGKEClusterCA` | Apply | Validates CA certificate output |
| `TestGKEClusterID` | Apply | Validates cluster ID output |
| `TestGKEClusterPrivateNodes` | Plan | Validates private cluster config |
| `TestGKEClusterWorkloadIdentity` | Plan | Validates Workload Identity config |
| `TestGKEClusterAdvancedDatapath` | Plan | Validates ADVANCED_DATAPATH |
| `TestGKEClusterGatewayAPI` | Plan | Validates Gateway API config |
| `TestGKEClusterReleaseChannel` | Plan | Validates release channel |
| `TestGKEClusterLogging` | Plan | Validates logging config |
| `TestGKEClusterMonitoring` | Plan | Validates monitoring config |
| `TestGKENodePoolShieldedConfig` | Plan | Validates shielded node config |
| `TestGKENodePoolAutoUpgrade` | Plan | Validates auto-upgrade management |
| `TestGKEClusterDeletionProtection` | Plan | Validates deletion_protection=false in test |
| `TestSystemNodePoolExists` | Plan | Validates system node pool exists |
| `TestNetworkingVPCExists` | Plan | Validates VPC exists in networking module |
| `TestIAMServiceAccountExists` | Plan | Validates service account in IAM module |
| `TestMonitoringAlertPolicies` | Plan | Validates logging bucket in monitoring module |

## Module-Aware Testing

Tests use `module.<name>.<resource>` prefixed resource addresses since resources are now inside modules. The `GetTerraformOptions` helper sets all required variables for the root module.

## Cleanup

Tests use `defer terraform.Destroy()` for automatic cleanup. If interrupted:

```bash
make destroy ENV=dev
```
