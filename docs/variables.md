# Variable Reference (Root Module)

All variables are defined in [`tf/variables.tf`](../tf/variables.tf) and include
`validation` blocks. Environment-specific values live in
`environments/<env>/terraform.tfvars`.

## Core

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `project` | string | (required) | GCP project ID (6-30 chars, lowercase/digits/hyphens) |
| `region` | string | `us-central1` | GCP region |
| `environment` | string | `production` | One of: dev, staging, production, test |
| `cluster_name` | string | `gke-cluster` | GKE cluster name |
| `deletion_protection` | bool | `true` | Enable cluster deletion protection |

## Networking

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `network_name` | string | `gke-network` | VPC network name |
| `subnet_cidr` | string | `10.0.0.0/20` | Primary subnet CIDR |
| `pods_cidr` | string | `10.4.0.0/14` | Secondary CIDR for pods |
| `services_cidr` | string | `10.8.0.0/20` | Secondary CIDR for services |
| `master_ipv4_cidr_block` | string | `172.16.0.0/28` | Master private endpoint CIDR |
| `authorized_network_cidr` | string | `0.0.0.0/0` | Master authorized network |

## Node pools

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `main_machine_type` | string | `e2-medium` | Main pool machine type |
| `main_min_count` | number | `1` | Main pool min nodes |
| `main_max_count` | number | `5` | Main pool max nodes |
| `main_disk_size_gb` | number | `100` | Main pool disk size (GB) |
| `main_disk_type` | string | `pd-balanced` | Main pool disk type |
| `system_machine_type` | string | `e2-medium` | System pool machine type |
| `system_min_count` | number | `1` | System pool min nodes |
| `system_max_count` | number | `3` | System pool max nodes |
| `system_disk_size_gb` | number | `100` | System pool disk size (GB) |
| `system_disk_type` | string | `pd-balanced` | System pool disk type |

## IAM / Kubernetes

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `namespace` | string | `default` | Workload Identity namespace |
| `namespaces` | list(string) | `["default","prod"]` | Namespaces to create |

## Monitoring / Backup / Security / Secrets

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `notification_email` | string | `""` | Alert notification email |
| `enable_alerting` | bool | `true` | Enable alert policies |
| `enable_backup` | bool | `true` | Enable GKE backup plan |
| `backup_retention_days` | number | `30` | Backup retention (days) |
| `enable_binary_authorization` | bool | `false` | Enable Binary Authorization |
| `secrets` | map(string) | `{}` | Secret Manager secrets to create |

## State bucket

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `state_bucket_name` | string | `tf-state` | GCS state bucket name |
