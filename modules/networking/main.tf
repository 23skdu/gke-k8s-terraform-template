resource "google_compute_network" "vpc" {
  name                    = "${var.network_name}-${var.environment}"
  auto_create_subnetworks = false
  project                 = var.project
}

resource "google_compute_subnetwork" "gke" {
  name                     = "${var.subnet_name}-${var.environment}"
  project                  = var.project
  region                   = var.region
  network                  = google_compute_network.vpc.id
  ip_cidr_range            = var.subnet_cidr
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = var.pods_cidr_name
    ip_cidr_range = var.pods_cidr
  }

  secondary_ip_range {
    range_name    = var.services_cidr_name
    ip_cidr_range = var.services_cidr
  }
}

# Firewall: allow internal communication between nodes
resource "google_compute_firewall" "allow_internal" {
  name    = "${var.network_name}-${var.environment}-allow-internal"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = [
    var.subnet_cidr,
    var.pods_cidr,
    var.services_cidr,
  ]
}

# Firewall: allow health checks from Google load balancers
resource "google_compute_firewall" "allow_health_checks" {
  name    = "${var.network_name}-${var.environment}-allow-health-checks"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  source_ranges = [
    "130.211.0.0/22",
    "35.191.0.0/16",
  ]

  target_tags = ["gke-node"]
}

# Firewall: allow master to communicate with nodes
resource "google_compute_firewall" "allow_master_to_nodes" {
  name    = "${var.network_name}-${var.environment}-allow-master"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["443", "10250"]
  }

  source_ranges = [var.master_ipv4_cidr_block]
  target_tags   = ["gke-node"]
}

# Firewall: allow nodes to communicate with master API
resource "google_compute_firewall" "allow_nodes_to_master" {
  name    = "${var.network_name}-${var.environment}-allow-nodes-master"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["443", "10250"]
  }

  source_ranges = [var.subnet_cidr]
  destination_ranges = [var.master_ipv4_cidr_block]
  target_tags   = ["gke-node"]
}

# Cloud Router for NAT
resource "google_compute_router" "router" {
  count   = var.enable_cloud_nat ? 1 : 0
  name    = "${var.network_name}-${var.environment}-router"
  project = var.project
  region  = var.region
  network = google_compute_network.vpc.id
}

# Cloud NAT for egress traffic
resource "google_compute_router_nat" "nat" {
  count                              = var.enable_cloud_nat ? 1 : 0
  name                               = "${var.network_name}-${var.environment}-nat"
  project                            = var.project
  router                             = google_compute_router.router[0].name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
