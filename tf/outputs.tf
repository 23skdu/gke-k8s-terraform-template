output "cluster_name" {
  description = "GKE cluster name"
  value       = module.gke.cluster_name
}

output "cluster_endpoint" {
  description = "GKE cluster endpoint"
  value       = module.gke.cluster_endpoint
  sensitive   = true
}

output "cluster_location" {
  description = "GKE cluster location"
  value       = module.gke.cluster_location
}

output "cluster_id" {
  description = "GKE cluster ID"
  value       = module.gke.cluster_id
}

output "cluster_ca_certificate" {
  description = "Base64 encoded CA certificate"
  value       = module.gke.cluster_ca_certificate
  sensitive   = true
}

output "service_account_email" {
  description = "Cluster service account email"
  value       = module.iam.service_account_email
}

output "network_name" {
  description = "VPC network name"
  value       = module.networking.network_name
}

output "subnet_name" {
  description = "Subnet name"
  value       = module.networking.subnet_name
}

output "backup_plan_id" {
  description = "Backup plan ID"
  value       = module.backup.backup_plan_id
}
