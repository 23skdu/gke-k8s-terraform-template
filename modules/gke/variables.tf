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

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{0,38}[a-z0-9]$", var.cluster_name))
    error_message = "Cluster name must be 1-40 characters, start with a letter, contain only lowercase letters, digits, and hyphens."
  }
}

variable "environment" {
  description = "Environment label"
  type        = string

  validation {
    condition     = contains(["dev", "staging", "production", "test"], var.environment)
    error_message = "Environment must be one of: dev, staging, production, test."
  }
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

  validation {
    condition     = can(cidrhost(var.master_ipv4_cidr_block, 0))
    error_message = "Master CIDR must be a valid IPv4 CIDR block."
  }
}

variable "authorized_network_cidrs" {
  description = "List of CIDR blocks allowed to access the GKE master endpoint"
  type        = list(string)
  default     = ["0.0.0.0/0"]

  validation {
    condition     = length(var.authorized_network_cidrs) > 0
    error_message = "At least one authorized network CIDR must be specified."
  }
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

  validation {
    condition     = var.main_min_count >= 0
    error_message = "Main pool min count must be >= 0."
  }
}

variable "main_max_count" {
  description = "Maximum nodes in main pool"
  type        = number
  default     = 5

  validation {
    condition     = var.main_max_count >= 1
    error_message = "Main pool max count must be >= 1."
  }
}

variable "main_disk_size_gb" {
  description = "Disk size in GB for main node pool"
  type        = number
  default     = 100

  validation {
    condition     = var.main_disk_size_gb >= 10
    error_message = "Main pool disk size must be >= 10 GB."
  }
}

variable "main_disk_type" {
  description = "Disk type for main node pool"
  type        = string
  default     = "pd-balanced"

  validation {
    condition     = contains(["pd-standard", "pd-ssd", "pd-balanced", "hyperdisk-balanced"], var.main_disk_type)
    error_message = "Disk type must be one of: pd-standard, pd-ssd, pd-balanced, hyperdisk-balanced."
  }
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

  validation {
    condition     = var.system_min_count >= 0
    error_message = "System pool min count must be >= 0."
  }
}

variable "system_max_count" {
  description = "Maximum nodes in system pool"
  type        = number
  default     = 3

  validation {
    condition     = var.system_max_count >= 1
    error_message = "System pool max count must be >= 1."
  }
}

variable "system_disk_size_gb" {
  description = "Disk size in GB for system node pool"
  type        = number
  default     = 100

  validation {
    condition     = var.system_disk_size_gb >= 10
    error_message = "System pool disk size must be >= 10 GB."
  }
}

variable "system_disk_type" {
  description = "Disk type for system node pool"
  type        = string
  default     = "pd-balanced"

  validation {
    condition     = contains(["pd-standard", "pd-ssd", "pd-balanced", "hyperdisk-balanced"], var.system_disk_type)
    error_message = "Disk type must be one of: pd-standard, pd-ssd, pd-balanced, hyperdisk-balanced."
  }
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

  validation {
    condition     = length(var.namespaces) > 0
    error_message = "At least one namespace must be specified."
  }
}
