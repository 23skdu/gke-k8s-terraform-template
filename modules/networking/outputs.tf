output "network_id" {
  description = "VPC network ID"
  value       = google_compute_network.vpc.id
}

output "network_name" {
  description = "VPC network name"
  value       = google_compute_network.vpc.name
}

output "subnet_id" {
  description = "Subnet ID"
  value       = google_compute_subnetwork.gke.id
}

output "subnet_name" {
  description = "Subnet name"
  value       = google_compute_subnetwork.gke.name
}

output "pods_range_name" {
  description = "Pods secondary range name"
  value       = var.pods_cidr_name
}

output "services_range_name" {
  description = "Services secondary range name"
  value       = var.services_cidr_name
}

output "master_ipv4_cidr_block" {
  description = "Master CIDR block"
  value       = var.master_ipv4_cidr_block
}
