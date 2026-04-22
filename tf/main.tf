terraform {
  required_version = ">= 1.10"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.33"
    }
  }

  backend "gcs" {
    bucket = "tf-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project
  region  = var.region
}

provider "kubernetes" {
  cluster_ca_certificate = base64decode(data.google_container_cluster.cluster.endpoint)
  host                   = "https://${data.google_container_cluster.cluster.endpoint}"
  token                 = data.google_container_cluster.cluster.access_tokens[0]
}

data "google_container_cluster" "cluster" {
  name     = var.cluster_name
  location = var.region
}