# Testing

## Overview

Tests are Go-based (Terratest) under `tests/`. Two modes:

- **Plan-only (unit)**: hermetic, no GCP credentials, no resources created.
  Uses a backend-free copy of `tf/` + `modules/` and a static fake token.
- **Integration**: runs against real `tf/` with GCS backend and creates real
  resources. Skipped with `go test -short`.

## Commands

```bash
make test-unit           # plan-only
make test-integration    # integration (costs money)
make test                # all
```

## Oracle Package

`tests/pkg/oracle` provides pure-Go reference implementations of the
validation rules in the Terraform configuration, including CIDR arithmetic
and regex mirrors. `TestHCLValidationMirrorsStayInSync` fails when HCL and
oracle regexes drift apart.

## Quality Gate

`TestOracleQualityGate` is a meta-test ensuring every `Test*` function makes
at least one value-level assertion, plan-only tests never call `Destroy`, and
applying tests call `integrationSkip`.

## CI

Plan-only tests run in GitHub Actions on every push/PR. Integration tests are
not run in CI by default.
