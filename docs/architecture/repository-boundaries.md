# Product and repository boundaries

Storm Harbor is delivered as independently versioned applications, services, and libraries. Repository names identify ownership; they do not create a second source of truth for identity, cloud state, or recovery evidence.

## Current repositories

| Repository | Owns | Does not own |
| --- | --- | --- |
| `storm-harbor-app` | public product website, login journey, MFA challenge UI, and the authenticated operational console | credentials, MFA secrets, recovery domain execution, or authoritative inventory |
| `storm-harbor-control-plane` | this technical documentation site, architecture, schemas, safe simulations, and executable contract checks | the customer console or a production runtime API |
| `api-auth` | global iot.EraX identity, credentials, email verification, MFA, and token issuance | Storm Harbor organisations, projects, roles, recovery plans, or cloud permissions |
| Maya Core and its Storm Harbor context | provider-independent operational intelligence, governed tools, memory, and explanations | direct storage access, unrestricted cloud access, or deterministic recovery authority |

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

Names describe a stable capability rather than a cloud vendor. GCP, Kubernetes, Terraform, GitOps, and database support belong behind versioned provider adapters unless a service has a genuinely independent lifecycle and bounded domain.

## Authentication and product access

`storm-harbor-app` uses a backend-for-frontend boundary when calling `api-auth`. The browser never owns an access token in local or session storage. The BFF stores product sessions and MFA challenges in secure `HttpOnly` cookies and propagates request and correlation identifiers.

Storm Harbor accepts any syntactically valid primary email; a corporate domain is not mandatory. Identity authentication alone never grants access to a Storm Harbor organisation or project. Product access requires a separate, authoritative membership, tenant, and RBAC decision.

The frontend may use a cookie-presence check only as a navigation convenience. Every API that exposes customer data or performs an operation must verify the session, tenant, permission, policy, and target independently.

## Maya boundary

Maya is the operational intelligence layer from [iot.EraX](https://ioterax.com). She may correlate findings, explain blocked gates, form hypotheses, and propose or invoke allowlisted recovery tools. Deterministic orchestration and policy remain authoritative.

Every Maya-initiated action must include:

- authenticated actor and effective tenant;
- bounded, typed tool and validated parameters;
- minimum short-lived permission;
- idempotency key, timeout, and kill-switch behaviour;
- human approval when policy or impact requires it;
- correlated input, decision, action, validation, and evidence.

Maya never receives an unrestricted shell, bypasses application authorization, or accesses databases and vector stores directly.
