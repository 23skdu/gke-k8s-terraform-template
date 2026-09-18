output "backup_plan_id" {
  description = "Backup plan ID"
  value       = var.enable_backup ? google_gke_backup_backup_plan.cluster_backup[0].id : ""
}
