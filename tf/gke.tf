resource "google_container_cluster" "cluster" {
  name     = var.cluster_name
  location = var.region

  # We manage our own node pool
  remove_default_node_pool = true
  initial_node_count       = 1

  # Networking - VPC-native cluster
  network    = var.network
  subnetwork = var.subnetwork

  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_range_name
    services_secondary_range_name = var.services_range_name
  }

  # Private cluster configuration
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = var.authorized_network_cidr
      display_name = "Authorized network"
    }
  }

  # Workload Identity
  workload_identity_config {
    workload_pool = "${var.project}.svc.id.goog"
  }

  # Datapath Provider (replaces legacy network policy)
  datapath_provider = "ADVANCED_DATAPATH"

  # Gateway API (replaces Ingress controller)
  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  # Release channel
  release_channel {
    channel = "REGULAR"
  }

  # Addons
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
  }

  # Security
  # Node security config (Shielded Nodes default on for new clusters)

  # Cluster logging
  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }

  # Cluster monitoring
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

  deletion_protection = false
}

resource "google_container_node_pool" "main" {
  name     = "main"
  location = var.region
  cluster  = google_container_cluster.cluster.name

  initial_node_count = var.node_count

  node_config {
    machine_type = var.machine_type

    service_account = google_service_account.cluster_sa.email

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    labels = {
      environment = var.environment
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
