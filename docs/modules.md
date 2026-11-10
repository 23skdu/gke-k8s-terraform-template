# Module Reference

Each module lives under `modules/<name>/` and exposes `main.tf`, `variables.tf`,
and `outputs.tf`. The root module wires them together in `tf/modules.tf`.

## networking

VPC, subnet, firewall rules, Cloud Router, and Cloud NAT.

**Resources:** `google_compute_network`, `google_compute_subnetwork`,
`google_compute_firewall` (x4), `google_compute_router`, `google_compute_router_nat`.

**Key inputs:** `project`, `region`, `environment`, `network_name`,
`subnet_cidr`, `pods_cidr`, `services_cidr`, `master_ipv4_cidr_block`,
`authorized_network_cidr`, `enable_cloud_nat`.

**Outputs:** `network_id`, `network_name`, `subnet_id`, `subnet_name`,
`pods_range_name`, `services_range_name`, `master_ipv4_cidr_block`,
`firewall_rules`, `nat_name`, `router_name`.

## iam

Service account, IAM role bindings, and Workload Identity binding.

**Resources:** `google_service_account`, `google_project_iam_member` (x4),
`google_service_account_iam_member`.

**Key inputs:** `project`, `cluster_name`, `environment`, `namespace`,
`iam_roles`.

**Outputs:** `service_account_email`, `service_account_id`,
`service_account_name`.

## gke

GKE cluster, main and system node pools, and namespaces.

**Resources:** `google_container_cluster`, `google_container_node_pool` (main,
system), `kubernetes_namespace_v1`.

**Key inputs:** `project`, `region`, `cluster_name`, `environment`,
`network_id`, `subnet_id`, `pods_range_name`, `services_range_name`,
`master_ipv4_cidr_block`, `authorized_network_cidrs`, `deletion_protection`,
`service_account_email`, node-pool sizing/disk variables, `namespaces`.

**Outputs:** `cluster_name`, `cluster_endpoint`, `cluster_location`,
`cluster_id`, `cluster_ca_certificate`.

## monitoring

Cloud Logging bucket config, notification channel, and alert policies.

**Resources:** `google_logging_project_bucket_config`,
`google_monitoring_notification_channel`, `google_monitoring_alert_policy` (x4).

**Key inputs:** `project`, `region`, `cluster_name`, `environment`,
`notification_email`, `enable_alerting`.

**Outputs:** `log_bucket_name`, `notification_channel_id`.

## backup

GKE Backup plan with retention.

**Resources:** `google_gke_backup_backup_plan`.

**Key inputs:** `project`, `region`, `cluster_name`, `environment`,
`cluster_id`, `enable_backup`, `backup_schedule`, `backup_retention_days`.

**Outputs:** `backup_plan_id`.

## security

Binary Authorization policy and attestor.

**Resources:** `google_binary_authorization_policy`,
`google_binary_authorization_attestor`, `google_container_analysis_note`.

**Key inputs:** `project`, `environment`, `enable_binary_authorization`,
`attestor_PROJECTS`.

**Outputs:** `policy_id`, `attestor_id`.

## secrets

Secret Manager secrets with IAM bindings.

**Resources:** `google_secret_manager_secret`, `google_secret_manager_secret_iam_member`.

**Key inputs:** `project`, `environment`, `secrets`,
`service_account_email`.

**Outputs:** `secret_ids`.

## statebucket

GCS bucket for Terraform state with versioning and lifecycle rules.

**Resources:** `google_storage_bucket`.

**Key inputs:** `project`, `region`, `bucket_name`, `environment`.

**Outputs:** `bucket_name`, `bucket_id`.

## kubernetes

Kubernetes namespaces.

**Resources:** `kubernetes_namespace_v1`.

**Key inputs:** `namespaces`, `environment`.

**Outputs:** `namespace_names`.
