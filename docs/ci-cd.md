# CI/CD

## GitHub Actions

Workflow: `.github/workflows/ci.yml`, triggered on push/PR to `main`.

Jobs:

1. **Terraform Lint & Validate** — fmt check, init, validate (root + modules).
2. **tflint** — Google plugin, root module.
3. **Checkov** — security scan (soft fail).
4. **Terratest Plan-Only** — plan-only test subset.

## Dependabot

`.github/dependabot.yml` configures daily updates for Go modules (`/tests`)
and Terraform providers (`/tf`).

## Pre-commit

`.pre-commit-config.yaml` installs hooks for terraform fmt, validate, tflint,
checkov, and terraform-docs. Install with `pre-commit install`.
