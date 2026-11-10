# GKE Kubernetes Terraform Template

[![CI](https://github.com/23skdu/gke-k8s-terraform-template/actions/workflows/ci.yml/badge.svg)](https://github.com/23skdu/gke-k8s-terraform-template/actions/workflows/ci.yml)

Production-ready **GKE Standard** cluster on Google Cloud, built with a modular
Terraform architecture, multi-environment support, and a hermetic Terratest
suite. Nine reusable modules cover networking, identity, cluster, monitoring,
backup, security, secrets, state, and Kubernetes resources.

## Documentation

Extended documentation lives in [`docs/`](docs/):

- [Architecture](docs/architecture.md) — design, module graph, networking model
- [Modules](docs/modules.md) — per-module reference
- [Variables](docs/variables.md) — root variable reference
- [Environments](docs/environments.md) — dev / staging / prod matrix
- [Deployment](docs/deployment.md) — deploy and teardown steps
- [Operations](docs/operations.md) — day-2 ops, Makefile, troubleshooting
- [Testing](docs/testing.md) — Terratest plan-only and integration tests
- [Security](docs/security.md) — security posture
- [CI/CD](docs/ci-cd.md) — pipelines and pre-commit hooks

## Features

- **9 reusable modules**: networking, iam, gke, monitoring, backup, security, secrets, statebucket, kubernetes
- **Multi-environment** with per-environment tfvars (`dev`, `staging`, `prod`)
- **GKE Standard cluster** with VPC-native networking and private nodes
- **Workload Identity** for secure pod-to-GCP authentication
- **Shielded Nodes** with secure boot and integrity monitoring
- **ADVANCED_DATAPATH** for network policy enforcement
- **Gateway API** for ingress management
- **Managed Prometheus** with alerting policies (CPU, memory, cluster health, auto-repair)
- **Cloud Logging** with system and workload components
- **Dedicated VPC** with 4 firewall rules and Cloud NAT
- **Two node pools**: main (autoscaling) + system (tainted for critical workloads)
- **Configurable disk** size and type per node pool
- **GKE Backup** plans for disaster recovery
- **Binary Authorization** for image verification (optional)
- **Secret Manager** integration with IAM bindings
- **GCS state bucket** with versioning and lifecycle rules
- **Input validation** on all variables (CIDR format, name patterns, email, ranges)
- **Deletion protection** enabled by default
- **CI/CD** via GitHub Actions (fmt, validate, tflint, checkov, terratest)
- **Pre-commit hooks** for local quality checks
- **Hermetic plan-only tests** that need no GCP credentials

## Structure

```
.
├── .github/workflows/ci.yml        # CI pipeline
├── .pre-commit-config.yaml         # Pre-commit hooks
├── .tflint.hcl                     # tflint config
├── .terraform-docs.yml             # terraform-docs config
├── CONTRIBUTING.md                 # Development guide
├── Makefile                        # Common operations
├── docs/                           # Extended documentation
├── modules/                        # Reusable Terraform modules
│   ├── networking/                 # VPC, subnets, firewall, Cloud NAT
│   ├── iam/                        # Service account, IAM roles, Workload Identity
│   ├── gke/                        # Cluster, main + system node pools
│   ├── monitoring/                 # Logging, alerting policies
│   ├── backup/                     # GKE backup plans
│   ├── security/                   # Binary Authorization
│   ├── secrets/                    # Secret Manager
│   ├── statebucket/                # GCS state bucket
│   └── kubernetes/                 # Namespaces
├── environments/                   # Per-environment tfvars
│   ├── dev/terraform.tfvars
│   ├── staging/terraform.tfvars
│   └── prod/terraform.tfvars
├── tf/                             # Root Terraform configuration
│   ├── main.tf                     # Provider config
│   ├── variables.tf                # Input variables (with validation)
│   ├── modules.tf                  # Module calls
│   ├── outputs.tf                  # Output values
│   ├── backend.tf                  # GCS backend declaration
│   └── terraform.tfvars.example
├── tests/                          # Terratest tests (Go)
│   ├── gke_test.go
│   ├── main_test.go                # Shared helpers + hermetic fixture
│   ├── modules_plan_test.go        # Plan-only module coverage
│   ├── oracle_gate_test.go         # Test-suite quality gate
│   ├── pkg/oracle/                 # Pure-Go validation oracles
│   ├── go.mod / go.sum
│   └── README.md
├── LICENSE
└── README.md
```

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.15
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [Go](https://go.dev/doc/install) >= 1.21 (for tests)
- [tflint](https://github.com/terraform-linters/tflint) (for linting)
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

# 2. Copy and edit tfvars
cp tf/terraform.tfvars.example environments/dev/terraform.tfvars
# Edit environments/dev/terraform.tfvars with your project ID

# 3. Deploy
make init ENV=dev
make plan ENV=dev
make apply ENV=dev

# 4. Get credentials
gcloud container clusters get-credentials gke-dev --region=us-central1 --project=YOUR_PROJECT
```

## Multi-Environment Deployment

| Environment | Cluster | Deletion Protection | Main Pool | System Pool | Backup | Alerts |
|-------------|---------|---------------------|-----------|-------------|--------|--------|
| `dev` | `gke-dev` | off | 1-3x e2-medium | 1-2x e2-small | off | off |
| `staging` | `gke-staging` | off | 1-5x e2-medium | 1-2x e2-medium | 14d | on |
| `prod` | `gke-prod` | **on** | 2-10x e2-standard-4 | 2-5x e2-medium | 30d | on |

```bash
make plan ENV=staging
make apply ENV=staging
```

See [Environments](docs/environments.md) for the full matrix.

## Modules

| Module | Resources | Key Outputs |
|--------|-----------|-------------|
| `networking` | VPC, subnet, 4 firewall rules, Cloud Router, Cloud NAT | `network_id`, `subnet_id`, `firewall_rules`, `nat_name` |
| `iam` | Service account, 4 IAM roles, Workload Identity binding | `service_account_email` |
| `gke` | Cluster, main node pool, system node pool, namespaces | `cluster_endpoint`, `cluster_id`, `cluster_ca_certificate` |
| `monitoring` | Log bucket, email channel, 4 alert policies | `log_bucket_name`, `notification_channel_id` |
| `backup` | GKE Backup plan with retention | `backup_plan_id` |
| `security` | Binary Authorization policy + attestor | `policy_id` |
| `secrets` | Secret Manager secrets with IAM | `secret_ids` |
| `statebucket` | GCS bucket with versioning | `bucket_name` |
| `kubernetes` | Kubernetes namespaces | `namespace_names` |

See [Modules](docs/modules.md) for the full reference.

## Variables

All variables include `validation` blocks. See [`tf/variables.tf`](tf/variables.tf)
and [Variables](docs/variables.md) for the full list.

## Testing

```bash
make test-unit           # Plan-only tests (no GCP resources)
make test-integration    # Full tests (creates real resources)
make test                # All tests
```

The suite includes plan-only coverage for every module plus an oracle quality
gate. See [Testing](docs/testing.md) and [tests/README.md](tests/README.md).

## Linting & Validation

```bash
make fmt         # Format all .tf files
make validate    # Validate root + all modules
make lint        # Run tflint
make checkov     # Security scan
```

## Pre-commit Hooks

```bash
pip install pre-commit
pre-commit install
```

## CI/CD

GitHub Actions on push/PR to `main`: format check, init, validate (root + modules),
tflint, checkov, and Terratest plan-only tests. See [CI/CD](docs/ci-cd.md).

## Clean Up

```bash
make destroy ENV=dev
```

## License

See [LICENSE](LICENSE) file.
