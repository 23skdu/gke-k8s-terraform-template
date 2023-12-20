resource "google_container_cluster" "cluster" {
  name     = "cluster"
  location = var.region
  remove_default_node_pool = true
  initial_node_count       = 1
}
resource "google_container_node_pool" "cluster-main-pool" {
  name       = "main"
  location   = var.region
  cluster    = google_container_cluster.cluster.name
  node_count = 1 
  node_config {
    machine_type = "e2-medium"
    service_account = google_service_account.cluster-sa.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
