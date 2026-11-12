# Engineering standards

Storm Harbor adopts reusable engineering practices while remaining independent of any adopter's private services, cloud accounts, identity or business data. The [shared agent guidance](https://github.com/Storm-Harbor/.github/blob/main/AGENTS.md) and its specifications define the project-owned rules.

## Applicable patterns

| Concern | Storm Harbor application |
| --- | --- |
| Contract and mapper ownership | versioned Go contract modules, hexagonal boundaries, pure explicit mappings and allowlisted public projections |
| Error and localization contracts | Go-native typed errors with stable codes/keys and safe arguments; localization at the presentation boundary |
| Technical demands | repository-owned issues, native hierarchy for new major outcomes, complete project metadata, linked PRs and user-controlled manual merges |
| E2E acceptance | actor journeys, negative/isolation/failure cases, exact artifacts and isolated dependencies; mocks do not prove integrated recovery |
| Multitenant isolation | trusted organization/project context inside persistence, search, RAG, cache, files, queues, controllers and evidence access |
| API and workload access | separate user authentication, project authorization and provider workload identity; no implicit privilege from network/account placement |
| Billing safety | sponsorship stays separate from product entitlements; integer amounts, TEST/LIVE separation and verified/idempotent events when required |
| Deployment and telemetry | health/readiness/shutdown contracts, bounded metrics, structured redacted logging and independently available recovery state |

## Deliberately excluded coupling

The project does not inherit another solution's Company aggregate, business registration policy, private package registry, token/header protocol, product catalog, default database, cloud deployment topology or project-board number. Those choices require Storm Harbor-owned contracts and decisions. A GCP deployment example is an adapter scenario, not a universal runtime requirement. A commercial subscription plan is distinct from a disaster-recovery plan.

The core language is Go 1.26.x. The public website retains ten locales and static browser assets; technical documentation is always English. Existing Node simulations and the Next.js authentication/BFF remain transitional, with explicit migration and qualification work. New policies do not turn those existing scaffolds into a measured, compliant runtime.

## Current publication evidence

The initial public website has artifact checks and browser scenarios for all locales, language selection, mobile navigation, static assets/styles, 404 responses and absent authenticated routes. It has no deployed backend/database. That scope does not qualify private onboarding, AI connectors, cloud adapters, controller execution, payment activation or production recovery.

The documentation is validated through a strict MkDocs build. Existing safe simulations retain their contract/guardrail checks and manual read-only execution boundary. GitHub issues and PRs hold exact candidate/CI evidence and pending work.

## Project-owned specifications

- [Contract ownership](https://github.com/Storm-Harbor/.github/blob/main/specs/SPECS-CONTRACTS.md)
- [Delivery and traceability](https://github.com/Storm-Harbor/.github/blob/main/specs/SPECS-DELIVERY.md)
- [Validation and acceptance](https://github.com/Storm-Harbor/.github/blob/main/specs/SPECS-ACCEPTANCE.md)
- [Privacy, access and operational boundaries](https://github.com/Storm-Harbor/.github/blob/main/specs/SPECS-PRIVACY-AND-ACCESS.md)
