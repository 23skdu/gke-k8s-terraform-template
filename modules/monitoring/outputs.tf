output "log_bucket_name" {
  description = "Cloud Logging bucket name"
  value       = google_logging_project_bucket_config.default.bucket_id
}

output "notification_channel_id" {
  description = "Notification channel ID (if email configured)"
  value       = var.notification_email != "" ? google_monitoring_notification_channel.email[0].name : ""
}
