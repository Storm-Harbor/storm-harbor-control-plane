# Project-owned AI, RAG and CLI integrations

## Purpose and current status

AI is optional and advisory. Every project can bring its own AI service, private retrieval API or bounded CLI through a versioned adapter. No AI provider, model family, knowledge service or corporate account is mandatory. The Go integration runtime and plugin registry are planned; current Node scripts only simulate provider selection and make no external calls.

The old Codex/Gemini profile identifiers and model strings in mock fixtures are test data, not verified current model recommendations or evidence of supported provider APIs. Production registries must validate the chosen provider's current capabilities and freeze exact adapter/model versions before activation.

## Integration contract

Each project chooses endpoint/transport, adapter version, identity reference, permitted task types and tools, input/output schemas, timeout, byte/cost limits, data classes, egress and retention policy. A CLI adapter runs in a bounded worker with a fixed executable and argument schema; it never exposes an unrestricted shell.

| Task | Minimum authorized input | Structured proposal |
| --- | --- | --- |
| Inventory classification | selected asset metadata | service/type/criticality with provenance |
| Dependency inference | approved technical relationships | suggested edges and confidence |
| Recovery-plan review | redacted frozen plan | gaps and proposed mitigations |
| Evidence summary | redacted phase results | findings linked to evidence references |
| Private retrieval | a scoped query and authorized context | minimal result references/excerpts allowed by the project |

Model selection, reasoning settings and task responsibility are separate concerns. Adapter capabilities are validated rather than forcing a universal reasoning enum onto every provider. Untrusted results cannot change authorization, controller code, approved targets or execution plans.

## Private knowledge boundary

The project retains its knowledge base. Retrieval can execute inside its environment and return only expressly authorized minimal results. Registration does not upload a corpus, create shared memory or allow another project's queries. Caches and tool state are scoped to organization/project and governed by retention.

Before each external call, enforce data classification, field allowlists, redaction, byte limits, destination policy and authenticated project scope. Secrets, credentials, raw Terraform state and kubeconfigs are prohibited. Retrieved text is untrusted data; embedded instructions cannot grant tool permissions. No alternate provider receives content automatically when the configured provider fails.

Evidence records redacted decisions, selected references, adapter/model version when applicable, latency/usage and digests. It does not persist raw knowledge documents, private internal reasoning or complete prompts by default. A digest alone does not authorize disclosure of its source.

## hub-ai and identity

`hub-ai`/Maya is an optional integration still under development. It uses the same contract as any other project-owned service. Connect to its authenticated, scoped API; do not directly read its databases, vector stores or repository files. `api-auth` is likewise an optional identity integration, not the universal Storm Harbor identity service.

## Execution authority and failure

AI may explain findings and propose a plan or allowlisted tool request. The deterministic Go orchestrator independently verifies actor, effective scope, schema, target, policy, approval, idempotency and kill switch. AI output cannot itself approve a recovery. An approved frozen recovery must run when AI is disabled or unavailable; failures create findings and use deterministic fallback behavior.

## Validation before activation

Test cross-project retrieval denial, prompt/response leakage canaries, output schema/reference validation, bounded CLI arguments and egress, injected instructions in retrieved text, deterministic behavior without AI, cost/time limits and compatibility across adapter upgrades. Private-corpus evaluations stay in the owning project's authorized environment.
