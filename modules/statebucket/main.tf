resource "google_storage_bucket" "tf_state" {
  name                        = "${var.bucket_name}-${var.environment}"
  location                    = var.region
  project                     = var.project
  force_destroy               = false
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type = "Delete"
    }
  }

  lifecycle_rule {
    condition {
      num_newer_versions = 5
    }
    action {
      type = "Delete"
    }
  }

  labels = {
    environment = var.environment
    managed-by  = "terraform"
  }
}
