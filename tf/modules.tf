# State bucket
module "statebucket" {
  source = "../modules/statebucket"

  project     = var.project
  region      = var.region
  bucket_name = var.state_bucket_name
  environment = var.environment
}

# Networking
module "networking" {
  source = "../modules/networking"

  project                = var.project
  region                 = var.region
  environment            = var.environment
  network_name           = var.network_name
  subnet_cidr            = var.subnet_cidr
  pods_cidr              = var.pods_cidr
  services_cidr          = var.services_cidr
  master_ipv4_cidr_block = var.master_ipv4_cidr_block
  authorized_network_cidr = var.authorized_network_cidr
}

# IAM
module "iam" {
  source = "../modules/iam"

  project      = var.project
  cluster_name = var.cluster_name
  environment  = var.environment
  namespace    = var.namespace
}

# GKE
module "gke" {
  source = "../modules/gke"

  project                 = var.project
  region                  = var.region
  cluster_name            = var.cluster_name
  environment             = var.environment
  network_id              = module.networking.network_id
  subnet_id               = module.networking.subnet_id
  pods_range_name         = module.networking.pods_range_name
  services_range_name     = module.networking.services_range_name
  master_ipv4_cidr_block  = module.networking.master_ipv4_cidr_block
  authorized_network_cidrs = [var.authorized_network_cidr]
  deletion_protection     = var.deletion_protection
  service_account_email   = module.iam.service_account_email
  main_machine_type       = var.main_machine_type
  main_min_count          = var.main_min_count
  main_max_count          = var.main_max_count
  main_disk_size_gb       = var.main_disk_size_gb
  main_disk_type          = var.main_disk_type
  system_machine_type     = var.system_machine_type
  system_min_count        = var.system_min_count
  system_max_count        = var.system_max_count
  system_disk_size_gb     = var.system_disk_size_gb
  system_disk_type        = var.system_disk_type
  namespaces              = var.namespaces
}

# Monitoring
module "monitoring" {
  source = "../modules/monitoring"

  project            = var.project
  region             = var.region
  cluster_name       = var.cluster_name
  environment        = var.environment
  notification_email = var.notification_email
  enable_alerting    = var.enable_alerting
}

# Backup
module "backup" {
  source = "../modules/backup"

  project               = var.project
  region                = var.region
  cluster_name          = var.cluster_name
  environment           = var.environment
  cluster_id            = module.gke.cluster_id
  enable_backup         = var.enable_backup
  backup_retention_days = var.backup_retention_days
}

# Security
module "security" {
  source = "../modules/security"

  project     = var.project
  environment = var.environment
  enable_binary_authorization = var.enable_binary_authorization
}

# Secrets
module "secrets" {
  source = "../modules/secrets"

  project               = var.project
  environment           = var.environment
  secrets               = var.secrets
  service_account_email = module.iam.service_account_email
}

# Kubernetes
module "kubernetes" {
  source = "../modules/kubernetes"

  namespaces  = var.namespaces
  environment = var.environment
}
