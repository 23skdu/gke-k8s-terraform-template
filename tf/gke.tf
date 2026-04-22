resource "google_container_cluster" "cluster" {
  name     = var.cluster_name
  location = var.region

  enable_autopilot = true

  remove_default_node_pool = true
  initial_node_count     = 1

  network_policy {
    enabled = true
  }

  addon_horizontal_pod_autoscaling {
    enabled = true
  }

  datapath_provider = "ADVANCED_DATAPATH"

  release_channel {
    channel = "REGULAR"
  }

  timeouts {
    create = "30m"
    update = "30m"
  }
}

resource "google_container_node_pool" "main" {
  name       = "main"
  location = var.region
  cluster  = google_container_cluster.cluster.name

  node_count = var.node_count

  node_config {
    machine_type = "e2-medium"

    service_account = google_service_account.cluster-sa.email

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ]

    labels = {
      environment = "production"
    }
  }

  management {
    auto_upgrade = true
    auto_repair = true
  }
}