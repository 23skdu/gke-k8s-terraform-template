# Architecture

## Overview

The template provisions a production-grade GKE cluster on Google Cloud using a
modular Terraform layout. The root module (`tf/`) composes nine reusable modules
under `modules/`, each owning a single concern. Environments differ only through
their `terraform.tfvars` files, so the same code deploys dev, staging, and prod.

## Module Dependency Graph

```
                    ┌─────────────┐
                    │ statebucket │  (GCS state bucket)
                    └─────────────┘

┌────────────┐     ┌────────────┐     ┌────────────┐
│ networking │────▶│    gke     │◀────│    iam     │
│  (VPC,     │     │ (cluster,  │     │ (SA, IAM,  │
│  subnets,  │     │  node      │     │  workload  │
│  firewall, │     │  pools)    │     │  identity) │
│  NAT)      │     └─────┬──────┘     └────────────┘
└────────────┘           │
                         ├──────────▶┌────────────┐
                         │           │ monitoring │
                         │           │ (logging,  │
                         │           │  alerts)   │
                         │           └────────────┘
                         │
                         ├──────────▶┌────────────┐
                         │           │   backup   │
                         │           │ (GKE       │
                         │           │  backups)  │
                         │           └────────────┘
                         │
                         ├──────────▶┌────────────┐
                         │           │  security  │
                         │           │ (Binary    │
                         │           │  Authz)    │
                         │           └────────────┘
                         │
                         ├──────────▶┌────────────┐
                         │           │  secrets   │
                         │           │ (Secret    │
                         │           │  Manager)  │
                         │           └────────────┘
                         │
                         └──────────▶┌────────────┐
                                     │ kubernetes │
                                     │(namespaces)│
                                     └────────────┘
```

## Networking Model

- A dedicated VPC with `auto_create_subnetworks = false`.
- One subnet in the chosen region carrying:
  - primary range: node IPs (`subnet_cidr`)
  - secondary range `pods`: pod IPs (`pods_cidr`)
  - secondary range `services`: service IPs (`services_cidr`)
- Four firewall rules:
  - `allow_internal` — node/pod/service ranges, all protocols
  - `allow_health_checks` — Google health-check ranges to tagged nodes
  - `allow_master_to_nodes` — master CIDR to nodes on 443/10250
  - `allow_nodes_to_master` — nodes to master CIDR on 443/10250
- Cloud Router + Cloud NAT (optional) for egress from private nodes.
- GKE master uses a private endpoint in `master_ipv4_cidr_block` (a /28).

## Cluster Model

- GKE **Standard** (not Autopilot), regional cluster.
- VPC-native networking, private nodes, Workload Identity enabled.
- `ADVANCED_DATAPATH` for network policy enforcement.
- Gateway API enabled.
- Shielded nodes (secure boot + integrity monitoring).
- Release channel (REGULAR) with auto-upgrade/repair on node pools.
- Two node pools:
  - **main**: autoscaling general workloads
  - **system**: tainted (`node-pool=system:PreferNoSchedule`) for critical
    workloads

## State Management

- Remote state is stored in a GCS bucket created by the `statebucket` module.
- The backend declaration lives in `tf/backend.tf` so tooling can generate a
  backend-free copy of the root module for hermetic plan-only tests.
- The bucket has versioning and lifecycle rules (delete objects older than 90
  days, keep 5 newer versions).

## Test Strategy

- Plan-only tests run against a backend-free copy of `tf/` plus `modules/`
  in a temporary directory with a static fake GCP token, so they need no
  credentials and never touch GCP.
- Integration tests run against the real `tf/` directory with its GCS backend
  and create real resources; they are skipped with `go test -short`.

See [Testing](testing.md) for details.
