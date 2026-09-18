variable "namespaces" {
  description = "List of Kubernetes namespaces to create"
  type        = list(string)
  default     = ["default"]
}

variable "environment" {
  description = "Environment label"
  type        = string
}
