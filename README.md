# GKE Kubernetes Terraform Template

Terraform template for provisioning a production-ready Google Kubernetes Engine (GKE) cluster on GCP with modular architecture and multi-environment support.

## Features

- **Modular architecture** with reusable `modules/` for GKE, IAM, networking, monitoring, backup, security, and secrets
- **Multi-environment support** with per-environment tfvars in `environments/`
- **GKE Standard cluster** with VPC-native networking
- **Private cluster** with restricted master endpoint access
- **Workload Identity** for secure pod-to-GCP authentication
- **Shielded Nodes** with secure boot and integrity monitoring
- **Datapath Provider** (advanced datapath) for network policy enforcement
- **Gateway API** support for ingress management
- **Managed Prometheus** monitoring with alerting policies
- **Cloud Logging** with system and workload components
- **Dedicated VPC** with firewall rules and Cloud NAT
- **System node pool** with taints for cluster-critical workloads
- **Main node pool** with autoscaling
- **GKE Backup** plans for disaster recovery
- **Binary Authorization** for image verification (optional)
- **Secret Manager** integration with IAM bindings
- **Deletion protection** enabled by default for production safety
- **GCS backend** with versioning for Terraform state
- **CI/CD pipeline** with GitHub Actions
- **Pre-commit hooks** for local quality checks
- **tflint** and **checkov** for static analysis
- **terraform-docs** for auto-generated documentation

## Structure

```
.
├── .github/workflows/ci.yml        # GitHub Actions CI pipeline
├── .pre-commit-config.yaml         # Pre-commit hooks config
├── .tflint.hcl                     # tflint configuration
├── .terraform-docs.yml             # terraform-docs configuration
├── Makefile                        # Common operations
├── modules/                        # Reusable Terraform modules
│   ├── gke/                        # GKE cluster and node pools
│   ├── iam/                        # Service account and IAM roles
│   ├── networking/                 # VPC, subnets, firewall, Cloud NAT
│   ├── monitoring/                 # Logging and alerting
│   ├── backup/                     # GKE backup plans
│   ├── security/                   # Binary Authorization
│   ├── secrets/                    # Secret Manager
│   └── kubernetes/                 # Kubernetes namespaces
├── environments/                   # Per-environment configuration
│   ├── dev/terraform.tfvars
│   ├── staging/terraform.tfvars
│   └── prod/terraform.tfvars
├── tf/                             # Root Terraform configuration
│   ├── main.tf                     # Provider configuration
│   ├── variables.tf                # Input variables
│   ├── outputs.tf                  # Output values
│   └── terraform.tfvars.example
├── tests/                          # Terratest tests (Go)
│   ├── gke_test.go                 # GKE cluster tests
│   ├── iam_test.go                 # IAM tests
│   ├── main_test.go                # Shared test helpers
│   └── README.md
├── LICENSE
└── README.md
```

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.15
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [Go](https://go.dev/doc/install) >= 1.21 (for tests)
- [tflint](https://github.com/terraform-linters/tflint) (for linting)
- [pre-commit](https://pre-commit.com/) (optional)
- GCP project with billing enabled

### Required APIs

```bash
gcloud services enable \
  container.googleapis.com \
  compute.googleapis.com \
  logging.googleapis.com \
  monitoring.googleapis.com \
  iam.googleapis.com \
  backupdr.googleapis.com \
  binaryauthorization.googleapis.com \
  secretmanager.googleapis.com
```

## Quick Start

```bash
# 1. Authenticate
gcloud auth application-default login

# 2. Initialize
make init ENV=dev

# 3. Plan
make plan ENV=dev

# 4. Apply
make apply ENV=dev

# 5. Get credentials
gcloud container clusters get-credentials gke-dev --region=us-central1 --project=YOUR_PROJECT
```

## Multi-Environment Deployment

Each environment has its own `terraform.tfvars` in `environments/`:

| Environment | Cluster | Deletion Protection | Autoscaling | Backup |
|-------------|---------|---------------------|-------------|--------|
| `dev` | `gke-dev` | off | 1-3 nodes | off |
| `staging` | `gke-staging` | off | 1-5 nodes | 14 days |
| `prod` | `gke-prod` | **on** | 2-10 nodes | 30 days |

Deploy to a specific environment:
```bash
make plan ENV=staging
make apply ENV=staging
```

## Modules

| Module | Description |
|--------|-------------|
| `networking` | VPC, subnets, firewall rules, Cloud NAT |
| `iam` | Service account, IAM roles, Workload Identity |
| `gke` | GKE cluster, main + system node pools, namespaces |
| `monitoring` | Cloud Logging, alerting policies, notification channels |
| `backup` | GKE Backup plans |
| `security` | Binary Authorization policies |
| `secrets` | Secret Manager secrets with IAM bindings |
| `kubernetes` | Kubernetes namespaces |

## Testing

```bash
# Run plan-only tests (no GCP resources created)
make test-unit

# Run full integration tests (creates real GCP resources)
make test-integration

# Run all tests
make test
```

See [tests/README.md](tests/README.md) for details.

## Linting & Validation

```bash
# Format
make fmt

# Validate
make validate

# tflint
make lint

# Checkov
make checkov
```

## Pre-commit Hooks

```bash
pip install pre-commit
pre-commit install
```

## CI/CD

GitHub Actions runs on push/PR to `main`:
- Terraform format, init, validate (root + modules)
- tflint static analysis
- Checkov security scanning
- Terratest plan-only tests

## Clean Up

```bash
make destroy ENV=dev
```

## License

See [LICENSE](LICENSE) file.
