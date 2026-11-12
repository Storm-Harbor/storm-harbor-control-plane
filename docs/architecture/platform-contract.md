# Platform contract

This is the accepted product direction. The repository currently contains documentation, a GCP-specific configuration example and safe Node.js simulations. It does not yet implement the Go runtime, multicloud adapters, extension registry or Kubernetes CRDs described here.

The first provider-neutral Go configuration producer is being delivered in [storm-lib-recovery-contracts #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1). It validates intent offline; App/executor integration and recovery remain planned. Follow [Getting Started](../getting-started/index.md), [GitHub execution](github-execution.md), [external secrets](external-secrets.md) and the [accepted decisions](accepted-decisions.md) for this initial slice.

## Independent organizations and private projects

Storm Harbor is a reusable, organization-independent recovery control plane. No specific identity provider, cloud account, source repository or knowledge service is required.

An organization can authorize selected private repositories without making their contents public or transferring them to the Storm Harbor GitHub organization. Every installation, resource, query, task, cache entry, credential reference and evidence bundle must carry an effective organization/project scope. Cross-scope access requires an explicit policy decision and independent authorization checks.

The public project contains source code, schemas and synthetic examples. Customer repositories, resource inventories, prompts, knowledge bases, secrets and operational evidence are private by default. Public telemetry and demonstration artifacts must use synthetic or expressly approved data.

## Cloud-independent and multicloud recovery

The domain models services, environments, resources, dependency edges, recovery points, objectives, policies, plans, runs and evidence independently of a cloud vendor. Provider identifiers remain namespaced opaque references; the core must not require a GCP project ID or another vendor's account structure.

| Adapter target | Provider-specific boundary | Current implementation |
| --- | --- | --- |
| AWS | accounts, regions, identities and service APIs | planned |
| Azure | tenants, subscriptions, regions and service APIs | planned |
| DigitalOcean | teams, projects, regions and service APIs | planned |
| Oracle Cloud Infrastructure | tenancies, compartments, regions and service APIs | planned |
| Google Cloud | organizations, folders, projects, regions and service APIs | GCP-shaped mock contracts only |

One journey may span multiple providers. Recovery from one provider to another is allowed only when the selected adapter versions declare compatible export/import, data, networking, identity and validation capabilities. A generic resource kind does not make managed services interchangeable. Unsupported mappings fail during planning, before approval or mutation.

The control plane and evidence store must survive the selected failure domains. A different project or region is insufficient for an explicitly modeled provider-wide outage; that scenario needs an independently available control plane, identity, backups and evidence path.

## Declarative YAML journey

The versioned YAML contract drives both CLI and guided UI. The UI edits and explains that contract instead of maintaining a second workflow definition.

1. Declare organization/project scope and authorized sources.
2. Connect provider, identity and optional AI adapters using secret references.
3. Discover inventory and confirm ownership and dependencies.
4. Define recovery objectives, targets, policies and permitted resource kinds.
5. Validate schemas and capabilities; resolve all sources and plugins to immutable revisions.
6. Build a dependency-ordered plan and exercise it in an isolated target.
7. Approve the exact plan digest, target, action scope and expiration.
8. Reconcile recovery phases with bounded retries, health gates and cancellation.
9. Verify service behavior and measured RTO/RPO, retain evidence, then continue discovery.

Configuration describes intent; status records observation. Reconciliation must be idempotent and may not repeat an external mutation merely because a status write failed. Credentials are never embedded in YAML. Configuration grants no permission by itself.

The current `dr.config.example.yaml` and schema remain the legacy GCP simulation contract. The provider-neutral resource API requires a new version, migration tooling and compatibility fixtures; changing a provider enum alone would misrepresent support.

## Go runtime baseline

Go **1.26.x** is the official language/release line for new applications, APIs, the control plane, controllers/operators, agents, CLI tools and SDKs. Go repositories must validate formatting, vet, build, OpenAPI where applicable, and meaningful tests with race/coverage collected in one equivalent execution.

Existing Node.js simulation scripts and the Next.js authentication/BFF scaffold are transitional implementations. Their replacement needs behavior and contract parity before retirement. Static website HTML/CSS/JavaScript and frontend build tools do not introduce a Node.js production runtime. Technical documentation remains in English; the public product website supports ten locales.

## Optional project-owned AI

Each project can register its own AI, RAG API or bounded CLI integration behind a versioned contract. The project controls endpoint, identity, permitted tools, data classes, retention and egress. Its knowledge base can stay inside its own environment, returning only authorized minimal results.

No complete knowledge-base replication, cross-tenant retrieval, unrestricted shell, background provider fallback or unapproved external prompt export is allowed. Outputs remain untrusted proposals until schema/reference validation and the applicable review pass. Evidence keeps redacted decisions and provenance, not raw retrieved documents by default.

`hub-ai`/Maya is an optional future integration under development. Neither Maya nor `api-auth` is a prerequisite for another organization to deploy or use Storm Harbor. Deterministic recovery and an already approved frozen plan must operate when AI is disabled or unavailable.

## Custom resources and controllers

Organizations can extend the declarative model with versioned resource definitions and Go controllers. See [Resource extensibility](extensibility.md). Native Kubernetes CRDs belong to the Kubernetes adapter; the external control plane does not require Kubernetes as its persistence or API substrate.
