variable "project" {
  description = "GCP project ID"
  type        = string

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{4,28}[a-z0-9]$", var.project))
    error_message = "Project ID must be 6-30 characters, lowercase letters, digits, and hyphens only."
  }
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"

  validation {
    condition     = can(regex("^[a-z]+-[a-z]+[0-9]$", var.region))
    error_message = "Region must be a valid GCP region (e.g. us-central1)."
  }
}

variable "environment" {
  description = "Environment label"
  type        = string
  default     = "production"

  validation {
    condition     = contains(["dev", "staging", "production", "test"], var.environment)
    error_message = "Environment must be one of: dev, staging, production, test."
  }
}

# Networking

variable "network_name" {
  description = "VPC network name"
  type        = string
  default     = "gke-network"

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{0,62}$", var.network_name))
    error_message = "Network name must start with a letter, contain only lowercase letters, digits, and hyphens."
  }
}

variable "subnet_cidr" {
  description = "Primary CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/20"

  validation {
    condition     = can(cidrhost(var.subnet_cidr, 0))
    error_message = "Subnet CIDR must be a valid IPv4 CIDR block."
  }
}

variable "pods_cidr" {
  description = "Secondary CIDR range for pods"
  type        = string
  default     = "10.4.0.0/14"

  validation {
    condition     = can(cidrhost(var.pods_cidr, 0))
    error_message = "Pods CIDR must be a valid IPv4 CIDR block."
  }
}

variable "services_cidr" {
  description = "Secondary CIDR range for services"
  type        = string
  default     = "10.8.0.0/20"

  validation {
    condition     = can(cidrhost(var.services_cidr, 0))
    error_message = "Services CIDR must be a valid IPv4 CIDR block."
  }
}

variable "master_ipv4_cidr_block" {
  description = "CIDR block for GKE master endpoint"
  type        = string
  default     = "172.16.0.0/28"

  validation {
    condition     = can(cidrhost(var.master_ipv4_cidr_block, 0))
    error_message = "Master CIDR must be a valid IPv4 CIDR block."
  }
}

variable "authorized_network_cidr" {
  description = "CIDR block for master authorized networks"
  type        = string
  default     = "0.0.0.0/0"

  validation {
    condition     = can(cidrhost(var.authorized_network_cidr, 0))
    error_message = "Authorized network CIDR must be a valid IPv4 CIDR block."
  }
}

# GKE

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
  default     = "gke-cluster"

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{0,38}[a-z0-9]$", var.cluster_name))
    error_message = "Cluster name must be 1-40 characters, start with a letter, contain only lowercase letters, digits, and hyphens."
  }
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
  description = "Disk type for main node pool (pd-standard, pd-ssd, pd-balanced)"
  type        = string
  default     = "pd-balanced"

  validation {
    condition     = contains(["pd-standard", "pd-ssd", "pd-balanced", "hyperdisk-balanced"], var.main_disk_type)
    error_message = "Disk type must be one of: pd-standard, pd-ssd, pd-balanced, hyperdisk-balanced."
  }
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

  validation {
    condition     = length(var.namespaces) > 0
    error_message = "At least one namespace must be specified."
  }
}

# Monitoring

variable "notification_email" {
  description = "Email for alert notifications"
  type        = string
  default     = ""

  validation {
    condition     = var.notification_email == "" || can(regex("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", var.notification_email))
    error_message = "Notification email must be a valid email address or empty."
  }
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

  validation {
    condition     = var.backup_retention_days >= 1
    error_message = "Backup retention must be >= 1 day."
  }
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

# State bucket

variable "state_bucket_name" {
  description = "GCS bucket name for Terraform state"
  type        = string
  default     = "tf-state"

  validation {
    condition     = can(regex("^[a-z][a-z0-9._-]{1,61}[a-z0-9]$", var.state_bucket_name))
    error_message = "Bucket name must be 3-63 characters, lowercase letters, digits, dots, hyphens, underscores."
  }
}
