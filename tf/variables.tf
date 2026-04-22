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

variable "location" {
  description = "GCP location (region or zone)"
  type        = string
  default     = "us-central1"
}

variable "node_count" {
  description = "Number of nodes in the node pool"
  type        = number
  default     = 1
}