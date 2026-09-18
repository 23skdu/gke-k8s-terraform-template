# GKE Kubernetes Terraform Template

Terraform template for provisioning a Google Kubernetes Engine (GKE) cluster on GCP with security and production best practices.

## Features

- **GKE Standard cluster** with VPC-native networking
- **Private cluster** with restricted master endpoint access
- **Workload Identity** for secure pod-to-GCP authentication
- **Shielded Nodes** with secure boot and integrity monitoring
- **Datapath Provider** (advanced datapath) for network policy enforcement
- **Gateway API** support for ingress management
- **Managed Prometheus** monitoring
- **Cloud Logging** with system and workload components
- **Auto-upgrade and auto-repair** node management
- **GCS backend** with versioning for Terraform state

## Structure

```
.
├── tf/                         # Terraform configuration
│   ├── main.tf                 # Provider and backend configuration
│   ├── variables.tf            # Input variables
│   ├── outputs.tf              # Output values
│   ├── gke.tf                  # GKE cluster and node pool
│   ├── iam.tf                  # IAM service account and Workload Identity
│   ├── gcsbuckets.tf           # GCS bucket for tfstate
│   ├── namespaces.tf           # Kubernetes namespaces
│   ├── logging.tf              # Cloud Logging configuration
│   └── terraform.tfvars.example
├── tests/                      # Terratest integration tests (Go)
│   ├── gke_test.go             # GKE cluster validation tests
│   ├── iam_test.go             # IAM and service account tests
│   ├── kubernetes_test.go      # Kubernetes namespace tests
│   ├── logging_test.go         # Cloud Logging tests
│   ├── storage_test.go         # GCS bucket tests
│   ├── main_test.go            # Shared test helpers
│   ├── go.mod                  # Go module dependencies
│   └── README.md               # Test instructions
├── LICENSE
└── README.md
```

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.15
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- GCP project with billing enabled
- A VPC network with two secondary IP ranges (for pods and services)

### Required APIs

Enable the following APIs in your GCP project:

```bash
gcloud services enable \
  container.googleapis.com \
  compute.googleapis.com \
  logging.googleapis.com \
  monitoring.googleapis.com \
  iam.googleapis.com
```

## Usage

1. Authenticate with GCP:
   ```bash
   gcloud auth application-default login
   gcloud auth login
   ```

2. Set your project:
   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

3. Copy the example variables file:
   ```bash
   cp tf/terraform.tfvars.example tf/terraform.tfvars
   ```

4. Edit `tf/terraform.tfvars` with your values:
   ```hcl
   project      = "my-gcp-project"
   region       = "us-central1"
   cluster_name = "gke-cluster"

   # Networking (must match your VPC/subnet configuration)
   network                = "default"
   subnetwork             = "default"
   pods_range_name        = "pods"
   services_range_name    = "services"
   master_ipv4_cidr_block = "172.16.0.0/28"
   authorized_network_cidr = "0.0.0.0/0"

   # Node pool
   node_count   = 1
   machine_type = "e2-medium"
   environment  = "production"

   # Workload Identity
   namespace = "default"
   ```

5. Initialize Terraform:
   ```bash
   cd tf
   terraform init
   ```

6. Plan and apply:
   ```bash
   terraform plan -out=tfplan
   terraform apply tfplan
   ```

## Getting Cluster Credentials

After provisioning, get credentials:
```bash
gcloud container clusters get-credentials CLUSTER_NAME --region=REGION --project=PROJECT
```

## Testing

This project includes automated tests using [Terratest](https://terratest.gruntwork.io/) to validate the Terraform configuration.

```bash
# Run all tests
cd tests && go test -v -timeout 30m ./...

# Run a specific test
cd tests && go test -v -run TestGKEClusterExists -timeout 30m ./...
```

See [tests/README.md](tests/README.md) for full test documentation including prerequisites, troubleshooting, and available test suites.

## Clean Up

```bash
terraform destroy
```

## License

See [LICENSE](LICENSE) file.
