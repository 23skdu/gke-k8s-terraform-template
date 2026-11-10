# Security

## Network

- Private GKE nodes; master endpoint restricted by authorized networks.
- Four firewall rules limit east-west and master access.
- Cloud NAT provides controlled egress.

## Identity

- Dedicated service account per cluster with least-privilege roles.
- Workload Identity maps Kubernetes service accounts to GCP service accounts.

## Node Hardening

- Shielded nodes (secure boot + integrity monitoring).
- Auto-repair and auto-upgrade enabled.

## Supply Chain

- Optional Binary Authorization verifies container images.
- Secret Manager stores sensitive values with IAM bindings.

## State

- Remote state in GCS with versioning and lifecycle rules.
- Uniform bucket-level access and public access prevention enforced.

## Static Analysis

- Checkov runs in CI (`soft_fail: true`); tflint with Google plugin.
- Pre-commit hooks run fmt, validate, tflint, checkov, terraform-docs.
