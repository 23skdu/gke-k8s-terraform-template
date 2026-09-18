project      = "my-gcp-project"
region       = "us-central1"
cluster_name = "gke-staging"
environment  = "staging"

# Networking
network_name          = "gke-network"
subnet_cidr           = "10.16.0.0/20"
pods_cidr             = "10.20.0.0/14"
services_cidr         = "10.24.0.0/20"
master_ipv4_cidr_block = "172.16.0.16/28"
authorized_network_cidr = "0.0.0.0/0"

# Cluster
deletion_protection = false

# Main node pool
main_machine_type  = "e2-medium"
main_min_count     = 1
main_max_count     = 5
main_disk_size_gb  = 100
main_disk_type     = "pd-balanced"

# System node pool
system_machine_type = "e2-medium"
system_min_count    = 1
system_max_count    = 2
system_disk_size_gb = 100
system_disk_type    = "pd-balanced"

# Namespaces
namespaces = ["default", "staging"]

# IAM
namespace = "default"

# Monitoring
notification_email = "team@example.com"
enable_alerting    = true

# Backup
enable_backup         = true
backup_retention_days = 14

# Security
enable_binary_authorization = false

# Secrets
secrets = {}

# State bucket
state_bucket_name = "tf-state"
