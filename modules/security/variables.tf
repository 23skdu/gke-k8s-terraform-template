variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "environment" {
  description = "Environment label"
  type        = string
}

variable "enable_binary_authorization" {
  description = "Enable Binary Authorization policy"
  type        = bool
  default     = false
}

variable "attestor_PROJECTS" {
  description = "List of projects allowed for image attestation"
  type        = list(string)
  default     = []
}
