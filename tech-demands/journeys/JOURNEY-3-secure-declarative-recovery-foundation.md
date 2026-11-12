# Journey 3 — Secure declarative recovery foundation

Native owner: [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3). GitHub owns lifecycle, assignees, project fields and PR links. Estimate: 20 focused engineering hours, summed from the direct child Epic; no umbrella overhead is counted twice.

## Outcome and actors

An organization maintainer describes private recovery intent in a central repository, chooses exercise or restore and a new/existing target/placement, and understands how the future GitHub App and Go executors will consume that declaration. Security administrators own external trust, vaults and approvals. Infrastructure/application owners retain their resource and release authority. Readers can follow the build through English Getting Started pages.

## Boundary and contracts

The first slice records accepted architecture and delivers offline validation in `storm-lib-recovery-contracts`. It preserves provider independence, explicit trusted organization/project scope, logical secret references and safe target lifecycle semantics. The legacy GCP simulation contract stays compatible. No production App, CLI/executor, cloud adapter, vault retrieval, infrastructure apply, restore or traffic cutover is delivered by this Journey.

Private repositories, customer cloud identifiers, vault paths, state and knowledge never become public examples. Authorized neighboring workspace references inform generic Terraform/state, application-release and Go-contract patterns; they are not mandatory integrations. The expressly excluded workspace is not inspected or changed.

## Hierarchy and delivery order

- [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4): [private configuration and execution contracts](../epics/EPIC-4-private-configuration-and-execution-contracts.md).
- Epic leaf [documentation #5](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/5): accepted decisions, adoption, YAML/policies/journey and GCP target design.
- Epic leaf [producer #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1): Go contract, schema, errors/catalogs and validation.

Documentation and producer preparation proceed independently. Future executable consumers require the producer's validated manual merge and immutable publication. PRs are merged manually by the user.

## Acceptance and evidence

Require reciprocal native hierarchy and complete Project #1 fields; safe synthetic examples; strict MkDocs build; scope/parser/lifecycle negative tests; combined Go race/statement-coverage ≥97%, formatting/vet/build/checksums; pinned CI; manual merge; immutable producer publication; issue/Project reconciliation. Go native coverage does not establish independently measured line or branch coverage.

Exact revisions, executed results and release blockers belong in owning issue checkpoints and linked checks. This record receives a completion marker only after all child acceptance and publication gates are satisfied.
