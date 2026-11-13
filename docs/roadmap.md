# Roadmap

The [platform contract](architecture/platform-contract.md) defines the accepted direction. Go 1.26.x is the runtime baseline; technical documentation stays in English. No production recovery engine, cloud adapter, extension registry or Kubernetes Operator exists yet.

The [complete control-plane plan and checklists](SPECS-CONTROL-PLANE.md) breaks this roadmap into implementation/adoption stages with prerequisites and acceptance gates.

| Milestone | Verifiable outcome | Current state |
| --- | --- | --- |
| M0: public foundations | website, English docs, schema/mock checks, clear status and community policies | website and English docs published; safe simulations available |
| M1: Go declarative core | versioned provider-neutral resources, schema validation, Go SDK/CLI and configuration migration | initial offline configuration library released as v0.1.0-alpha.1; CLI/migration/engine still planned |
| M2: private onboarding and inventory | organization/project isolation, configurable identity, authorized private sources and dependency graph | planned |
| M3: adapters and AI integration | capability/version registry; independent cloud adapters and optional project-owned AI/RAG/CLI | planned |
| M4: extensions and Kubernetes | custom kinds/Go controllers, admission/versioning, native CRDs and Operator | planned |
| M5: isolated recovery | disposable targets, verified restore, cleanup under partial failure and independent evidence | planned |
| M6: authorized real recovery | immutable approvals, ephemeral access, deterministic execution, kill switch and runbooks | planned |
| M7: validated multicloud mappings | supported source/target combinations across AWS, Azure, DigitalOcean, OCI and GCP | planned |

Cloud independence and multicloud topology belong in M1 contracts. M7 proves concrete cross-cloud combinations; earlier adapters must not require a GCP-shaped core. Managed-service portability is explicit and tested rather than inferred from provider names.

## Follow current construction

[Getting Started](getting-started/index.md) is the adoption entry point. [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3) and [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4) track the first slice: accepted App/runners/vault/target architecture, adoption documentation and the [Go producer #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1). GitHub owns live issue/PR status. This slice does not install the App, execute workflows against cloud resources or perform real recovery.

The initial producer has passed manual merge, main validation and immutable publication. Next are an independently owned offline Go CLI, the customer-owned private recovery repository, GitHub App authentication/onboarding, runner/vault integration and a GCP adapter with safe sandbox acceptance. Optional repository analysis can propose manifest PRs after the read-only identity and data-policy gates. Infrastructure bootstrap and retained restore/cutover remain separately authorized phases. Documentation and examples evolve with each delivered contract.

Each private recovery repository may also publish an optional [dashboard](getting-started/recovery-dashboard.md) with run history, gap-analysis checklists, phase outcomes and dependency/journey diagrams. Stage 08A follows the scoped reporting/evidence contracts and verifies actual website privacy; it does not wait for the full operational console. Authoritative historical evidence remains outside the failure domain.

## Migration work

- Preserve existing Node simulation fixtures until Go equivalents demonstrate behavior and guardrail parity.
- Replace the GCP-specific configuration through a new version with migration fixtures, not an expanded provider enum alone.
- Replace the current Next.js authentication/BFF runtime with a Go boundary that supports configurable identity and preserves secure session behavior.
- Keep public website generation separate from application runtimes and exclude private/authenticated routes from Pages artifacts.
- Introduce cloud integrations independently with an honest discovery/simulation/recovery support matrix.

## Criteria for leaving mock mode

- reviewed threat model and organization/project isolation tests;
- proven cleanup, cancellation and idempotency during partial failure;
- tested quotas, budget and kill switch;
- minimum permissions and secret references documented for each provider/phase;
- control plane, identity and evidence available outside the modeled failure domains;
- verified source-to-target data, network and service compatibility;
- approved incident/rollback runbooks and measured restore objectives;
- at least three consecutive isolated exercises meeting their defined objectives.

## Community and sustainability

Prepare [sponsorship](community/sponsorship.md) using verified payment destinations, and maintain [CNCF readiness](community/cncf-readiness.md) with actual implementation/adoption evidence. These tracks do not imply active payments, CNCF acceptance or production support.

## Harbor Lens

The [Harbor Lens plan](SPECS-HARBOR-LENS.md) records our own separate frontend, synthetic sandbox, future read-only observation, Harbor groups and independently approved operations. [Journey #16](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/16) starts with specification and synthetic UI. Core Go producers remain here; the product website stays in its existing repository.
