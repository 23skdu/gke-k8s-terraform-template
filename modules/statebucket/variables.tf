variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "bucket_name" {
  description = "GCS bucket name for Terraform state"
  type        = string
}

variable "environment" {
  description = "Environment label"
  type        = string
}
