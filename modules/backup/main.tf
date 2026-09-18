resource "google_gke_backup_backup_plan" "cluster_backup" {
  count    = var.enable_backup ? 1 : 0
  name     = "${var.cluster_name}-${var.environment}-backup"
  project  = var.project
  location = var.region
  cluster  = var.cluster_id

  retention_policy {
    backup_delete_lock_days = var.backup_retention_days
    backup_retain_days      = var.backup_retention_days
  }

  backup_config {
    include_volume_data = true
    include_secrets     = true
    all_namespaces     = true
  }
}
