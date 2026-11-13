# Harbor Lens — Plan and acceptance checklist

Harbor Lens is Storm Harbor's own operational frontend and public reference application. The owner accepted this direction on 2026-11-13, replacing the earlier Harbor Commons proposal. OpenLens and Terraform dashboards inspire the experience; they are not adopted runtimes, forks or authentication systems. Gitea is not the selected application.

This is accepted architecture, not a claim of implemented production capabilities. GitHub owns live status: [Journey #16](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/16) → [Epic #17](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/17) → [specification #18](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/18) and [sandbox #1](https://github.com/Storm-Harbor/harbor-lens/issues/1).

## Ownership

| Repository | Responsibility |
| --- | --- |
| `storm-harbor-control-plane` | Go APIs, CLI, GitHub App/bots, orchestration, collectors, adapters, controllers and contract ownership |
| `harbor-lens` | Separate browser application, presentation and sandbox; future consumer of published core APIs |
| `storm-harbor-app` | Existing multilingual product website at `stormharbor.io`; stays in its repository |
| Customer recovery repository | Private YAML intent and optional historical report site; no fork of the engine |
| Existing contract module | Immutable `RecoveryConfiguration v1alpha1` compatibility dependency |

Backend components target Go 1.26.x. Browser HTML/CSS/JavaScript and development/build tools are separate. No new core repository is needed for a collector or Lens API. The first UI slice proceeds independently because it uses synthetic local fixtures and consumes no unpublished API.

## Harbors and federation

A **Harbor** is a registered control-plane instance or managed environment boundary, bound to authenticated organization/project scope and explicit capabilities. It is not necessarily a Kubernetes cluster, cloud account or public endpoint. Final resource/identifier contracts remain to be designed.

A **Harbor group** is a logical federation of explicitly authorized Harbors. Membership does not merge permissions or grant cross-organization access. Reads, counts, traversal, cache, evidence and actions enforce each member's effective scope before data crosses a boundary.

```joint
flowchart LR
  Lens[Harbor Lens] --> API[Scoped Go API]
  API --> A[Authorized Harbor A]
  API --> B[Authorized Harbor B]
  A --> Cloud[Cloud and Kubernetes adapters]
  B --> IaC[Infrastructure and metrics adapters]
```

This is the planned live flow. The initial sandbox has no API, identity service or collector.

## Three distinct modes

| Mode | Data and effects | Gate |
| --- | --- | --- |
| Sandbox | Synthetic inventory, metrics, diagrams and simulated runs; no external effects | Visible synthetic identity and isolated state |
| Observation | Customer-selected sources, bounded read-only inventory/metrics | Published contracts, authenticated scope, provider read grants and collection/egress/retention policy |
| Operation | Exact approved plans executed through core orchestration | Independent policy, frozen plan, approval and scoped workload identity |

Inspired by [Stripe sandboxes](https://docs.stripe.com/sandboxes), isolation covers identity, storage, caches, jobs, evidence and external destinations. A mode switch cannot grant authority. Synthetic outcomes are not readiness evidence. A real cloud drill is a separately authorized, budgeted operation with target lifecycle and cleanup controls. This analogy introduces no payment dependency.

## Views and private data

The intended views cover authorized Harbor inventory, cloud/Kubernetes/databases, desired/managed/observed state, Terraform ownership/drift, freshness and unknown states, accessible topology, YAML journeys, target environment/region, exercise versus restore, create versus reuse, phase outcomes, cancellation, gaps and history. Actual measured RTO/RPO must remain distinct from targets and simulations.

Optional opt-in AI agents keep the same permitted repositories, data policy and proposal review gates. Agent selection never grants authority or introduces a shared knowledge corpus.

Customers explicitly select repositories, sources and capabilities. GitHub App access does not grant provider IAM or permission to export code to an AI provider. Planned Go collectors may run inside customer networks, including outbound-only/self-hosted execution. Secrets stay in customer-selected vaults; prefer short-lived workload identity. Static credentials require separately managed rotation/lifecycle.

Never put raw Terraform state, kubeconfig, tokens, credentials, raw logs, prompts or knowledge bases into browser fixtures, public Pages, logs or frontend storage. Connections require destination/protocol allowlists, SSRF/DNS/private-network and egress policy, bounded polling, timeout/cancellation, retry, cardinality, bytes, cost and retention. Sensitive/high-cardinality customer identifiers cannot become metric labels. Reauthorize downloads, caches, jobs and evidence after membership/policy changes.

The separate [repository dashboard](getting-started/recovery-dashboard.md), Stage 08A, remains a historical evidence projection with independently enforced hosting audience. A private repository alone does not establish private Pages hosting.

## Contracts before live consumers

Publish reviewed contracts for Harbor registration/capabilities/revocation; organization/project and Harbor group membership; connections/workload identity/data policy; mode-bound identities/stores/destinations; allowlisted inventory/health/freshness/topology/metrics; safe errors/pagination/cancellation; and exact-plan operations/idempotency/history/evidence.

Core owns these semantics. Producer validation, manual merge, immutable publication and acceptance precede live consumers. Current `RecoveryConfiguration recovery.stormharbor.io/v1alpha1` accepts none of `spec.lens`, `spec.harbors`, `spec.mode`, `spec.metrics` or `spec.ai`. This plan is not executable YAML; do not add speculative fields to the alpha validator.

## Delivery checklist

### L0 — Recorded scope

- [x] Accept our own separate Harbor Lens application; preserve website ownership.
- [x] Register native Journey/Epic and owning issues in Project #1.
- [x] Create `harbor-lens` without live integrations.
- [ ] Manually merge/publish this specification and verify its public URL.

### L1 — Synthetic sandbox

Estimate: specification 6 focused hours plus frontend 16 focused hours; not elapsed work or production qualification.

- [ ] Build original responsive overview, inventory, metrics, topology, recovery and history views.
- [ ] Keep a persistent sandbox/synthetic label; live connections, observation and operation stay unavailable.
- [ ] Provide synthetic Harbor creation and simulation without credentials, network APIs or browser persistence.
- [ ] Label client-side scope filtering as demonstration, never backend authorization.
- [ ] Maintain locale parity, keyboard navigation, mobile menu, bounded safe labels and accessible graph alternatives.
- [ ] Exercise the actual built artifact: negative paths, isolation, run history, cancellation/reset, i18n and responsive behavior.
- [ ] Measure statements, branches, functions and lines independently; target 100% each and record gaps honestly.
- [ ] Prepare source PR and local preview; deployment and live integrations are separate deliveries.

### L2 — Authoritative contracts and identity

- [ ] Publish reviewed registration, membership, mode, projection and operation contracts.
- [ ] Enforce scope before storage, retrieval, counts, graphs, exports and cache access.
- [ ] Test colliding IDs, switched/concurrent scope, revocation, jobs, idempotency and group membership.
- [ ] Separate sandbox/observation/operation identity, storage, destinations and evidence.

### L3 — Read-only observation

- [ ] Implement scoped Go connection/collector ports and initial adapters.
- [ ] Add customer-controlled workload identity, endpoint/egress policy and collection budgets.
- [ ] Handle denied, stale and partial sources, timeout/cancellation and bounded retries without inventing healthy metrics.
- [ ] Validate isolated dependencies/safe provider sandboxes; prove no mutation or foreign data exposure.

### L4 — Harbor groups

- [ ] Implement explicit membership, per-member projections and access revalidation.
- [ ] Preserve independent failure domains, quotas and freshness during partial outages.
- [ ] Require a separately owned administrative contract for any cross-organization aggregation.

### L5 — Approved operations

- [ ] Connect core planning/approval/execution only after producer publication gates.
- [ ] Freeze scope, source revision, target, journey/schema/controller versions and exact-plan digest.
- [ ] Test replay, concurrency, cancellation, partial failure, rollback/compensation and evidence retention.
- [ ] Keep observation grants distinct from mutation; browser state and AI output never authorize actions.

### L6 — Public reference and production qualification

- [ ] Define synthetic public workloads and private deployment bindings; begin with GCP Kubernetes/Cloud SQL.
- [ ] Add a second cloud and tested source/target compatibility before claiming multicloud recovery coverage.
- [ ] Cover Kubernetes, databases, IaC, backup, secrets, networking, regions and recovery targets through explicit scenarios.
- [ ] Recover Lens with independent artifacts, identity, backups/evidence and a qualified CLI while Lens is unavailable.
- [ ] Meet production API/frontend coverage dimensions and applicable E2E scenarios; report deployment and recovery evidence separately.

## Acceptance evidence

| Boundary | Required evidence |
| --- | --- |
| Browser sandbox | Actual artifact, no live destinations/persisted credentials, safe input, keyboard/mobile/i18n, synthetic history |
| Scope/federation | Absent foreign rows/counts/graphs/cache/evidence, revocation and concurrent-scope scenarios |
| Collection | Read-only capability, bounded egress/cardinality, cancellation, stale/partial results and retention |
| Operation | Independent exact-plan approval, scoped idempotency, replay/concurrency denial and immutable evidence |
| Resilience | Lens outage cannot prevent emergency CLI, vault access or evidence recovery |

Design references only: [OpenLens](https://github.com/MuhammedKalkan/OpenLens), [Terraboard](https://github.com/camptocamp/terraboard), [Stripe sandboxes](https://docs.stripe.com/sandboxes). No runtime or code from these projects is adopted.
