project      = "my-gcp-project"
region       = "us-central1"
cluster_name = "gke-prod"
environment  = "production"

# Networking
network_name          = "gke-network"
subnet_cidr           = "10.32.0.0/20"
pods_cidr             = "10.36.0.0/14"
services_cidr         = "10.40.0.0/20"
master_ipv4_cidr_block = "172.16.0.32/28"
authorized_network_cidr = "10.0.0.0/8"

# Cluster
deletion_protection = true

# Main node pool
main_machine_type = "e2-standard-4"
main_min_count    = 2
main_max_count    = 10

# System node pool
system_machine_type = "e2-medium"
system_min_count    = 2
system_max_count    = 5

# Namespaces
namespaces = ["default", "prod", "monitoring", "logging"]

# IAM
namespace = "default"

# Monitoring
notification_email = "ops-team@example.com"
enable_alerting    = true

# Backup
enable_backup        = true
backup_retention_days = 30

# Security
enable_binary_authorization = true

# Secrets
secrets = {
  "db-password" = "Database connection password"
  "api-key"     = "External API key"
}
