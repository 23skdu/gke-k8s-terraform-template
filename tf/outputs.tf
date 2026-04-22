output "cluster_name" {
  description = "GKE cluster name"
  value       = google_container_cluster.cluster.name
}

output "cluster_endpoint" {
  description = "GKE cluster endpoint"
  value       = google_container_cluster.cluster.endpoint
  sensitive   = true
}

output "cluster_location" {
  description = "GKE cluster location"
  value       = google_container_cluster.cluster.location
}

output "cluster_id" {
  description = "GKE cluster ID"
  value       = google_container_cluster.cluster.id
}

output "service_account_email" {
  description = "Cluster service account email"
  value       = google_service_account.cluster-sa.email
}