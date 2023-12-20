resource "google_storage_bucket" "tf-state" {
  name          = "tf-state"
  location      = var.location
  force_destroy = true
  public_access_prevention = "enforced"
  lifecycle_rule {
    condition { age = 3 }
    action { type = "Delete" }
  }
}
