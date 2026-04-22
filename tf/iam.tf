resource "google_service_account" "cluster_sa" {
  account_id   = "${var.cluster_name}-sa"
  display_name = "GKE Cluster Service Account"
  description = "Service account for GKE cluster nodes"
}

resource "google_project_iam_member" "cluster_sa_roles" {
  for_each = toset([
    "roles/container.nodeSelector",
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles monitoring.viewer",
  ])

  project = var.project
  role    = each.value
  member  = "serviceAccount:${google_service_account.cluster_sa.email}"
}