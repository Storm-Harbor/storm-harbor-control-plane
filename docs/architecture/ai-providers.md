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

## Assisted repository onboarding

The planned initial use is a scoped, opt-in analysis of selected repository revisions to propose recovery YAML in a PR. The bot identifies infrastructure/application dependencies and reports missing objectives, backup information and target choices. It does not invent missing facts, execute recovery or retrieve secrets. Operators manually register vault bindings without putting credential values in Git or model context.

The owner's preferred deep-analysis preset is an approved OpenAI GPT model with `xhigh` where the selected agent/model integration supports it. This remains a project choice; it is not a requirement for Copilot, Claude Code or private agents. The actual agent/model must be pinned and evaluated, and unsupported settings cannot trigger a silent fallback. The [master plan's Stage 04A](../SPECS-CONTROL-PLANE.md#stage-04a-optional-assisted-analysis-and-manifest-drafting) owns the complete checklist, source references and read/proposal/execute separation. This early onboarding branch does not make the full AI/RAG runtime or autonomous recovery available.

## Customer-selected analysis agents

The owner explicitly selected agent-independent, opt-in repository analysis on 2026-11-13. Each project can choose its approved agent through a bounded adapter. The choices below are planned integrations, not a list of already working SDK/CLI integrations or licensed account capabilities.

| Planned choice | Required integration boundary |
| --- | --- |
| GitHub Copilot | Approved integration, entitlement/organization policy, actual model/tool capabilities and selected repository access |
| OpenAI Codex | Project-approved supported integration, explicit identity, fixed tool/input/output policy and verified data destination |
| Claude Code | Project-approved supported integration, explicit identity, fixed tool/input/output policy and verified data destination |
| Customer-owned agent or private AI/RAG/CLI | The same scoped proposal contract, isolation, bounded tools and customer data/retention controls |

The analysis contract must bind authenticated organization/project scope, selected repositories/revisions/paths, task, allowed input data, execution placement, agent identity/version, tool allowlist, timeout/byte/cost ceilings, authorized destination and a structured manifest/findings output. A project may keep analysis on a private runner. Repository read access does not by itself authorize sending source or knowledge to an external agent service.

Storm Harbor supplies only the permitted inputs/tools and independently validates proposals. Installing or choosing an agent does not grant access beyond the bot's repository allowlist, enable local shell access, import ambient developer credentials/plugins or authorize cloud/vault retrieval. Where an agent cannot satisfy a required boundary, that adapter remains unavailable for the task. Exact authentication, supported models, settings and capabilities are verified against the selected integration's official documentation and account policy during implementation/activation; none is assumed interchangeable.

The reviewed project profile selects one adapter explicitly. There is no automatic Copilot-to-Codex/Claude/private-provider fallback. Record safe provenance and capability decisions, not unrestricted prompts or internal reasoning. Secret bindings remain operator-managed; candidate manifests and findings go through the same constrained PR/human-review path for every agent. The current `RecoveryConfiguration` alpha accepts neither `spec.ai` nor `spec.agent`; the future versioned analysis profile is a separate contract.

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
