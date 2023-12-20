resource "google_logging_project_sink" "default-bucket" {
  name        = "_Default"
  destination = "logging.googleapis.com/projects/rsd-dev/locations/global/buckets/_Default"
  filter = "NOT LOG_ID(\"cloudaudit.googleapis.com/activity\") AND NOT LOG_ID(\"externalaudit.googleapis.com/activity\") AND NOT LOG_ID(\"cloudaudit.googleapis.com/system_event\") AND NOT LOG_ID(\"externalaudit.googleapis.com/system_event\") AND NOT LOG_ID(\"cloudaudit.googleapis.com/access_transparency\") AND NOT LOG_ID(\"externalaudit.googleapis.com/access_transparency\")"
}
resource "google_logging_project_bucket_config" "default-sink-config" {
    project    = var.project
    location  = "global"
    retention_days = 7
    bucket_id = "_Default"
}
