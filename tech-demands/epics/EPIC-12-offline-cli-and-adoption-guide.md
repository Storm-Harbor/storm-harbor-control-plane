# Epic 12 — Offline CLI and adoption guide

Native owner: [Epic #12](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/12). Parent: [Journey #11](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/11), [durable record](../journeys/JOURNEY-11-offline-recovery-configuration-validation.md). GitHub owns current lifecycle and linked PR evidence.

## Deliveries and estimate

| Native leaf | Acceptance | Focused hours |
| --- | --- | --- |
| [Core CLI #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14) | Version/help, bounded regular-file validation, local expected scope, safe output/exit behavior, raw input fingerprint, executable tests, module/CI/platform and release gates | 16 |
| [Adoption docs #13](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/13) | Exact candidate usage, current-versus-published boundary, error/exit reference, native records, strict documentation and navigation acceptance | 4 |

Pure umbrella total: 20 hours; no implementation overhead is counted again. Both leaves now belong to the primary control-plane repository following the owner's 2026-11-13 decision. The originally created CLI issue was transferred with its native identity/history and project linkage; the unused auxiliary repository was archived without runtime code.

## Integration and exit

Consume immutable `storm-lib-recovery-contracts v0.1.0-alpha.1`; never use a local `replace`, producer source import or moving branch. Configuration semantics stay in that producer. CLI report/output is a presentation contract and grants no authority. No agent is mandatory; future opt-in Copilot/Codex/Claude Code/private-agent analysis uses scoped proposals and separate gates.

User manual merge, accepted main checks, actual executable/platform evidence, immutable binary release integrity/provenance and published English adoption guide are required before completing the Epic. Preserve incomplete coverage dimensions and missing publication gates explicitly. A candidate CI artifact is not a released executor or evidence of recovery.
