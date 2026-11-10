# Documentation

This directory contains the extended documentation for the
**GKE Kubernetes Terraform Template**. The root [README](../README.md) is a
quick-start overview; the documents here go deeper into architecture, modules,
configuration, operations, testing, and security.

## Contents

| Document | Description |
|----------|-------------|
| [Architecture](architecture.md) | High-level design, module graph, networking model |
| [Modules](modules.md) | Per-module reference (resources, inputs, outputs) |
| [Variables](variables.md) | Root module variable reference |
| [Environments](environments.md) | dev / staging / prod configuration matrix |
| [Deployment](deployment.md) | Step-by-step deployment and teardown |
| [Operations](operations.md) | Day-2 operations, Makefile targets, troubleshooting |
| [Testing](testing.md) | Terratest suite, plan-only vs integration |
| [Security](security.md) | Security posture, Binary Authorization, IAM, network |
| [CI/CD](ci-cd.md) | GitHub Actions pipelines and pre-commit hooks |

## Quick Links

- Root module entry point: [`tf/`](../tf)
- Reusable modules: [`modules/`](../modules)
- Environment tfvars: [`environments/`](../environments)
- Tests: [`tests/`](../tests)
- Contributing guide: [`CONTRIBUTING.md`](../CONTRIBUTING.md)
- License: [`LICENSE`](../LICENSE)
