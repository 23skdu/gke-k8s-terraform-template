resource "google_service_account" "cluster_sa" {
  account_id   = "${var.cluster_name}-sa"
  display_name = "GKE Cluster Service Account"
  description  = "Service account for GKE cluster nodes"
}

resource "google_project_iam_member" "cluster_sa_roles" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/stackdriver.resourceMetadata.writer",
  ])

  project = var.project
  role    = each.value
  member  = "serviceAccount:${google_service_account.cluster_sa.email}"
}

resource "google_service_account_iam_member" "workload_identity" {
  service_account_id = google_service_account.cluster_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project}.svc.id.goog[${var.namespace}/${var.cluster_name}-sa]"
}
