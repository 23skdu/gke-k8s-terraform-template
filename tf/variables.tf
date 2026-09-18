variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment label"
  type        = string
  default     = "production"
}

# Networking

variable "network_name" {
  description = "VPC network name"
  type        = string
  default     = "gke-network"
}

variable "subnet_cidr" {
  description = "Primary CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/20"
}

variable "pods_cidr" {
  description = "Secondary CIDR range for pods"
  type        = string
  default     = "10.4.0.0/14"
}

variable "services_cidr" {
  description = "Secondary CIDR range for services"
  type        = string
  default     = "10.8.0.0/20"
}

variable "master_ipv4_cidr_block" {
  description = "CIDR block for GKE master endpoint"
  type        = string
  default     = "172.16.0.0/28"
}

variable "authorized_network_cidr" {
  description = "CIDR block for master authorized networks"
  type        = string
  default     = "0.0.0.0/0"
}

# GKE

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
  default     = "gke-cluster"
}

variable "deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = true
}

# Main node pool

variable "main_machine_type" {
  description = "Machine type for main node pool"
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
  description = "Machine type for system node pool"
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

# IAM

variable "namespace" {
  description = "Kubernetes namespace for Workload Identity binding"
  type        = string
  default     = "default"
}

# Namespaces

variable "namespaces" {
  description = "List of Kubernetes namespaces to create"
  type        = list(string)
  default     = ["default", "prod"]
}

# Monitoring

variable "notification_email" {
  description = "Email for alert notifications"
  type        = string
  default     = ""
}

variable "enable_alerting" {
  description = "Enable monitoring alerts"
  type        = bool
  default     = true
}

# Backup

variable "enable_backup" {
  description = "Enable GKE backup plans"
  type        = bool
  default     = true
}

variable "backup_retention_days" {
  description = "Backup retention in days"
  type        = number
  default     = 30
}

# Security

variable "enable_binary_authorization" {
  description = "Enable Binary Authorization"
  type        = bool
  default     = false
}

# Secrets

variable "secrets" {
  description = "Map of Secret Manager secrets to create"
  type        = map(string)
  default     = {}
}
