resource "google_storage_bucket" "tf_state" {
  name          = "${var.project}-tf-state"
  location      = var.location
  force_destroy = false
  public_access_prevention = "enforced"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      age = 30
    }
    action {
      type = "Delete"
    }
  }
}