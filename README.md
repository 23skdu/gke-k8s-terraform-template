# GKE Kubernetes Terraform Template

Terraform template for provisioning a Google Kubernetes Engine (GKE) cluster on GCP.

## Structure

```
.
├── tf/                      # Terraform configuration
│   ├── main.tf              # Provider and backend configuration
│   ├── variables.tf         # Input variables
│   ├── outputs.tf          # Output values
│   ├── gke.tf              # GKE cluster and node pool
│   ├── iam.tf              # IAM service account
│   ├── gcsbuckets.tf       # GCS bucket for tfstate
│   ├── namespaces.tf       # Kubernetes namespaces
│   ├── logging.tf          # Cloud Logging configuration
│   └── terraform.tfvars.example
├── manifests/              # Kubernetes manifests (optional)
├── LICENSE
└── README.md
```

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.10
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- GCP project with billing enabled

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
   location     = "us-central1"
   node_count   = 1
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

## Clean Up

```bash
terraform destroy
```

## License

See [LICENSE](LICENSE) file.