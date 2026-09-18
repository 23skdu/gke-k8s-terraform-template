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
modules/         # Reusable Terraform modules
  networking/    # VPC, subnets, firewall, Cloud NAT
  iam/           # Service accounts, IAM roles
  gke/           # GKE cluster, node pools
  monitoring/    # Logging, alerting
  backup/        # GKE backup plans
  security/      # Binary Authorization
  secrets/       # Secret Manager
  statebucket/   # GCS state bucket
  kubernetes/    # Namespaces
environments/    # Per-environment tfvars
tf/              # Root module (calls modules/)
tests/           # Terratest tests
```

## Adding a New Module

1. Create `modules/<name>/` with `main.tf`, `variables.tf`, `outputs.tf`
2. Add validation blocks to all variables
3. Add module call in `tf/modules.tf`
4. Add corresponding variables in `tf/variables.tf`
5. Add outputs in `tf/outputs.tf`
6. Add tests in `tests/<name>_test.go`
7. Update environment tfvars in `environments/*/terraform.tfvars`

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
