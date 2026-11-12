# Epic 4 — Private configuration and execution contracts

Native owner: [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4). Parent: [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3), [durable record](../journeys/JOURNEY-3-secure-declarative-recovery-foundation.md). GitHub owns current lifecycle and PR links. Estimate: 20 focused engineering hours = documentation 8 + producer 12; no duplicated parent implementation effort.

## Deliveries

| Owning issue | Contract and acceptance |
| --- | --- |
| [Documentation #5](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/5) | English Getting Started, exact YAML fixtures, policies/journey, GCP target creation/reuse/placement, App/runners/external secrets architecture, accepted decisions, strict documentation build |
| [storm-lib-recovery-contracts #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1) | New `RecoveryConfiguration` version, trusted scope equality, logical secret aliases, exercise/restore and create/existing lifecycle invariants, strict bounded decoder, safe typed errors and ten catalogs, schema conformance, race/coverage/vet/build/module evidence and immutable publication |

## Decisions and exclusions

Use Go 1.26.x and a provider-independent domain; separate configuration, identity and secret storage. Selected private repository access does not authorize product/cloud operations. Installation policy cannot be relaxed by YAML. The security strategy is owner-approved; further strategy approval is not a blocker to preparing these deliveries. Manual PR merge remains required.

No cloud/vault/production job execution, GitHub App hosting, CLI consumer or automatic project provisioning is introduced here. Configuration checks are not membership authentication, plan approval, provider capability verification or actual source protection. Customer identifiers stay in private local reference notes.

## Integration and exit

Published documentation must describe current versus planned capabilities accurately, preserve existing simulation fixtures and cross-link the module owner and native hierarchy. Future runtime consumers adopt only a validated immutable producer release after manual merge. Evidence checkpoints carry candidate revision/check URLs and remaining merge/release gates. Close the Epic only after both children and all publication acceptance are complete.
