variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "environment" {
  description = "Environment label"
  type        = string
}

variable "secrets" {
  description = "Map of secrets to create (key = secret name, value = description)"
  type        = map(string)
  default     = {}
}

variable "service_account_email" {
  description = "Service account to grant secret access"
  type        = string
}
