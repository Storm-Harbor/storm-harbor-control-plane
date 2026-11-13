# Product and repository boundaries

Storm Harbor is delivered as independently versioned applications, services, and libraries. Repository names identify ownership; they do not create a second source of truth for identity, cloud state, or recovery evidence.

The [platform contract](platform-contract.md) is the accepted target. New executable components use Go 1.26.x. The existing Next.js BFF and Node.js simulations remain transitional until a tested migration replaces them. The static product website and its frontend build tools are separate from production runtime services.

## Current repositories

| Repository | Owns | Does not own |
| --- | --- | --- |
| `storm-harbor-app` | existing product website and transitional login/MFA/application scaffold; the website remains published at `stormharbor.io` | recovery domain execution, authoritative inventory or the future separately delivered monitoring-center app |
| `storm-harbor-control-plane` | architecture/docs, legacy safe simulations and planned primary Go core: CLI, GitHub integration, orchestration, APIs, adapters and controllers/CRDs | product website presentation or the future monitoring-center user interface; planned runtime is not implemented by this ownership decision |
| `harbor-lens` | separate operational frontend and synthetic sandbox; future core API consumer | Go core, product website, live authorization or recovery execution in the first prototype |
| `storm-lib-recovery-contracts` | released alpha Go producer: scoped provider-neutral configuration, strict YAML/schema validation and safe typed errors | authentication, GitHub App hosting, cloud/vault calls, execution or native CRD installation |
| `api-auth` | optional identity integration currently used by the legacy BFF | a mandatory identity provider for all Storm Harbor adopters; product authorization or recovery execution |
| `hub-ai` / Maya and project-owned AI/RAG/CLI | optional, scoped intelligence behind a versioned adapter; hub-ai is still under development | mandatory AI, replicated customer knowledge, unrestricted cloud access or deterministic recovery authority |

The owner selected `storm-harbor-control-plane` as the primary implementation repository on 2026-11-13. It will own the Go CLI, GitHub App/bots, orchestration, APIs, adapters, controllers/CRDs and internal packages. These remain separate domain/application/adapter boundaries inside the repository; multiple core commands or deployment units do not require multiple repositories. Sites and user-facing applications are separate. The product website remains in `storm-harbor-app`, published at `stormharbor.io`. The separate `harbor-lens` application starts with a synthetic sandbox and will consume published core APIs; see the [Harbor Lens plan](../SPECS-HARBOR-LENS.md). GitHub App/bots refer to integrations, not that future user-facing application. Technical source documentation may accompany code; the existing documentation publisher is not moved implicitly.

The already released `storm-lib-recovery-contracts` remains a pinned immutable dependency for compatibility. Preserve its published module paths/tags; any consolidation or extraction needs an explicit migration. No new core repository is created by default.

## Naming convention

Use capability-oriented Go command/package names within the primary repository. Existing externally versioned reusable packages retain their `storm-lib-*` identity. A service name describes a boundary, not a requirement to create another repository.

Examples of capability-oriented names are:

```text
storm-lib-recovery-contracts
storm-lib-provider-sdk
storm-inventory
storm-orchestrator
storm-evidence
storm-policy
```

Names describe a stable capability rather than a cloud vendor. AWS, Azure, DigitalOcean, Oracle Cloud, GCP, Kubernetes, Terraform, GitOps and database support belong behind versioned adapters unless a service has an independent lifecycle and bounded domain.

## Authentication and product access

The current `storm-harbor-app` scaffold uses a backend-for-frontend boundary when calling `api-auth`. Its browser never owns an access token in local or session storage. The BFF stores product sessions and MFA challenges in secure `HttpOnly` cookies and propagates request and correlation identifiers. This is an existing integration, not a universal identity contract. The Go migration must provide a provider-independent identity boundary and preserve session/security behavior.

Storm Harbor accepts any syntactically valid primary email; a corporate domain is not mandatory. Identity authentication alone never grants access to a Storm Harbor organisation or project. Product access requires a separate, authoritative membership, tenant, and RBAC decision.

The frontend may use a cookie-presence check only as a navigation convenience. Every API that exposes customer data or performs an operation must verify the session, tenant, permission, policy, and target independently.

## Maya boundary

Maya is an optional operational intelligence integration developed in `hub-ai`. Another project can bring its own AI, private RAG API or bounded CLI instead. The integration may correlate findings, explain blocked gates and propose allowlisted tools. Deterministic orchestration and policy remain authoritative. Knowledge stays in the project's chosen environment; retrieval results and outbound data require explicit scope and data policy.

Every Maya-initiated action must include:

- authenticated actor and effective tenant;
- bounded, typed tool and validated parameters;
- minimum short-lived permission;
- idempotency key, timeout, and kill-switch behaviour;
- human approval when policy or impact requires it;
- correlated input, decision, action, validation, and evidence.

Maya never receives an unrestricted shell, bypasses application authorization, or accesses databases and vector stores directly.
