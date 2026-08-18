# Security architecture

## Trust boundaries

```joint
flowchart LR
  U[User and approver]
  GH[GitHub App]
  CP[DR Control Plane]
  WIF[Workload Identity Federation]
  GCP[GCP protected resources]
  K8S[Kubernetes API]
  VAULT[Secret Manager and KMS]
  EVID[Immutable evidence]

  U --> GH
  GH --> CP
  CP --> WIF
  WIF --> GCP
  WIF --> K8S
  CP --> VAULT
  CP --> EVID

  classDef identity fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef control fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef target fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef proof fill:#eeeeee,stroke:#666666,color:#1d1d1b,stroke-width:2px;
  class U,GH,WIF identity;
  class CP control;
  class GCP,K8S,VAULT target;
  class EVID proof;
```

## Mandatory controls

| Control | Expected implementation |
| --- | --- |
| Tenant isolation | tenant ID in every key, policy, and query; tests against IDOR |
| Authentication | GitHub App and corporate identity; short sessions |
| Cloud access | OIDC/WIF, no persistent keys, restricted audience and subject |
| Authorization | product RBAC + policy by action/resource/environment |
| Approval | requester/approver segregation and approval expiration |
| Secrets | resource references only; redaction before logging/evidence |
| Supply chain | actions and images pinned by digest/SHA, SBOM, and signing |
| Runtime | egress allowlist, network policy, read-only filesystem |
| Evidence | hash, signature, retention, WORM, and an independent project |
| Emergency | global kill switch and audited, time-bound break-glass access |

## Authorization order

A feature flag only enables code. An operation simultaneously requires:

1. the capability enabled in the plan;
2. permission under the installation policy ceiling;
3. a compatible, trusted adapter;
4. the target and labels on the allowlist;
5. a risk/cost/change plan within limits;
6. valid approval for recovery mode;
7. an ephemeral identity authorized for the current phase;
8. the kill switch disabled immediately before mutation.

## IAM snapshot

Snapshots include accounts, roles, bindings, conditions, providers, and Workload Identity, but never private keys. The backup writer and recovery reader should use different identities whenever possible. The encrypted bundle contains the schema version, source revision, digest, and dependency order.

For small payloads, Secret Manager can store versions. For larger snapshots, use versioned Cloud Storage with a retention lock and KMS encryption, keeping only the URI and digest in Secret Manager.
