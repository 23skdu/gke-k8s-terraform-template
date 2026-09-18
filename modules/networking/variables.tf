variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "environment" {
  description = "Environment label"
  type        = string
}

variable "network_name" {
  description = "VPC network name"
  type        = string
  default     = "gke-network"
}

variable "subnet_name" {
  description = "Subnet name"
  type        = string
  default     = "gke-subnet"
}

variable "subnet_cidr" {
  description = "Primary CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/20"
}

variable "pods_cidr_name" {
  description = "Name of the secondary IP range for pods"
  type        = string
  default     = "pods"
}

variable "pods_cidr" {
  description = "Secondary CIDR range for pods"
  type        = string
  default     = "10.4.0.0/14"
}

variable "services_cidr_name" {
  description = "Name of the secondary IP range for services"
  type        = string
  default     = "services"
}

variable "services_cidr" {
  description = "Secondary CIDR range for services"
  type        = string
  default     = "10.8.0.0/20"
}

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

variable "enable_cloud_nat" {
  description = "Enable Cloud NAT for egress traffic"
  type        = bool
  default     = true
}
