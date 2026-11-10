# Environments

Three environments are provided: `dev`, `staging`, and `prod`. Each has a
`terraform.tfvars` file under `environments/`. The environment is selected with
`make <target> ENV=<env>`.

## Configuration Matrix

| Setting | dev | staging | prod |
|---------|-----|---------|------|
| Cluster name | `gke-dev` | `gke-staging` | `gke-prod` |
| Environment label | `dev` | `staging` | `production` |
| Deletion protection | off | off | on |
| Main pool | 1-3 x e2-medium | 1-5 x e2-medium | 2-10 x e2-standard-4 |
| Main disk | 50 GB pd-balanced | 100 GB pd-balanced | 200 GB pd-balanced |
| System pool | 1-2 x e2-small | 1-2 x e2-medium | 2-5 x e2-medium |
| Backup | off | 14 days | 30 days |
| Alerting | off | on | on |
| Authorized networks | 0.0.0.0/0 | 0.0.0.0/0 | 10.0.0.0/8 |
| Subnet CIDR | 10.0.0.0/20 | 10.16.0.0/20 | 10.32.0.0/20 |

## Adding an Environment

1. Create `environments/<name>/terraform.tfvars`.
2. Set `environment = "<name>"` (must be one of dev/staging/production/test).
3. Adjust CIDRs so they do not overlap with other environments.
4. Run `make init ENV=<name> && make plan ENV=<name>`.

## CIDR Planning

Each environment uses a disjoint slice of 10.0.0.0/8. Keep pod (10.x.4.0/14)
and service (10.x.8.0/20) ranges inside the environment's allocation and away
from the master CIDR (172.16.0.0/28 family).
