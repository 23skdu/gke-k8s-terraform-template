resource "google_service_account" "cluster_sa" {
  account_id   = "${var.cluster_name}-${var.environment}-sa"
  display_name = "GKE Cluster Service Account (${var.environment})"
  description  = "Service account for GKE cluster nodes"
  project      = var.project
}

resource "google_project_iam_member" "cluster_sa_roles" {
  for_each = toset(var.iam_roles)

  project = var.project
  role    = each.value
  member  = "serviceAccount:${google_service_account.cluster_sa.email}"
}

resource "google_service_account_iam_member" "workload_identity" {
  service_account_id = google_service_account.cluster_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project}.svc.id.goog[${var.namespace}/${var.cluster_name}-sa]"
}
