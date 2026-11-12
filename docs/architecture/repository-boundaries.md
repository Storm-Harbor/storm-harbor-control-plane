# Product and repository boundaries

Storm Harbor is delivered as independently versioned applications, services, and libraries. Repository names identify ownership; they do not create a second source of truth for identity, cloud state, or recovery evidence.

The [platform contract](platform-contract.md) is the accepted target. New executable components use Go 1.26.x. The existing Next.js BFF and Node.js simulations remain transitional until a tested migration replaces them. The static product website and its frontend build tools are separate from production runtime services.

## Current repositories

| Repository | Owns | Does not own |
| --- | --- | --- |
| `storm-harbor-app` | public product website, login journey, MFA challenge UI, and the authenticated operational console | credentials, MFA secrets, recovery domain execution, or authoritative inventory |
| `storm-harbor-control-plane` | this technical documentation site, architecture, schemas, safe simulations, and executable contract checks | the customer console or a production runtime API |
| `api-auth` | optional iot.EraX identity integration currently used by the legacy BFF | a mandatory identity provider for all Storm Harbor adopters; product authorization or recovery execution |
| `hub-ai` / Maya and project-owned AI/RAG/CLI | optional, scoped intelligence behind a versioned adapter; hub-ai is still under development | mandatory AI, replicated customer knowledge, unrestricted cloud access or deterministic recovery authority |

The Storm Harbor runtime will be introduced through independent repositories as its bounded capabilities become executable.

## Naming convention

Reusable Storm Harbor packages use the `storm-lib-*` repository and package-family prefix. Independently deployable microservices use the `storm-*` prefix.

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

Maya is an optional operational intelligence integration from [iot.EraX](https://ioterax.com), developed in `hub-ai`. Another project can bring its own AI, private RAG API or bounded CLI instead. The integration may correlate findings, explain blocked gates and propose allowlisted tools. Deterministic orchestration and policy remain authoritative. Knowledge stays in the project's chosen environment; retrieval results and outbound data require explicit scope and data policy.

Every Maya-initiated action must include:

- authenticated actor and effective tenant;
- bounded, typed tool and validated parameters;
- minimum short-lived permission;
- idempotency key, timeout, and kill-switch behaviour;
- human approval when policy or impact requires it;
- correlated input, decision, action, validation, and evidence.

Maya never receives an unrestricted shell, bypasses application authorization, or accesses databases and vector stores directly.
