terraform {
#  backend "gcs" { }
  required_version = ">= 1.0"
  required_providers {
    google = {  version = "5.8.0" }
    kubernetes = { version = "2.24.0" }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
   }
}
variable "project" { default = "" }
variable "region" { default = "" }
variable "location" { default = "" }
provider "google" {
  region  = var.region
  project = var.project
}
provider "kubernetes" {
  config_path    = "~/.kube/config"
}
provider "kubectl" {
  config_path    = "~/.kube/config"
}
