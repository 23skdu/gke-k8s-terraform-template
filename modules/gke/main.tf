resource "google_container_cluster" "cluster" {
  name     = "${var.cluster_name}-${var.environment}"
  location = var.region
  project  = var.project

  remove_default_node_pool = true
  initial_node_count       = 1

  network    = var.network_id
  subnetwork = var.subnet_id

  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_range_name
    services_secondary_range_name = var.services_range_name
  }

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }

  master_authorized_networks_config {
    dynamic "cidr_blocks" {
      for_each = var.authorized_network_cidrs
      content {
        cidr_block   = cidr_blocks.value
        display_name = "Authorized network ${cidr_blocks.key}"
      }
    }
  }

  workload_identity_config {
    workload_pool = "${var.project}.svc.id.goog"
  }

  datapath_provider = "ADVANCED_DATAPATH"

  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  release_channel {
    channel = "REGULAR"
  }

  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }
    http_load_balancing {
      disabled = false
    }
    network_policy_config {
      disabled = false
    }
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }

  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }

  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS"]
    managed_prometheus {
      enabled = true
    }
  }

  timeouts {
    create = "30m"
    update = "30m"
  }

  deletion_protection = var.deletion_protection
}

# Main workload node pool
resource "google_container_node_pool" "main" {
  name     = "main"
  location = var.region
  cluster  = google_container_cluster.cluster.name
  project  = var.project

  autoscaling {
    min_node_count = var.main_min_count
    max_node_count = var.main_max_count
  }

  node_config {
    machine_type = var.main_machine_type
    disk_size_gb = var.main_disk_size_gb
    disk_type    = var.main_disk_type

    service_account = var.service_account_email

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    labels = {
      environment = var.environment
      pool        = "main"
    }

    tags = ["gke-node"]
  }

  management {
    auto_upgrade = true
    auto_repair  = true
  }

  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
  }
}

# System node pool for cluster-critical workloads
resource "google_container_node_pool" "system" {
  name     = "system"
  location = var.region
  cluster  = google_container_cluster.cluster.name
  project  = var.project

  autoscaling {
    min_node_count = var.system_min_count
    max_node_count = var.system_max_count
  }

  node_config {
    machine_type = var.system_machine_type
    disk_size_gb = var.system_disk_size_gb
    disk_type    = var.system_disk_type

    service_account = var.service_account_email

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    labels = {
      environment = var.environment
      pool        = "system"
    }

    taint {
      key    = "node-pool"
      value  = "system"
      effect = "PREFER_NO_SCHEDULE"
    }

    tags = ["gke-node"]
  }

  management {
    auto_upgrade = true
    auto_repair  = true
  }

  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
  }
}
