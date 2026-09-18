.PHONY: help init plan apply destroy fmt validate lint test test-unit test-integration clean

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
	checkov -d $(TF_DIR) --framework terraform --quiet

test-unit: ## Run plan-only unit tests (no GCP resources)
	cd tests && go test -v -timeout 10m -run "TestGKEClusterPrivateNodes|TestGKEClusterWorkloadIdentity|TestGKEClusterAdvancedDatapath|TestGKEClusterGatewayAPI|TestGKEClusterReleaseChannel|TestGKEClusterLogging|TestGKEClusterMonitoring|TestGKENodePoolShieldedConfig|TestGKENodePoolAutoUpgrade|TestGKEClusterDeletionProtection|TestSystemNodePoolExists|TestNetworkingVPCExists|TestIAMServiceAccountExists|TestMonitoringAlertPolicies" ./...

test-integration: ## Run full integration tests (creates real GCP resources)
	cd tests && go test -v -timeout 30m ./...

test: ## Run all tests
	cd tests && go test -v -timeout 30m ./...

clean: ## Clean generated files
	cd $(TF_DIR) && rm -rf .terraform .terraform.lock.hcl tfplan
	cd tests && rm -rf .terraform
