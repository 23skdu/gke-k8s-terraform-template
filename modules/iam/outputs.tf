output "service_account_email" {
  description = "Cluster service account email"
  value       = google_service_account.cluster_sa.email
}

output "service_account_id" {
  description = "Cluster service account ID"
  value       = google_service_account.cluster_sa.account_id
}

output "service_account_name" {
  description = "Cluster service account fully qualified name"
  value       = google_service_account.cluster_sa.name
}
