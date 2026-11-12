# Product overview

## Definition

Storm Harbor is an **organization-independent, declarative recovery control plane**. It continuously assesses whether services can meet RTO/RPO, plans and exercises recovery, and is designed to restore environments through authorized, auditable reconciliation.

The [platform contract](../architecture/platform-contract.md) makes cloud independence, private projects, optional project-owned AI/RAG/CLI, Go 1.26.x runtime components and YAML journeys explicit. AWS, Azure, DigitalOcean, Oracle Cloud and GCP are adapter targets; production support for these providers is not implemented in this scaffold.

## Functional scope

```joint
flowchart TB
  E[Experience<br/>portal, API, CLI, and GitHub App]
  G[Governance<br/>services, owners, tiers, RTO, and RPO]
  C[Control plane<br/>inventory, graph, policy, and orchestration]
    A[Versioned adapters<br/>sources, clouds, data, Kubernetes, optional AI]
  R[Recovery targets<br/>foundation, data, cluster, and workloads]
  P[Proof<br/>metrics, logs, artifacts, and evidence]

  E --> G --> C --> A --> R --> P

  classDef experience fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef governance fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef control fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef execution fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  classDef proof fill:#eeeeee,stroke:#666666,color:#1d1d1b,stroke-width:2px;
  class E experience;
  class G governance;
  class C control;
  class A,R execution;
  class P proof;
```

## ITSM/ITSCM scope

The product uses service-continuity practices without attempting to replace a complete ITSM suite.

| Included | Outside the core |
| --- | --- |
| minimum service and owner catalog | complete service desk |
| criticality, impact, and recovery tier | general ticket management |
| dependency mapping | enterprise-wide CMDB |
| RTO, RPO, and recovery plans | every ITIL process |
| exercises, outcomes, and exceptions | complete change management |
| readiness, evidence, and auditing | IT financial management |

Future integrations may open incidents, record changes, and synchronize CMDB attributes while maintaining clear ownership for each field.

## Principles

1. **Federated:** each domain maintains its authoritative source.
2. **Outside the failure domain:** the orchestrator survives the recovered environment.
3. **Read-only before mutation:** discovery and planning precede privileged tokens.
4. **Fail closed:** missing evidence, policy, or identity blocks execution.
5. **Idempotent:** repeating a phase produces the same desired state.
6. **Evidence first:** every decision and transition remains explainable after an incident.
7. **Extensible:** providers and engines integrate through versioned adapter contracts.
8. **Private by default:** each organization retains its repositories, knowledge and operational data.
9. **Declarative:** UI and CLI follow the same versioned YAML/resource contracts.
10. **Vendor-independent:** cloud, identity and AI choices belong behind adapters.

## Intended experience

A user authorizes selected private repositories, connects cloud accounts through scoped adapters, and optionally connects project-owned AI/RAG/CLI. The guided UI and CLI share a versioned YAML configuration: discover, confirm inventory, set objectives, plan, validate, simulate, approve, recover, verify and reconcile. Source revisions and adapter digests are frozen for every run. Cross-cloud recovery requires explicit compatibility and health checks; real recovery is never enabled implicitly.
