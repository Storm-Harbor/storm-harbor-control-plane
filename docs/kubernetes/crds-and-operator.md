# CRDs and Kubernetes Operator

## Role in the architecture

The Operator is the product's Kubernetes specialist. It runs a control loop inside each managed cluster and reconciles Custom Resources. The external orchestrator creates or monitors those resources, but remains responsible for the global DAG that includes GCP, Terraform, IAM, and databases.

```joint
flowchart LR
  CP[External DR Control Plane]
  API[Kubernetes API]
  OP[DR Operator]
  PLAN[DisasterRecoveryPlan]
  TARGET[RecoveryTarget]
  RUN[RecoveryRun]
  FLUX[FluxCD]
  WORK[Workloads and data services]

  CP --> API
  API --> PLAN
  API --> TARGET
  API --> RUN
  PLAN --> OP
  TARGET --> OP
  RUN --> OP
  OP --> API
  OP --> FLUX
  FLUX --> WORK

  classDef outside fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:3px;
  classDef api fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef crd fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef runtime fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  class CP outside;
  class API api;
  class PLAN,TARGET,RUN crd;
  class OP,FLUX,WORK runtime;
```

!!! danger "Architectural boundary"
    If the cluster disappears, the Operator, Flux, and CRDs disappear with it. The external control plane first recovers the foundation and GKE; only then does it reinstall the Operator and recreate the Custom Resources from the frozen execution plan.

## API group and versioning

The proposed provider-neutral API group is `recovery.controlplane.io`. The first contract may use `v1alpha1`; production requires a conversion strategy and compatibility policy before promotion to `v1beta1`/`v1`.

| CRD | Scope | Source | Purpose |
| --- | --- | --- | --- |
| `DisasterRecoveryPlan` | Namespaced | materialized from the global plan | Kubernetes desired state, objectives, and allowed phases |
| `RecoveryTarget` | Namespaced | discovery + confirmation | target cluster/namespace, constraints, and references |
| `RecoveryRun` | Namespaced | external orchestrator | an immutable run, its approval, and status |

Avoid cluster-scoped CRDs in the MVP. A dedicated management namespace reduces the blast radius and simplifies multi-tenant RBAC.

## DisasterRecoveryPlan

```yaml
apiVersion: recovery.controlplane.io/v1alpha1
kind: DisasterRecoveryPlan
metadata:
  name: checkout-production
  namespace: dr-system
spec:
  serviceRef: checkout
  source:
    repository: platform/gitops
    revision: 4c8d9e7f1a2b3c4d5e6f7890123456789012abcd
    path: clusters/recovery/checkout
  objectives:
    rto: 60m
    rpo: 15m
  strategy:
    fluxBootstrap: true
    orderedHealthGates:
      - infrastructureReady
      - fluxReady
      - workloadsAvailable
      - serviceReachable
  policyRef:
    name: critical-services
status:
  observedGeneration: 3
  conditions: []
```

Rules:

- `revision` is an immutable SHA, never a branch during a run;
- `spec` describes intent; controllers never write to `spec`;
- `status.observedGeneration` shows which generation was processed;
- secret references use `SecretKeySelector`, never an inline value;
- CEL validation prevents dangerous combinations at the API server.

## RecoveryTarget

```yaml
apiVersion: recovery.controlplane.io/v1alpha1
kind: RecoveryTarget
metadata:
  name: checkout-dr-sandbox
  namespace: dr-system
spec:
  provider: gcp
  projectRef: sample-dr-sandbox
  cluster:
    name: checkout-dr-sandbox
    location: europe-west1
    privateEndpointRequired: true
    workloadIdentityRequired: true
  allowedNamespaces:
    - checkout
  disposableLabels:
    purpose: disaster-recovery-simulation
    disposable: "true"
status:
  discoveredClusterUid: 8d8ab676-39b1-49b4-a585-bc9de6b6a45a
  lastVerifiedAt: "2026-08-18T10:00:00Z"
  conditions: []
```

Before each action, the Operator confirms that the UID, project, cluster, and labels still match the target. An identical name is not proof of identity.

## RecoveryRun

