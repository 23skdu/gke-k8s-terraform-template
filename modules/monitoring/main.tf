resource "google_logging_project_bucket_config" "default" {
  project        = var.project
  location       = "global"
  bucket_id      = "_Default"
  retention_days = 30
}

# Notification channel for alerts
resource "google_monitoring_notification_channel" "email" {
  count        = var.notification_email != "" ? 1 : 0
  project      = var.project
  display_name = "Alert Email (${var.environment})"
  type         = "email"

  labels = {
    email_address = var.notification_email
  }
}

# Alert: Cluster health
resource "google_monitoring_alert_policy" "cluster_health" {
  count        = var.enable_alerting ? 1 : 0
  project      = var.project
  display_name = "GKE Cluster Health (${var.environment})"
  combiner     = "OR"

  conditions {
    display_name = "Cluster has unhealthy nodes"
    condition_threshold {
      filter     = "resource.type = \"k8s_cluster\" AND metric.type = \"kubernetes.io/node/condition/ready\""
      duration   = "300s"
      comparison = "COMPARISON_LT"
      threshold_value = 1

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_MEAN"
      }

      trigger {
        count = 1
      }
    }
  }

  notification_channels = var.notification_email != "" ? [google_monitoring_notification_channel.email[0].name] : []

  alert_strategy {
    auto_close = "604800s"
  }
}

# Alert: High CPU usage
resource "google_monitoring_alert_policy" "high_cpu" {
  count        = var.enable_alerting ? 1 : 0
  project      = var.project
  display_name = "GKE High CPU Usage (${var.environment})"
  combiner     = "OR"

  conditions {
    display_name = "Node CPU utilization above 80%"
    condition_threshold {
      filter     = "resource.type = \"k8s_node\" AND metric.type = \"kubernetes.io/node/cpu/allocatable_utilization\""
      duration   = "300s"
      comparison = "COMPARISON_GT"
      threshold_value = 0.8

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_MEAN"
      }

      trigger {
        count = 1
      }
    }
  }

  notification_channels = var.notification_email != "" ? [google_monitoring_notification_channel.email[0].name] : []
}

# Alert: High memory usage
resource "google_monitoring_alert_policy" "high_memory" {
  count        = var.enable_alerting ? 1 : 0
  project      = var.project
  display_name = "GKE High Memory Usage (${var.environment})"
  combiner     = "OR"

  conditions {
    display_name = "Node memory utilization above 85%"
    condition_threshold {
      filter     = "resource.type = \"k8s_node\" AND metric.type = \"kubernetes.io/node/memory/allocatable_utilization\""
      duration   = "300s"
      comparison = "COMPARISON_GT"
      threshold_value = 0.85

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_MEAN"
      }

      trigger {
        count = 1
      }
    }
  }

  notification_channels = var.notification_email != "" ? [google_monitoring_notification_channel.email[0].name] : []
}

# Alert: Node auto-repair failure
resource "google_monitoring_alert_policy" "auto_repair_failure" {
  count        = var.enable_alerting ? 1 : 0
  project      = var.project
  display_name = "GKE Node Auto-Repair Failure (${var.environment})"
  combiner     = "OR"

  conditions {
    display_name = "Node auto-repair failed"
    condition_threshold {
      filter     = "resource.type = \"gce_instance\" AND metric.type = \"compute.googleapis.com/instance/autoscaler/ending_scale_down_event\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      threshold_value = 0

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_RATE"
      }

      trigger {
        count = 1
      }
    }
  }

  notification_channels = var.notification_email != "" ? [google_monitoring_notification_channel.email[0].name] : []
}
