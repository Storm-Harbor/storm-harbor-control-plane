# Storm Harbor Technical Documentation

[![DR simulation](https://github.com/Storm-Harbor/storm-harbor-control-plane/actions/workflows/disaster-recovery-simulation.yml/badge.svg)](https://github.com/Storm-Harbor/storm-harbor-control-plane/actions/workflows/disaster-recovery-simulation.yml)
[![Documentation](https://img.shields.io/badge/docs-MkDocs-526CFE)](docs/index.md)
[![Status](https://img.shields.io/badge/status-safe%20simulation-0B7285)](#project-status)

This repository owns the **Storm Harbor technical documentation site**, architecture contracts, configuration schemas, and safe recovery simulations. Storm Harbor is an organization-independent, declarative recovery control plane designed for private projects and multicloud environments spanning AWS, Azure, DigitalOcean, Oracle Cloud and Google Cloud through versioned adapters.

The [platform contract](docs/architecture/platform-contract.md) defines the accepted direction: YAML journeys, Go 1.26.x runtime components, extensible resources/controllers, optional project-owned AI/RAG/CLI integrations, and recoverability with auditable execution. **These are design goals; the current executable code is a GCP-shaped Node.js simulation scaffold.**

The customer-facing website, authentication experience, and operational console live in `storm-harbor-app`. Runtime libraries follow `storm-lib-*`; independently deployable runtime services follow `storm-*`.

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
- provider-specific discovery through versioned cloud adapters;
- read-only Terraform drift detection and remote backend correlation;
- GitOps integration to restore FluxCD and workloads at an immutable revision;
- CRDs for Kubernetes plans, targets, and runs;
- Kubernetes Operator with idempotent control loops and observable status;
- isolated simulations with guardrails, health gates, and mandatory cleanup;
- structured evidence for each phase and a recovery-readiness score;
- adapter architecture ready for AWS, Azure, and other GitOps engines;
- optional project-owned AI, RAG or bounded CLI integrations, with private knowledge, redaction and human review.

The capability list describes the target product. No production cloud adapter, Go runtime or installable CRD is delivered by the current scaffold. [Resource extensibility](docs/architecture/extensibility.md) and the [roadmap](docs/roadmap.md) define the remaining deliveries.

## Recovery flow

```text
preflight → foundation → IAM → Terraform → database → GKE
          → Operator/CRDs → FluxCD → workloads → validation → evidence
```

Each transition depends on a health gate. A failure, timeout, policy violation, or kill switch stops new mutations; cleanup and evidence collection continue.

## Project status

The repository is a **technical documentation and safe-simulation scaffold**, not the runtime control-plane backend. The current pipelines:

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

Follow the [complete control-plane plan and stage checklists](docs/SPECS-CONTROL-PLANE.md) for delivered foundations, prerequisites, acceptance gates and the path from offline validation to private-adopter exercises and retained recovery.

The MkDocs portal contains an executive overview, diagrams, architecture, the federated model, inventory, ITSCM, security, CRDs/Operator, pipelines, and roadmap.

```bash
./scripts/serve-docs.sh
```

The script creates or reuses `.venv-docs`, installs the pinned MkDocs versions, and starts the server. Additional arguments are forwarded to `mkdocs serve`.

Open `http://127.0.0.1:8000`. The sidebar has collapsible sections and a persistent menu button on desktop as well.

To validate the built navigation, run a strict MkDocs build, then `npm ci`,
`npx playwright install chromium` and `npm run test:docs-browser`. The pinned
browser suite reads the real `site/` artifact at its configured Pages origin to
exercise instant navigation, mobile submenus, keyboard access and storage failure.

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

1. Publish the website, English technical documentation and community foundations.
2. Version the provider-neutral YAML/resource contracts and Go 1.26.x SDK/CLI.
3. Deliver tenant-isolated onboarding, inventory and private-source discovery.
4. Introduce cloud and optional AI/RAG/CLI adapters with explicit capability matrices.
5. Implement extensible resources/controllers, native Kubernetes CRDs and a Go Operator.
6. Prove isolated recovery, cleanup and evidence, including supported multicloud mappings.
7. Enable approved real recovery only after the published readiness gates pass.

See [sponsorship](docs/community/sponsorship.md) and [CNCF readiness](docs/community/cncf-readiness.md) for sustainability and community preparation. Storm Harbor has no CNCF affiliation at this stage.

## Contributing

Changes that introduce real operations must document the threat model, minimum permissions, idempotency, rollback, kill switch, expected evidence, and a test plan for a disposable environment.
