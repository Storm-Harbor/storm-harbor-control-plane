# Disaster Recovery Control Plane

[![DR simulation](https://github.com/ferreiraad/itscm-mms/actions/workflows/disaster-recovery-simulation.yml/badge.svg)](https://github.com/ferreiraad/itscm-mms/actions/workflows/disaster-recovery-simulation.yml)
[![Documentation](https://img.shields.io/badge/docs-MkDocs-526CFE)](docs/index.md)
[![Status](https://img.shields.io/badge/status-safe%20simulation-0B7285)](#project-status)

**Disaster Recovery Control Plane** discovers dependencies, continuously assesses recoverability, and orchestrates Disaster Recovery across GitHub, Terraform, Google Cloud, and Kubernetes.

The product applies a **Federated Source of Truth (FSoT)**: each system remains authoritative in its domain, while the control plane discovers, normalizes, correlates, and reconciles those sources into a unified view of readiness, RTO, RPO, plans, runs, and evidence.

## The problem it solves

DR plans are often disconnected from the actual infrastructure. Repositories describe intent, APIs expose observed state, backups live in another domain, and test evidence is scattered. The result is a difficult question to answer: **can this service be recovered now, within the agreed RTO and RPO?**

The Disaster Recovery Control Plane connects this information without replacing its original sources:

| Domain | Authoritative source | Control plane responsibility |
| --- | --- | --- |
| DR policy and plan | DR repository | validate, version, and approve |
| Desired infrastructure | Terraform repository | correlate modules, plans, and revisions |
| Terraform state | remote backend | verify availability, lineage, and drift |
| Kubernetes desired state | GitOps repository | resolve the immutable recovery revision |
| Actual cloud state | GCP APIs / Cloud Asset Inventory | discover assets and relationships |
| Actual Kubernetes state | Kubernetes API | reconcile CRDs, workloads, and health |
| Secrets | Secret Manager + KMS | store references only and use ephemeral identity |
| Evidence | external immutable storage | sign, retain, and query |
| Operational index | PostgreSQL | materialize inventory, graph, runs, and metrics |

## Architecture in one sentence

The **control plane outside the recoverable cluster** coordinates discovery, inventory, policies, and recovery; the **Operator inside Kubernetes** reconciles only Kubernetes resources through CRDs.

```text
Federated Sources          DR Control Plane              Recovery targets
GitHub / TF Backend  ───►  Discover + Correlate   ───►  GCP foundation
GCP / Kubernetes API ───►  Policy + ITSCM         ───►  Terraform
Vault / Evidence      ───►  Orchestrate + Prove   ───►  GKE + DR Operator
```

This separation avoids the paradox of hosting the recovery brain inside the cluster it must rebuild.

## Product capabilities

- GitHub App onboarding with explicitly authorized repositories;
- inventory of services, owners, environments, assets, and dependencies;
- mapping of criticality, RTO, RPO, recovery tier, and plan per service;
- GCP-first discovery through APIs and Cloud Asset Inventory;
- read-only Terraform drift detection and remote backend correlation;
- GitOps integration to restore FluxCD and workloads at an immutable revision;
- CRDs for Kubernetes plans, targets, and runs;
- Kubernetes Operator with idempotent control loops and observable status;
- isolated simulations with guardrails, health gates, and mandatory cleanup;
- structured evidence for each phase and a recovery-readiness score;
- adapter architecture ready for AWS, Azure, and other GitOps engines;
- configurable AI providers for Codex Cloud and Gemini, with tiers, reasoning levels, task types, redaction, and human review.

## Recovery flow

```text
preflight → foundation → IAM → Terraform → database → GKE
          → Operator/CRDs → FluxCD → workloads → validation → evidence
```

Each transition depends on a health gate. A failure, timeout, policy violation, or kill switch stops new mutations; cleanup and evidence collection continue.

## Project status

The repository is an **architecture and safe-simulation scaffold**. The current pipelines:

- run only through manual dispatch with the literal confirmation `SIMULATE`;
- retain `contents: read` and do not request a cloud OIDC token;
- validate disposable targets, RPO, and supplied paths;
- run a deterministic phase-based mock engine;
- produce JSON evidence and a Job Summary;
- simulate AI provider/model selection without reading credentials or calling APIs;
- do not authenticate to GCP, run Terraform, or modify Kubernetes.

This makes it possible to evolve contracts, guardrails, and observability before enabling real operations.

## Run locally

The runtime is pinned to **Node.js 26.6.0**.

```bash
nvm install
nvm use
node --run check:runtime
node --run test:disaster-recovery
```

To simulate locally without external access:

```bash
DR_CONFIRMATION=SIMULATE \
DR_DRY_RUN=true \
DR_SCENARIO=terraform-drift \
DR_PHASES=immutable-source,backend-init,refresh-only-plan,classify,policy-evaluation \
DR_TERRAFORM_SCENARIO=resource-deleted \
node scripts/simulate-recovery.mjs
```

## Customer documentation

The MkDocs portal contains an executive overview, diagrams, architecture, the federated model, inventory, ITSCM, security, CRDs/Operator, pipelines, and roadmap.

```bash
./scripts/serve-docs.sh
```

The script creates or reuses `.venv-docs`, installs the pinned MkDocs versions, and starts the server. Additional arguments are forwarded to `mkdocs serve`.

Open `http://127.0.0.1:8000`. The sidebar has collapsible sections and a persistent menu button on desktop as well.

- [Documentation home](docs/index.md)
- [Control plane architecture](docs/architecture/control-plane.md)
- [Editable draw.io XML](docs/assets/diagrams/control-plane.drawio.xml)
- [CRDs and Operator](docs/kubernetes/crds-and-operator.md)
- [Simulation pipelines](docs/operations/pipelines.md)
- [Configuration plan](config/dr.config.example.yaml)

## Essential guardrails

- ephemeral identity through GitHub OIDC and Workload Identity Federation;
- no private keys in snapshots, repositories, or evidence;
- real targets subject to an approval environment and organizational policy;
- immutable source, allowlist, disposable labels, and cost/change limits;
- database restoration always targets a new instance;
- control plane and evidence store outside the recovered failure domain;
- a feature flag never replaces authorization or policy;
- every mutation must be idempotent, auditable, and interruptible.

## Next milestones

1. GitHub App, onboarding, and installation policy.
2. Inventory API and PostgreSQL dependency graph.
3. Discovery adapters for GitHub, Terraform, GCP, and Kubernetes.
4. `DisasterRecoveryPlan`, `RecoveryTarget`, and `RecoveryRun` CRDs.
5. Kubernetes Operator and external orchestrator with durable queues.
6. Simulations in a disposable project and an immutable evidence store.
7. Real recovery with approval, break-glass access, and auditing.

## Contributing

Changes that introduce real operations must document the threat model, minimum permissions, idempotency, rollback, kill switch, expected evidence, and a test plan for a disposable environment.
