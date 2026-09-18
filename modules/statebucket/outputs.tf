output "bucket_name" {
  description = "State bucket name"
  value       = google_storage_bucket.tf_state.name
}

output "bucket_id" {
  description = "State bucket ID"
  value       = google_storage_bucket.tf_state.id
}
