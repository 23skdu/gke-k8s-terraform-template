# Operations

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make help` | Show available targets |
| `make init ENV=<env>` | `terraform init` |
| `make plan ENV=<env>` | Plan and save `tfplan` |
| `make apply ENV=<env>` | Apply saved plan |
| `make destroy ENV=<env>` | Destroy resources |
| `make fmt` | Format all `.tf` files |
| `make validate` | Validate root + modules (backend-free) |
| `make lint` | Run tflint |
| `make checkov` | Security scan |
| `make test-unit` | Plan-only tests |
| `make test-integration` | Full integration tests |
| `make test` | All tests |
| `make clean` | Remove generated files |

## Troubleshooting

- **Backend initialization required**: run `terraform init -reconfigure` in `tf/`.
- **Provider auth errors**: ensure `GOOGLE_OAUTH_ACCESS_TOKEN` or ADC is set.
- **Plan-only tests failing on backend**: confirm `tf/backend.tf` is skipped by
  the fixture (it is, via `fixtureSkip`).
- **Node pool taints**: system pool uses `node-pool=system:PreferNoSchedule`;
  schedule critical workloads with matching tolerations.
- **CIDR overlaps**: use the oracle package (`tests/pkg/oracle`) to check
  CIDR overlap before applying.

## Day-2 Changes

1. Edit `environments/<env>/terraform.tfvars` or module code.
2. Run `make fmt && make validate && make lint`.
3. Run `make plan ENV=<env>` and review.
4. Run `make apply ENV=<env>`.
