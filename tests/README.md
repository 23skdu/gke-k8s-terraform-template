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

3. Copy the example variables file and edit with your values:
   ```bash
   cp tf/terraform.tfvars.example tf/terraform.tfvars
   ```

## Running Tests

From the `tests/` directory:

```bash
# Run all tests
go test -v ./...

# Run a specific test
go test -v -run TestGKEClusterExists ./...

# Run tests in parallel (default behavior)
go test -v -parallel 4 ./...

# Run with a timeout (tests can take a while due to GCP resource creation)
go test -v -timeout 30m ./...
```

### Useful Flags

| Flag | Description |
|------|-------------|
| `-v` | Verbose output showing test progress |
| `-run <pattern>` | Run only tests matching the regex pattern |
| `-timeout <duration>` | Set test timeout (default: 10m, recommended: 30m+) |
| `-parallel <n>` | Number of tests to run concurrently |
| `-count=1` | Disable test caching |

## Test Suite

Tests are organized by resource type:

| File | Tests |
|------|-------|
| `gke_test.go` | GKE cluster creation, endpoint, location, CA certificate, private nodes, workload identity, datapath, gateway API, release channel, logging, monitoring, node pool config, deletion protection |
| `iam_test.go` | Service account creation, email output, IAM role assignments, workload identity binding |
| `kubernetes_test.go` | Kubernetes namespace validation |
| `logging_test.go` | Cloud Logging sink configuration |
| `storage_test.go` | GCS bucket for Terraform state |
| `main_test.go` | Shared `GetTerraformOptions` helper and test setup |

## Test Types

- **Plan-only tests** (`InitAndPlanAndShowWithStruct`): Validate resource configuration without creating real infrastructure. These are fast and safe.
- **Apply tests** (`InitAndApply`): Create real GCP resources and validate outputs. These are slower and incur costs.

## Cleanup

Tests use `defer terraform.Destroy()` to clean up resources automatically after each test. If tests are interrupted, manually destroy:

```bash
cd tf
terraform destroy
```

## Cost Considerations

Apply tests provision real GKE clusters and associated resources. Expect GCP charges during test runs. Run plan-only tests first for quick validation.

## Troubleshooting

- **Timeout errors**: Increase the timeout with `-timeout 30m` or higher.
- **Authentication errors**: Ensure `gcloud auth application-default login` is run and credentials are valid.
- **Permission errors**: The authenticated account needs sufficient IAM roles (e.g., Editor or Owner) on the target GCP project.
- **Quota errors**: Check GCP project quotas for compute resources, IPs, and GKE clusters.
