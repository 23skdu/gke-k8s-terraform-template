# Remote state backend.
#
# Kept in its own file so that tooling (tests, CI dry-runs, plan fixtures) can
# generate a backend-free copy of the root module without editing any file that
# also holds provider or resource configuration.
terraform {
  backend "gcs" {
    bucket = "tf-state"
    prefix = "terraform/state"
  }
}
