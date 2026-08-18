# Control plane architecture

## Primary decision

The control plane runs outside the GKE environment it may need to recover. It does not replace GitHub, Terraform, GCP, Kubernetes, or the evidence store: it maintains a materialized model with provenance and freshness, and orchestrates operations through adapters.

[Download the editable draw.io XML diagram](../assets/diagrams/control-plane.drawio.xml){ .md-button .md-button--primary }

## Layered diagram

```joint
flowchart TB
  subgraph L1[EXPERIENCE]
    UI[Customer portal]
    API[Public API and CLI]
    APP[GitHub App]
  end

  subgraph L2[FEDERATED SOURCES OF TRUTH]
    DR[DR repository<br/>policies and plans]
    TG[Terraform repository<br/>desired infrastructure]
    TB[Remote backend<br/>Terraform state]
    GO[GitOps repository<br/>desired workloads]
    GA[GCP APIs and CAI<br/>actual cloud state]
    KA[Kubernetes API<br/>actual cluster state]
    VA[Secret Manager and KMS<br/>secret authority]
  end

  subgraph L3[DR CONTROL PLANE — OUTSIDE RECOVERY TARGET]
    IN[Ingestion and adapters]
    IV[Inventory service]
    DG[Dependency graph]
    PE[Policy and ITSCM engine]
    OR[Recovery orchestrator]
    EV[Evidence service]
    PR[Plugin registry]
    IN --> IV --> DG --> PE --> OR --> EV
    PR --> IN
    PR --> OR
  end

  subgraph L4[EXECUTION]
    TF[Terraform runner]
    CP[GCP foundation adapter]
    DB[Database recovery adapter]
    KO[Kubernetes DR Operator]
    FX[FluxCD bootstrap adapter]
  end

  subgraph L5[PERSISTENCE AND OBSERVABILITY]
    PG[(PostgreSQL)]
    Q[Durable queue]
    OT[OpenTelemetry]
    ES[Immutable evidence store]
  end

  UI --> API
  APP --> IN
  DR --> IN
  TG --> IN
  TB --> IN
  GO --> IN
  GA --> IN
  KA --> IN
  VA --> OR
  OR --> TF
  OR --> CP
  OR --> DB
  OR --> KO
  OR --> FX
  IV --> PG
  DG --> PG
  OR --> Q
  OR --> OT
  EV --> ES

  classDef experience fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef source fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef core fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef execute fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  classDef persist fill:#eeeeee,stroke:#666666,color:#1d1d1b,stroke-width:2px;
  class UI,API,APP experience;
  class DR,TG,TB,GO,GA,KA,VA source;
  class IN,IV,DG,PE,OR,EV,PR core;
  class TF,CP,DB,KO,FX execute;
  class PG,Q,OT,ES persist;
```

All connectors are straight in the portal and orthogonal in the draw.io XML. Colors identify responsibility, not operational state.

## Components and contracts

| Component | Responsibility | Must not |
| --- | --- | --- |
| Ingestion/adapters | fetch deltas, normalize IDs, record provenance | decide recovery alone |
| Inventory service | materialize services, assets, owners, and environments | replace the authoritative source |
| Dependency graph | store typed relationships and generate DAGs | infer a critical dependency without confidence |
| Policy/ITSCM engine | assess gaps, tiers, RTO/RPO, and readiness | grant cloud permission |
| Orchestrator | freeze the plan, execute phases and health gates | reside only in the recoverable cluster |
| Evidence service | record inputs, decisions, metrics, and outputs | store secrets or tokens |
| Plugin registry | resolve adapter capabilities and versions | load an unsigned plugin in production |

## GCP-first topology

The initial deployment should use a **dedicated management/DR project**, separate from protected projects. The database, queue, signing keys, and evidence should span regions or projects according to the threat model. For SaaS, the control plane can move to a vendor-operated account without changing domain contracts.

```joint
flowchart LR
  subgraph M[MANAGEMENT PROJECT]
    C[Control plane]
    P[(PostgreSQL HA)]
    Q[Durable queue]
    K[KMS signing keys]
  end
  subgraph V[SECURITY VAULT PROJECT]
    S[Secret Manager]
    E[Retention-locked evidence]
  end
  subgraph T[PROTECTED PROJECT]
    G[GKE and workloads]
    D[Cloud SQL]
    I[Cloud resources]
  end
  subgraph R[RECOVERY PROJECT]
    RG[Recovery GKE]
    RD[Restored database]
    RI[Recovered foundation]
  end

  C --> P
  C --> Q
  C --> K
  C --> S
  C --> I
  C --> G
  C --> D
  C --> RI
  C --> RD
  C --> RG
  C --> E

  classDef management fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef vault fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef protected fill:#fdebed,stroke:#b0000c,color:#1d1d1b,stroke-width:2px;
  classDef recovery fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  class C,P,Q,K management;
  class S,E vault;
  class G,D,I protected;
  class RG,RD,RI recovery;
```

## Orchestration as a state machine

A run must be persisted before it starts and may advance only through valid transitions.

```joint
stateDiagram-v2
  [*] --> Pending
  Pending --> Planning: configuration accepted
  Planning --> AwaitingApproval: plan and policy pass
  Planning --> Rejected: policy fails
  AwaitingApproval --> Running: approval and identity granted
  AwaitingApproval --> Cancelled: denied or expired
  Running --> Paused: kill switch or operator pause
  Paused --> Running: explicit resume
  Running --> Validating: execution phases pass
  Running --> Failed: phase or timeout fails
  Validating --> Succeeded: service health passes
  Validating --> Failed: validation fails
  Rejected --> Evidencing
  Cancelled --> Evidencing
  Failed --> Evidencing
  Succeeded --> Evidencing
  Evidencing --> [*]
```

Cleanup is a parallel `always` track, not an optional transition. A cleanup failure must appear as its own outcome without obscuring the original cause.

## Control plane availability

- stateless, horizontally scalable API;
- PostgreSQL with PITR and restoration tests;
- durable queue with an idempotency key per command;
- leases to prevent two executors from running the same phase;
- short-lived tokens obtained only when needed;
- latest executable plan exported in signed form for a break-glass scenario;
- metrics and evidence outside the primary failure domain.
