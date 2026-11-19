.PHONY: help init plan apply destroy fmt validate lint test test-unit test-integration clean scan scan-checkov scan-trivy scan-gitleaks scan-secrets scan-fs

ENV ?= dev
TF_DIR = tf

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

init: ## Initialize Terraform
	cd $(TF_DIR) && terraform init

plan: ## Plan Terraform changes (ENV=dev|staging|prod)
	cd $(TF_DIR) && terraform plan -var-file=../environments/$(ENV)/terraform.tfvars -out=tfplan

apply: ## Apply Terraform changes (ENV=dev|staging|prod)
	cd $(TF_DIR) && terraform apply tfplan

destroy: ## Destroy Terraform resources (ENV=dev|staging|prod)
	cd $(TF_DIR) && terraform destroy -var-file=../environments/$(ENV)/terraform.tfvars

fmt: ## Format Terraform files
	cd $(TF_DIR) && terraform fmt -recursive
	cd modules && terraform fmt -recursive

validate: ## Validate Terraform configuration
	cd $(TF_DIR) && terraform init -backend=false && terraform validate

lint: ## Run tflint
	cd $(TF_DIR) && tflint --config=../.tflint.hcl

checkov: ## Run Checkov security scan
	checkov -d $(TF_DIR) --framework terraform --quiet --skip-check CKV_GKE_12,CKV_GKE_13,CKV_GKE_14,CKV_GKE_15,CKV_GKE_16

scan-checkov: ## Run Checkov IaC security scan (all checks)
	checkov -d $(TF_DIR) --framework terraform --quiet --skip-check CKV_GKE_12,CKV_GKE_13,CKV_GKE_14,CKV_GKE_15,CKV_GKE_16

scan-trivy: ## Run Trivy IaC misconfiguration scan
	trivy config $(TF_DIR) --severity HIGH,CRITICAL --exit-code 1 --skip-version-check
	trivy config modules --severity HIGH,CRITICAL --exit-code 1 --skip-version-check
scan-gitleaks: ## Run Gitleaks secret scan
	gitleaks detect --source . --no-git --redact

scan-secrets: ## Run secret scanning (Gitleaks)
	gitleaks detect --source . --no-git --redact

scan-fs: ## Run Trivy filesystem vulnerability and secret scan
	trivy fs . --scanners vuln,secret --severity HIGH,CRITICAL --exit-code 0 --skip-version-check

scan: scan-checkov scan-trivy scan-gitleaks ## Run all security scans
test-unit: ## Run plan-only unit tests (no GCP resources)
	cd tests && go test -v -timeout 10m -run "TestGKEClusterPrivateNodes|TestGKEClusterWorkloadIdentity|TestGKEClusterAdvancedDatapath|TestGKEClusterGatewayAPI|TestGKEClusterReleaseChannel|TestGKEClusterLogging|TestGKEClusterMonitoring|TestGKENodePoolShieldedConfig|TestGKENodePoolAutoUpgrade|TestGKEClusterDeletionProtection|TestSystemNodePoolExists|TestNetworkingVPCExists|TestIAMServiceAccountExists|TestMonitoringAlertPolicies" ./...

test-integration: ## Run full integration tests (creates real GCP resources)
	cd tests && go test -v -timeout 30m ./...

test: ## Run all tests
	cd tests && go test -v -timeout 30m ./...

clean: ## Clean generated files
	cd $(TF_DIR) && rm -rf .terraform .terraform.lock.hcl tfplan
	cd tests && rm -rf .terraform