```yaml
apiVersion: recovery.controlplane.io/v1alpha1
kind: RecoveryRun
metadata:
  name: checkout-20260818-001
  namespace: dr-system
  labels:
    recovery.controlplane.io/mode: simulation
spec:
  planRef:
    name: checkout-production
  targetRef:
    name: checkout-dr-sandbox
  mode: Simulation
  executionPlanDigest: sha256:2cf24dba5fb0a30e...
  approval:
    id: approval-01
    expiresAt: "2026-08-18T12:00:00Z"
  requestedPhases:
    - BootstrapFlux
    - ReconcileWorkloads
    - ValidateService
status:
  phase: ReconcilingWorkloads
  startedAt: "2026-08-18T10:01:14Z"
  phaseRuns: []
  conditions: []
```

After creation, operational fields in `RecoveryRun.spec` must be immutable through CEL/admission. Cancellation uses an explicit field (`spec.cancel: true`) or a control-plane subresource/API, never arbitrary editing of the executed plan.

## Control loop

```joint
flowchart TD
  W[Watch RecoveryRun event]
  G[Read latest object and generation]
  L[Acquire per-run lease]
  P[Validate plan, target, policy and kill switch]
  D[Compute next idempotent action]
  A[Apply through Kubernetes API]
  H[Evaluate health gate]
  S[Patch status and emit event]
  R[Requeue with bounded backoff]
  F[Finalize evidence references]

  W --> G --> L --> P --> D --> A --> H --> S
  S --> R --> G
  H --> F

  classDef observe fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef decide fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef mutate fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  classDef record fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  class W,G,L observe;
  class P,D,H decide;
  class A mutate;
  class S,R,F record;
```

### Reconciler requirements

- process the latest resource version;
- use a lease/idempotency key to prevent duplicate execution;
- do not keep essential state only in memory;
- use server-side apply with a dedicated field manager;
- limit concurrency by tenant, cluster, and failure domain;
- retry only transient errors, with backoff and jitter;
- classify permanent errors as `Degraded=True`;
- observe `generation`, not loops caused only by `status`;
- respect the global deadline and per-health-gate timeout;
- check the kill switch before every mutating step.

## Conditions and phases

Use standardized, independent conditions:

| Type | Meaning |
| --- | --- |
| `Accepted` | schema, references, and policy are valid |
| `TargetVerified` | target identity and guardrails were confirmed |
| `Progressing` | an action is in progress |
| `Ready` | all requested health gates passed |
| `Degraded` | permanent error or unmet objective |
| `EvidencePublished` | references and digests were persisted externally |

Each condition includes `status`, `reason`, `message`, `observedGeneration`, and `lastTransitionTime`. Messages must not contain tokens, secret payloads, or kubeconfig data.

## Finalizers

The `recovery.controlplane.io/run-protection` finalizer exists only when asynchronous work must finish before removal:

1. block new steps;
2. revoke temporary credentials/leases;
3. execute cleanup allowed by the execution plan;
4. publish the latest status/evidence reference;
5. remove the finalizer.

A finalizer must not remain stuck indefinitely. There must be a deadline, visible reason, and audited administrative procedure. External resources are never deleted merely because a Custom Resource was removed, except under an explicit policy for a disposable target.

## Minimum RBAC

The Operator does not need `cluster-admin`. Separate ServiceAccounts by controller/capability when justified by the blast radius. The core normally requires:

- get/list/watch/patch on its CRDs and `/status`;
- get/create/update on Leases;
- get/list/watch on required workloads and Flux objects;
- patch restricted to allowlisted namespaces;
- create on Events;
- no generic Secret reads; only explicitly referenced names.

## Installation after cluster loss

```joint
sequenceDiagram
  participant C as External control plane
  participant G as GCP adapter
  participant K as New GKE API
  participant O as DR Operator
  participant F as FluxCD

  C->>G: recover network, cluster and identity
  G-->>C: cluster endpoint and verified UID
  C->>K: install pinned Operator bundle
  K-->>C: deployment available
  C->>K: apply Plan, Target and immutable Run
  K->>O: watch RecoveryRun
  O->>F: bootstrap immutable Git revision
  F-->>O: sources and kustomizations ready
  O-->>K: Ready=True and phase metrics
  C->>K: read status and evidence references
```

## Required tests

- envtest/unit tests for defaulting, validation, and transitions;
- idempotency under repeated reconciliation;
- leader loss in the middle of a phase;
- API throttling and watch reconnection;
- target with the same name and a different UID;
- missing or unsigned Git revision;
- timeout and cancellation during Flux reconciliation;
- finalizer with an unavailable external dependency;
- prevention of secret leakage in status, logs, and Events;
- CRD upgrade with objects from earlier versions.
