variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
  default     = "gke-cluster"
}

variable "environment" {
  description = "Environment label for the cluster nodes"
  type        = string
  default     = "production"
}

variable "namespace" {
  description = "Kubernetes namespace for Workload Identity binding"
  type        = string
  default     = "default"
}

variable "machine_type" {
  description = "Machine type for the node pool"
  type        = string
  default     = "e2-medium"
}

# Networking

variable "network" {
  description = "VPC network name"
  type        = string
  default     = "default"
}

variable "subnetwork" {
  description = "VPC subnetwork name"
  type        = string
  default     = "default"
}

variable "pods_range_name" {
  description = "Name of the secondary IP range for pods"
  type        = string
  default     = "pods"
}

variable "services_range_name" {
  description = "Name of the secondary IP range for services"
  type        = string
  default     = "services"
}

# Private cluster

variable "master_ipv4_cidr_block" {
  description = "CIDR block for the GKE master's private endpoint"
  type        = string
  default     = "172.16.0.0/28"
}

variable "authorized_network_cidr" {
  description = "CIDR block allowed to access the GKE master endpoint"
  type        = string
  default     = "0.0.0.0/0"
}

# Cluster behavior

variable "deletion_protection" {
  description = "Enable deletion protection for the GKE cluster. Set to false in non-production environments."
  type        = bool
  default     = true
}

# Node pool autoscaling

variable "node_pool_min_count" {
  description = "Minimum number of nodes in the node pool"
  type        = number
  default     = 1
}

variable "node_pool_max_count" {
  description = "Maximum number of nodes in the node pool"
  type        = number
  default     = 5
}

# Namespaces

variable "namespaces" {
  description = "List of Kubernetes namespaces to create"
  type        = list(string)
  default     = ["default", "prod"]
}
