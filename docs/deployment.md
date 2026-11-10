# Deployment

## Prerequisites

- Terraform >= 1.15
- gcloud CLI authenticated (`gcloud auth application-default login`)
- A GCP project with billing enabled
- Required APIs enabled (see root README)

## Deploy

```bash
# 1. Copy and edit tfvars
cp tf/terraform.tfvars.example environments/dev/terraform.tfvars
# edit project ID, region, etc.

# 2. Initialize
make init ENV=dev

# 3. Plan
make plan ENV=dev

# 4. Apply
make apply ENV=dev

# 5. Fetch credentials
gcloud container clusters get-credentials gke-dev --region=us-central1 --project=YOUR_PROJECT
```

## Teardown

```bash
make destroy ENV=dev
```

## State Bucket Bootstrap

The `statebucket` module creates the GCS bucket used as the Terraform backend.
On the very first deployment you may need to apply with a local backend first,
then migrate state into GCS. See `tf/backend.tf` and `tf/main.tf`.
