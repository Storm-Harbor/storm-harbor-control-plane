# Epic 4 — Private configuration and execution contracts

Native owner: [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4). Parent: [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3), [durable record](../journeys/DONE-JOURNEY-3-secure-declarative-recovery-foundation.md). GitHub owns current lifecycle and PR links. Estimate: 20 focused engineering hours = documentation 8 + producer 12; no duplicated parent implementation effort.

Completed: **2026-11-12T17:05:03Z**, after both native children and their publication/acceptance gates were verified. Full future stages are tracked separately in the [master plan](../../docs/SPECS-CONTROL-PLANE.md).

## Deliveries

| Owning issue | Contract and acceptance |
| --- | --- |
| [Documentation #5](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/5) | English Getting Started, exact YAML fixtures, policies/journey, GCP target creation/reuse/placement, App/runners/external secrets architecture, accepted decisions, strict documentation build |
| [storm-lib-recovery-contracts #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1) | New `RecoveryConfiguration` version, trusted scope equality, logical secret aliases, exercise/restore and create/existing lifecycle invariants, strict bounded decoder, safe typed errors and ten catalogs, schema conformance, race/coverage/vet/build/module evidence and immutable publication |

## Decisions and exclusions

Use Go 1.26.x and a provider-independent domain; separate configuration, identity and secret storage. Selected private repository access does not authorize product/cloud operations. Installation policy cannot be relaxed by YAML. The security strategy is owner-approved; further strategy approval is not a blocker to preparing these deliveries. Manual PR merge remains required.

No cloud/vault/production job execution, GitHub App hosting, CLI consumer or automatic project provisioning is introduced here. Configuration checks are not membership authentication, plan approval, provider capability verification or actual source protection. Customer identifiers stay in private local reference notes.

## Integration and exit

Published documentation describes current versus planned capabilities, preserves the legacy simulation fixtures and cross-links the module owner and native hierarchy. [v0.1.0-alpha.1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/releases/tag/v0.1.0-alpha.1) is published immutably at accepted commit `315acfaaf13e7606525476764f204566682c2ab6`; the public Go proxy origin and signed module/go.mod checksums are verified. Both children and this Epic are closed with native Project Done. Future runtime consumers can pin this release only within their own scoped delivery and acceptance gates.
