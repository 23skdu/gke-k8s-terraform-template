# Contributing

## Development Setup

1. Install prerequisites:
   - [Terraform](https://www.terraform.io/downloads) >= 1.15
   - [Go](https://go.dev/doc/install) >= 1.21
   - [tflint](https://github.com/terraform-linters/tflint)
   - [pre-commit](https://pre-commit.com/)
   - [gcloud CLI](https://cloud.google.com/sdk/docs/install)

2. Set up pre-commit hooks:
   ```bash
   pip install pre-commit
   pre-commit install
   ```

## Project Structure

```
modules/         # Reusable Terraform modules (9 total)
  networking/    # VPC, subnets, 4 firewall rules, Cloud NAT
  iam/           # Service accounts, IAM roles, Workload Identity
  gke/           # Cluster, main + system node pools, namespaces
  monitoring/    # Logging bucket, alerting policies
  backup/        # GKE backup plans
  security/      # Binary Authorization
  secrets/       # Secret Manager with IAM bindings
  statebucket/   # GCS state bucket with versioning
  kubernetes/    # Kubernetes namespaces
environments/    # Per-environment tfvars (dev, staging, prod)
tf/              # Root module (main.tf, variables.tf, modules.tf, outputs.tf)
tests/           # Terratest tests (Go)
```

## Adding a New Module

1. Create `modules/<name>/` with:
   - `main.tf` - Resource definitions
   - `variables.tf` - Input variables with `description`, `type`, and `validation` blocks
   - `outputs.tf` - Output values

2. Add module call in `tf/modules.tf`:
   ```hcl
   module "<name>" {
     source = "../modules/<name>"
     # ... variables
   }
   ```

3. Add corresponding variables in `tf/variables.tf` with validation blocks

4. Add outputs in `tf/outputs.tf`

5. Add tests in `tests/<name>_test.go`:
   - Use `plan.ResourcePlannedValuesMap["module.<name>.<resource>"]`
   - Use `GetTerraformOptions(t)` for consistent test config

6. Update environment tfvars in `environments/*/terraform.tfvars`

7. Run `make fmt validate lint test-unit` to verify

## Code Standards

- All variables must have `description`, `type`, and `validation` blocks
- Use `for_each` over `count` for dynamic resources
- Mark sensitive outputs with `sensitive = true`
- Follow [Terraform style guide](https://developer.hashicorp.com/terraform/language/style)

## Testing

```bash
# Plan-only tests (no GCP resources)
make test-unit

# Full integration tests (creates real resources)
make test-integration
```

## Commit Messages

Use conventional commits:
- `feat:` new feature
- `fix:` bug fix
- `docs:` documentation
- `refactor:` code restructuring
- `test:` adding tests
- `chore:` maintenance

## Pull Request Process

1. Create a feature branch from `main`
2. Make changes following the code standards
3. Run `make fmt validate lint test-unit`
4. Update README if adding/changing features
5. Submit PR with description of changes
