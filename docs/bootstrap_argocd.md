# Bootstrap Argo CD into the GKE cluster

This document describes a reproducible, security-conscious way to bootstrap
[Argo CD](https://argo-cd.readthedocs.io/) into a cluster provisioned by this
Terraform template, and to hand cluster management over to GitOps.

> **Scope.** Terraform in this repository provisions the *cluster* (GKE,
> networking, IAM, monitoring, backup, security, secrets). Argo CD and the
> workloads it manages are intentionally kept out of Terraform so that the
> cluster's day-2 state is owned by Git rather than by `terraform apply`.
> Terraform is used only to create the prerequisites Argo CD needs (a namespace,
> a Workload Identity binding, and optionally a Secret Manager secret).

## 1. Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Git repository (this repo or a dedicated gitops repo)       │
│  ├── bootstrap/          # root Argo CD Application          │
│  ├── platform/           # cluster-wide components          │
│  └── apps/               # workload Applications            │
└──────────────┬──────────────────────────────────────────────┘
               │  (SSH deploy key / GitHub App)
               ▼
┌─────────────────────────────────────────────────────────────┐
│  GKE cluster                                                │
│  ├── argocd namespace  (Argo CD control plane)              │
│  ├── root Application  (App of Apps)                        │
│  └── child Applications (synced from Git)                   │
└─────────────────────────────────────────────────────────────┘
```

Bootstrap flow:

1. Terraform creates the cluster and the `argocd` namespace.
2. An operator installs Argo CD with Helm (version pinned).
3. Argo CD is pointed at the Git repository.
4. A single root `Application` (App of Apps) is applied; Argo CD then takes
   over and reconciles everything under `platform/` and `apps/`.

## 2. Prerequisites

- Cluster created: `make apply ENV=<env>` and credentials fetched:
  ```bash
  gcloud container clusters get-credentials <cluster> --region=<region> --project=<project>
  ```
- Tooling installed locally:
  - `kubectl` (>= 1.28)
  - `helm` (>= 3.12)
  - `argocd` CLI (optional but recommended)
- A Git repository reachable from the cluster (GitHub/GitLab/etc.) that will
  hold the Argo CD manifests.
- (Recommended) A dedicated GCP service account for Argo CD if it must touch
  GCP APIs (see § 7 Workload Identity).

## 3. Step 1 — Create the Argo CD namespace

Create the namespace with Terraform (preferred) or kubectl. The template's
`modules/kubernetes` can be extended, or apply directly:

```bash
kubectl create namespace argocd
kubectl label namespace argocd \
  environment=<env> \
  managed-by=argocd
```

Enable Pod Security Standards on the namespace:

```bash
kubectl label namespace argocd \
  pod-security.kubernetes.io/enforce=baseline \
  pod-security.kubernetes.io/enforce-version=latest \
  pod-security.kubernetes.io/warn=restricted \
  pod-security.kubernetes.io/warn-version=latest
```

## 4. Step 2 — Install Argo CD with Helm

Pin the chart version and use a values file committed to Git.

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

helm upgrade --install argocd argo/argo-cd \
  --namespace argocd \
  --version 7.7.0 \
  --values gitops/argocd/values-<env>.yaml \
  --wait --timeout 10m
```

Recommended `values-<env>.yaml` (commit under `gitops/argocd/`):

```yaml
global:
  domain: argocd.<env>.example.com

configs:
  params:
    # Disable anonymous access; require SSO/OIDC in prod.
    server.insecure: "false"
  cm:
    # Keep the UI read-only for unauthenticated users.
    users.anonymous.enabled: "false"

server:
  replicas: 2
  autoscaling:
    enabled: true
    minReplicas: 2
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 256Mi

controller:
  replicas: 1
  resources:
    limits:
      cpu: "1"
      memory: 1Gi

repoServer:
  replicas: 2
  resources:
    limits:
      cpu: 500m
      memory: 512Mi

redis:
  enabled: true

dex:
  enabled: false   # enable only if using Dex for SSO; prefer external OIDC

# High availability (prod only)
ha:
  enabled: true
```

> **HashiCorp/Terraform alternative.** If you prefer Terraform to own the
> installation, add the `helm` provider to `tf/main.tf` and a
> `helm_release` resource. This is discouraged for the Argo CD control plane
> itself because it creates a circular dependency (Terraform needs the cluster,
> Argo CD needs Terraform outputs). Prefer the manual/Helm step above, then
> let Argo CD self-manage via the App of Apps pattern (§ 6).

## 5. Step 3 — Access the Argo CD API/UI

Port-forward (initial bootstrap only):

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
# open https://localhost:8080
```

Retrieve the initial admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d; echo
```

Log in with the CLI:

```bash
argocd login localhost:8080 --username admin --password <password> --insecure
```

**Immediately rotate or disable the local admin account** once SSO is
configured (§ 8).

## 6. Step 4 — Connect the Git repository

Create a read-only deploy key in the Git provider and store the private key
as a Kubernetes secret in the `argocd` namespace.

Using the Argo CD CLI:

```bash
argocd repo add git@github.com:<org>/<gitops-repo>.git \
  --ssh-private-key-path ~/.ssh/argocd_deploy_key \
  --name gitops
```

Or declaratively (committed as a SealedSecret / ExternalSecret):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: repo-gitops
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  type: git
  url: git@github.com:<org>/<gitops-repo>.git
  sshPrivateKey: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    ...
    -----END OPENSSH PRIVATE KEY-----
```

> **Security.** Prefer a GitHub App or a read-only deploy key scoped to the
> GitOps repository only. Never reuse a personal access token. Rotate the key
> on a schedule and on any suspected compromise.

## 7. Step 5 — Apply the root Application (App of Apps)

Create `bootstrap/root-app.yaml` in the GitOps repository:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: root
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: git@github.com:<org>/<gitops-repo>.git
    targetRevision: main
    path: bootstrap
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

Apply it once:

```bash
kubectl apply -f bootstrap/root-app.yaml
```

From this point on, Argo CD reconciles the cluster. Add child `Application`
manifests under `bootstrap/` for each platform component; the root app will
pick them up automatically.

## 8. Step 6 — Bootstrap platform components

Recommended components (each as an Argo CD `Application`):

| Component | Purpose |
|-----------|---------|
| `namespaces` | Create workload namespaces with labels/PSS |
| `ingress-nginx` or Gateway API | Ingress |
| `cert-manager` | TLS certificates |
| `external-secrets` | Sync Secret Manager secrets into Kubernetes |
| `prometheus`/`grafana` | Observability (if not using Managed Prometheus) |
| `kyverno` or OPA Gatekeeper | Policy enforcement |
| `sealed-secrets` | Encrypted secrets in Git |

Example child application:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: external-secrets
  namespace: argocd
spec:
  project: default
  source:
    repoURL: git@github.com:<org>/<gitops-repo>.git
    targetRevision: main
    path: platform/external-secrets
    helm:
      valueFiles:
        - values-<env>.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: external-secrets
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

## 9. Step 7 — Workload Identity for Argo CD (optional)

If Argo CD (or External Secrets) must call GCP APIs, bind a Kubernetes service
account to a dedicated GCP service account using the template's Workload
Identity pattern:

```bash
# GCP service account
gcloud iam service-accounts create argocd-<env> \
  --project=<project> --display-name="Argo CD <env>"

# Workload Identity binding (mirrors modules/iam)
gcloud iam service-accounts add-iam-policy-binding \
  argocd-<env>@<project>.iam.gserviceaccount.com \
  --role=roles/iam.workloadIdentityUser \
  --member="serviceAccount:<project>.svc.id.goog[argocd/argocd-server]"
```

Annotate the Kubernetes service account:

```bash
kubectl annotate serviceaccount argocd-server -n argocd \
  iam.gke.io/gcp-service-account=argocd-<env>@<project>.iam.gserviceaccount.com
```

Grant only the least-privilege roles the component needs (e.g.
`roles/secretmanager.secretAccessor` for External Secrets).

## 10. Step 8 — SSO, RBAC, and hardening

1. **Enable SSO** (OIDC via your IdP) in `argocd-cm` / `argocd-rbac-cm`.
2. **Disable the local admin** account:
   ```yaml
   # argocd-cm
   admin.enabled: "false"
   ```
3. **RBAC**: grant `role:readonly` by default; restrict `role:admin` to a
   small group. Use Argo CD projects to isolate teams.
4. **TLS**: terminate TLS at an Ingress/LoadBalancer with a managed
   certificate; set `server.insecure: "false"`.
5. **Network policy**: restrict Argo CD egress to the Git provider and GCP
   APIs; restrict ingress to the Ingress controller.
6. **Audit**: enable Argo CD audit logging and ship to Cloud Logging.
7. **Secrets**: never commit plaintext secrets. Use Sealed Secrets or
   External Secrets + Secret Manager.
8. **Image provenance**: enable Binary Authorization (already optional in
   this template) and require signed images for the `argocd` namespace.

## 11. Step 9 — Self-management (optional)

Once stable, let Argo CD manage its own installation by creating an
`Application` pointing at `gitops/argocd/` with the Helm chart. This removes
the manual Helm step from runbooks. Enable only after the root app is healthy.

## 12. DevSecOps checklist

- [ ] Argo CD chart version pinned and updated via Dependabot.
- [ ] Deploy key / GitHub App scoped read-only to the GitOps repo.
- [ ] No plaintext secrets in Git (gitleaks clean).
- [ ] IaC scans (checkov/trivy) green: `make scan`.
- [ ] Pod Security Standards enforced on all namespaces.
- [ ] Network policies applied to `argocd` namespace.
- [ ] RBAC least privilege; local admin disabled.
- [ ] Audit logs enabled and reviewed.
- [ ] Backup of Argo CD (etcd/Redis) included in GKE Backup plan.

## 13. Troubleshooting

| Symptom | Likely cause / fix |
|---------|--------------------|
| `argocd-server` CrashLoopBackOff | Check resource limits; ensure Redis reachable. |
| Application stuck `OutOfSync` | Run `argocd app sync <name> --prune`; inspect diff. |
| `Permission denied (publickey)` on repo | Deploy key not registered or wrong secret type label. |
| Workload Identity not working | Verify service account annotation and IAM binding; check `GKE_METADATA` mode. |
| Cluster unreachable from Argo CD | Confirm `https://kubernetes.default.svc` and Argo CD RBAC. |

## 14. References

- Argo CD documentation: https://argo-cd.readthedocs.io/
- Argo CD Helm chart: https://github.com/argoproj/argo-helm
- App of Apps pattern: https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/
- External Secrets Operator: https://external-secrets.io/
- GKE Workload Identity: https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity
- This template: `docs/security.md`, `docs/architecture.md`, `Makefile` (`make scan`)
