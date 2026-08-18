# Federated Source of Truth

## The term

**Federated Source of Truth (FSoT)** is a descriptive architectural pattern, not a single formal specification. It describes a model in which different systems remain authoritative in their domains while a coordinating layer provides interoperability, correlation, and a unified view.

Definition adopted by the product:

> Each system remains authoritative for its domain. The DR Control Plane discovers, normalizes, correlates, and reconciles those sources into a unified operational view of recovery readiness.

## Authority matrix

| Entity/field | Authority | Copy in the control plane | Conflict rule |
| --- | --- | --- | --- |
| recovery plan and policy | DR repository | validated document + commit SHA | Git wins; execution freezes the SHA |
| desired infrastructure | Terraform repository | modules, addresses, and SHA | Git wins |
| Terraform lineage/state | remote backend | serial, lineage, digest, and freshness | backend wins |
| desired workloads | GitOps repository | normalized objects and SHA | Git wins |
| existing GCP resource | GCP API/CAI | asset with timestamp and provenance | latest API observation wins |
| existing Kubernetes object | Kubernetes API | UID, generation, status, and timestamp | API wins |
| secret | Secret Manager | resource reference and version only | vault wins; value is never copied |
| evidence | evidence store | index, digest, and URI | immutable object wins |
| owner/RTO/RPO | governed product record or ITSM integration | value and provenance | per-attribute ownership policy |

## Reconciliation

```joint
sequenceDiagram
  participant S as Authoritative source
  participant A as Adapter
  participant N as Normalizer
  participant I as Inventory
  participant P as Policy engine
  participant E as Evidence

  S->>A: snapshot or delta + cursor
  A->>N: provider payload + provenance
  N->>I: canonical assets and relations
  I->>I: upsert by canonical identity
  I->>P: changed service graph
  P->>P: evaluate freshness, drift, RTO and RPO
  P->>E: decision + inputs + rule version
```

## Why not use GitHub for everything

GitHub is excellent for desired state, review, history, and approval. It is not the right place for Terraform state, runtime state, metrics, secrets, distributed locks, or evidence that must survive the loss or compromise of the GitHub organization itself.

The control plane is not “the new truth” either. PostgreSQL acts as a **materialized operational index**. Each record needs:

- `source_system` and `source_resource_id`;
- `observed_at` and `expires_at`;
- `source_revision`, etag, or resource version;
- `content_digest`;
- a confidence level when the relationship is inferred;
- reconciliation state and the latest error.

## Consistency and stale data

The product accepts eventual consistency for inventory, but not for critical decisions. Before a recovery run:

1. freeze commits and versions of declarative sources;
2. force a refresh of the sources used in the plan;
3. reject data beyond the freshness budget;
4. reevaluate policy against the frozen snapshot;
5. assign a digest to the execution plan;
6. record any later change as new drift without silently changing the active run.

!!! warning "Safety rule"
    A missing or unavailable critical source produces an `UNKNOWN` state, never `HEALTHY`. In real recovery, `UNKNOWN` blocks execution unless an explicitly governed break-glass procedure is used.
