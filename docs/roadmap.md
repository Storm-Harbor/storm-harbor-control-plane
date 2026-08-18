# Roadmap

## Product phases

```joint
flowchart LR
  M0[M0<br/>contracts and safe mocks]
  M1[M1<br/>onboarding and inventory]
  M2[M2<br/>continuous discovery]
  M3[M3<br/>Kubernetes Operator]
  M4[M4<br/>disposable recovery tests]
  M5[M5<br/>approval-gated recovery]
  M6[M6<br/>multi-cloud adapters]
  M0 --> M1 --> M2 --> M3 --> M4 --> M5 --> M6

  classDef done fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef next fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef later fill:#f4f4f4,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  class M0 done;
  class M1,M2 next;
  class M3,M4,M5,M6 later;
```

## Definition of done by milestone

| Milestone | Verifiable outcome |
| --- | --- |
| M0 | tested schemas, documentation, safe pipelines, and mock evidence |
| M1 | tenant-isolated installation, allowlist, and editable inventory |
| M2 | periodic reconciliation with provenance, freshness, and drift |
| M3 | versioned CRDs, idempotent Operator, and tested status and finalizers |
| M4 | disposable GCP project created, recovered, validated, and removed |
| M5 | real execution protected by approvals, WIF, policy, and break-glass access |
| M6 | new providers without changes to the central domain model |

## Criteria for leaving mock mode

- reviewed threat model;
- proven isolation and cleanup during partial failure;
- tested quotas, budget, and kill switch;
- minimum permissions documented by phase;
- evidence store outside the failure domain;
- approved incident and rollback runbooks;
- at least three consecutive simulations within the defined objective.
