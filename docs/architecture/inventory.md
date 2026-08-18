# Inventory and dependency graph

## Central unit: service

The inventory starts with the business service, not the cloud resource. Technical assets exist to explain whether that service can be recovered.

```joint
flowchart LR
  S[checkout-service<br/>critical | RTO 60 | RPO 15]
  R[Git repositories]
  T[Terraform modules]
  P[GCP project]
  N[VPC and DNS]
  K[GKE cluster]
  W[Namespace and workloads]
  D[Cloud SQL]
  B[Recovery points]
  G[GitOps objects]

  S --> R
  R --> T
  T --> P
  P --> N
  N --> K
  K --> W
  S --> D
  D --> B
  R --> G
  G --> W

  classDef service fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:3px;
  classDef code fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef cloud fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef data fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  class S service;
  class R,T,G code;
  class P,N,K,W cloud;
  class D,B data;
```

## Initial relational model

PostgreSQL is sufficient for the MVP. The graph is represented by entities and typed relationships, with recursive queries when necessary.

### SQLite in GitHub?

A SQLite file can technically be committed, but it should not be the product's runtime database:

| Use | Decision | Reason |
| --- | --- | --- |
| local development | allowed | simple setup and a single writer |
| single-node appliance | allowed with external backup | deliberately nonconcurrent operation |
| rebuildable discovery cache | allowed outside Git | loss does not affect the authoritative source |
| simulation artifact | allowed if redacted | immutable snapshot for download and diagnosis |
| shared inventory | PostgreSQL | concurrency, transactions, backup, HA, and queries |
| `.db` file versioned in Git | not recommended | binary with no diff/merge and with locking or concurrent-write issues |

GitHub remains suitable for declarative YAML/JSON, schemas, policies, and recovery plans. A pipeline should not modify and commit SQLite back to the repository: besides causing conflicts, that mixes runtime state with desired state and requires unnecessary write permission.

| Table | Essential fields |
| --- | --- |
| `tenants` | id, name, policy_set_id |
| `services` | id, name, tier, criticality, rto, rpo, owner_id |
| `environments` | id, service_id, name, classification |
| `assets` | id, provider, type, canonical_ref, source, observed_at |
| `dependencies` | source_id, target_id, type, criticality, confidence |
| `repositories` | installation_id, role, owner, name, path, revision |
| `recovery_plans` | service_id, version, source_revision, digest |
| `recovery_points` | asset_id, type, created_at, integrity_status |
| `recovery_runs` | plan_id, state, mode, started_at, finished_at |
| `phase_runs` | run_id, phase, attempt, state, metrics, evidence_uri |
| `findings` | service_id, rule_id, severity, status, first_seen, last_seen |

## Relationships

`dependency_type` needs operational semantics:

- `HOSTED_ON`: workload → cluster;
- `DEPENDS_ON`: service → database;
- `PROVISIONED_BY`: asset → Terraform module;
- `RECONCILED_BY`: workload → GitOps object;
- `AUTHENTICATES_WITH`: workload → service account;
- `NETWORK_REACHES`: workload → endpoint;
- `RECOVERED_FROM`: asset → recovery point;
- `OWNED_BY`: service → team.

Each relationship has `source`, `confidence`, `observed_at`, and, when inferred, `explanation`. Low-confidence critical relationships generate a human confirmation task.

## Readiness by service

The score helps with prioritization but must never hide binary controls.

| Dimension | Example gate |
| --- | --- |
| Coverage | all critical assets are linked to a plan |
| Recoverability | recovery point is valid and within the RPO |
| Reproducibility | sources and modules are resolved to immutable revisions |
| Access | recovery identity and permissions are validated |
| Test recency | latest exercise falls within the defined cadence |
| Performance | measured RTO/RPO are within the objective |
| Evidence | bundle is intact, signed, and retained |

Example presentation: `Readiness 84/100 — BLOCKED`, because a single mandatory gate (`database recovery point integrity`) is red. The score does not turn a block into approval.
