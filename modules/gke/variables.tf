variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
}

variable "environment" {
  description = "Environment label"
  type        = string
}

variable "network_id" {
  description = "VPC network self link"
  type        = string
}

variable "subnet_id" {
  description = "Subnet self link"
  type        = string
}

variable "pods_range_name" {
  description = "Pods secondary range name"
  type        = string
}

variable "services_range_name" {
  description = "Services secondary range name"
  type        = string
}

variable "master_ipv4_cidr_block" {
  description = "CIDR block for GKE master endpoint"
  type        = string
}

variable "authorized_network_cidr" {
  description = "CIDR block for master authorized networks"
  type        = string
  default     = "0.0.0.0/0"
}

variable "deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = true
}

# Main node pool

variable "main_machine_type" {
  description = "Machine type for the main node pool"
  type        = string
  default     = "e2-medium"
}

variable "main_min_count" {
  description = "Minimum nodes in main pool"
  type        = number
  default     = 1
}

variable "main_max_count" {
  description = "Maximum nodes in main pool"
  type        = number
  default     = 5
}

# System node pool

variable "system_machine_type" {
  description = "Machine type for the system node pool"
  type        = string
  default     = "e2-medium"
}

variable "system_min_count" {
  description = "Minimum nodes in system pool"
  type        = number
  default     = 1
}

variable "system_max_count" {
  description = "Maximum nodes in system pool"
  type        = number
  default     = 3
}

# Service account

variable "service_account_email" {
  description = "Service account email for node pools"
  type        = string
}

# Namespaces

variable "namespaces" {
  description = "List of Kubernetes namespaces to create"
  type        = list(string)
  default     = ["default", "prod"]
}
