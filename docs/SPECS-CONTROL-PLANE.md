# Storm Harbor Control Plane — Plan and checklist

This is the implementation and adoption plan for Storm Harbor. It connects the product decisions, delivery stages, security boundaries and acceptance checklists so maintainers and adopters can follow the work from configuration validation to verified recovery.

**Evidence baseline: 2026-11-12.** A checked item means the stated outcome has supporting evidence. An unchecked item is required or planned work, not an available capability. GitHub issues, native parent/sub-issue relationships, linked PRs and releases remain authoritative for live status; this document is an acceptance ledger, not a second project board. Update a checkbox only with its owning issue and concrete validation/publication evidence.

The current planning delivery is [documentation #9](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/9). The initial foundation is recorded by [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3), [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4), [adoption documentation #5](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/5) and [Go producer #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1).

## 1. Product outcome and boundaries

An organization describes recovery intent in its own private repository. Storm Harbor validates that intent, discovers authorized state, builds a frozen plan, obtains the required approval and executes bounded recovery phases in the organization's permitted environment. It verifies service/data outcomes and retains evidence outside the failure domain.

- Organizations, projects, repositories, infrastructure and knowledge stay private by default.
- GitHub holds reviewed desired state. Cloud APIs, Terraform state backends, backup stores, vaults and evidence stores retain authority over their own domains.
- The design is independent of cloud, identity provider, AI provider and a mandatory Kubernetes installation.
- AWS, Azure, DigitalOcean, Oracle Cloud Infrastructure and GCP are adapter targets. Cross-cloud recovery requires tested source-to-target compatibility.
- Go **1.26.x** is the runtime baseline for SDKs, CLIs, applications, APIs, agents and controllers. Browser assets/build tools and the existing transitional simulations have their documented boundaries.
- Recovery is a versioned YAML journey. User-defined resource kinds and controllers are part of the target design; native Kubernetes CRDs belong to the Kubernetes adapter.
- Each private recovery repository may have its own optional dashboard: run history, gap-analysis checklists, phase outcomes and dependency/journey diagrams using the Storm Harbor visual identity. Publication is a scoped projection of evidence, with independently verified audience and hosting controls.
- AI/RAG/CLI integrations are optional, project-owned and advisory. They cannot authorize an operation, change the approved plan or gain access to another project's knowledge.
- Technical documentation is English. The public website maintains all configured locales, currently 25 language/regional options.

### Current implementation boundary

The public website, English documentation, safe simulations and initial offline Go configuration library are delivered. The GitHub App, executable recovery CLI, live cloud/vault adapters, durable orchestrator, production recovery API and installable CRDs are future deliveries. A valid YAML document is not proof of authenticated membership, provider permission or recoverability.

## 2. Stages and dependencies at a glance

Stages identify acceptance boundaries, not calendar promises. Independent work can proceed together when it does not consume an unpublished producer. Security, evidence, documentation and supply-chain controls begin with the first executable component and remain required throughout.

| Stage | Outcome | Prerequisites | Evidence baseline |
| --- | --- | --- | --- |
| 00 | Public foundations and engineering ownership | None | Delivered |
| 01 | Initial versioned Go configuration producer | 00 | Alpha released and verified |
| 02 | Offline Go CLI and immutable executable artifacts | 01 | [CLI source merged #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14); signed binary release/qualification pending |
| 03 | Private adopter recovery repository and reviewed configuration | 02 | Approved adoption step; repository not created yet |
| 04 | GitHub App installation, events and trusted product scope | 01; 03 for pilot onboarding | Planned |
| 04A | Optional repository analysis and manifest proposals | 02–04; explicit model/data/tool-policy gates | Planned early onboarding branch |
| 05 | Bounded workflow execution and runner/workload identity | 02, 04 | Planned |
| 06 | External bindings and secret resolution | 04, 05 | Planned |
| 07 | Read-only inventory, dependencies and readiness | 04–06 for protected sources | Planned |
| 08 | Policy, immutable planning and durable orchestration | 02, 04–07 | Planned |
| 08A | Optional repository recovery dashboard and historical reports | 03, 07–08; approved projection and publication audience | Planned; independent of the operational console |
| 09 | GCP target creation/reuse and infrastructure-owner integration | 05, 06, 08 | Design documented; adapter planned |
| 10 | Isolated end-to-end recovery exercise | 03–09 and independent evidence path | Planned |
| 11 | Authorized retained restore, cutover and emergency execution | 10 and approved incident controls | Planned |
| 12 | Extensible resources, controllers and Kubernetes integration | Versioned admission/status/run contracts from 08; safe lifecycle evidence | Planned |
| 13 | Additional cloud adapters and tested multicloud mappings | 08–11 patterns; acceptance per provider/mapping | Planned |
| 14 | Optional private AI/RAG/CLI integration | Trusted scope and deterministic tools from 04–08 | Planned; not a recovery prerequisite |
| 15 | Go operational APIs and qualified customer console | Published domain/run/evidence contracts; working workflows | Planned; static website already delivered |
| 16 | Operational qualification and supported releases | Cross-cutting controls from 02 onward; relevant completed capabilities | Planned production qualification |
| 17 | Community, funding and CNCF readiness | Independent of recovery execution | Foundations delivered; activation/submission separate |

The first usable execution milestone is a **bounded, isolated GCP exercise through a private adopter repository**. It does not require waiting for a console, AI integration or every cloud adapter. A retained real restore has additional gates and must not be inferred from a successful exercise.

Stage ranges refer to the numbered core stages. Letter-suffixed Stages 04A and 08A are optional branches and do not become prerequisites for recovery merely because they appear between core stages.

Stage 04A is an optional early use of a bounded AI adapter for read-only repository analysis and reviewed manifest proposals. It does not wait for the full recovery engine and does not enable mutations. Stage 14 extends the model to broader private retrieval/operational assistance after deterministic tool and orchestration contracts exist.

## 3. Ownership and contract evolution

| Owner | Existing responsibility |
| --- | --- |
| `storm-harbor-control-plane` | Primary repository for Go CLI, GitHub App/bots, orchestration, APIs, adapters, controllers/CRDs and internal packages; also technical documentation, legacy simulations and demand records |
| `storm-lib-recovery-contracts` | Initial shared Go configuration semantics, strict YAML/schema validation and typed errors |
| `storm-harbor-app` | Static multilingual product website and transitional application scaffold |
| `harbor-lens` | Separate operational frontend and synthetic sandbox; future scoped core API consumer |
| `.github` | Public profile, engineering guidance, contribution/security policies and delivery practices |
| Customer-owned private repository | Recovery declarations, reviewed policies/runbooks and authorized source references; never shared customer knowledge or secret values |

The owner selected `storm-harbor-control-plane` as the primary repository for core implementation on 2026-11-13. Organize Go commands and internal domain/application/adapter packages there; separate runtime components do not imply separate repositories. Sites and user-facing applications remain separate. The product website stays in `storm-harbor-app` at `stormharbor.io`; a future monitoring-center application will have its own delivery and consume the control-plane APIs. GitHub App/bots are core integrations, distinct from that user-facing application. Existing technical source documentation can accompany implementation; this decision does not silently migrate its current publisher. The already published `storm-lib-recovery-contracts` version remains an immutable compatibility dependency. New core repository extraction requires a concrete separately approved need; it is not the default plan. Native owning issues and producer publication gates still apply.

### Initial alpha versus planned contract families

| Contract family | Available in the initial alpha | Further requirements |
| --- | --- | --- |
| `RecoveryConfiguration` | Scope, purpose, source alias, target intent and logical secret aliases | Explicit version/migration for incompatible changes |
| Scope | Identifier shape and equality to independently supplied trusted scope | Authentication, authoritative membership, revocation and operation authorization |
| Provider/placement | Logical aliases | Customer-controlled binding, actual cloud identity, environment/region/network and capability verification |
| Secrets | Alias, scope and bounded purpose | Vault/version binding, workload identity, least privilege and runtime resolution |
| Journey/dependency graph | Not implemented | Versioned steps, dependencies, health gates, cancellation and compensation semantics |
| Policy/approval | Not implemented | Installation ceilings, immutable plan/target binding, expiry and independent enforcement |
| Run/evidence | Not implemented | Durable state, idempotency, checkpoints, audit, integrity and authorized retrieval |
| Custom kinds/controllers/AI profiles | Not implemented | Admission, ownership, version conversion, isolation and bounded tools |

Future contract names and fields are design work, not accepted extensions to the current document. In particular, `spec.policy`, `spec.journey`, `spec.ai`, raw vault paths, inline credentials and arbitrary `status` fields are not accepted by the initial strict alpha parser. Preserve the separate legacy GCP simulation schema until an explicit migration is qualified.

## 4. Stage checklists

### Stage 00 — Public foundations

**Outcome:** contributors and adopters can discover the project, understand its current limits and find the owning work.

- [x] Publish [stormharbor.io](https://stormharbor.io/) as a static product website with 25 locale options and RTL support for Arabic, Hebrew and Urdu.
- [x] Publish English [Getting Started](getting-started/index.md), YAML, policies, journey, GCP-target and external-secret guides.
- [x] Correct documentation navigation for desktop/mobile, keyboard access, instant navigation and storage failure; publish refreshed custom assets.
- [x] Update the [organization profile](https://github.com/Storm-Harbor) for independent branding and current repository boundaries.
- [x] Establish shared engineering, contribution, security, contract/delivery/acceptance and private-scope guidance.
- [x] Use native owning issues and the maintainer project with complete assignee, type, labels, priority, size, focused-hour estimate and status fields.

**Exit evidence:** website/docs deployments and live URL checks; 25 locale URLs with correct language/direction, 35 website browser scenarios, seven docs navigation scenarios, public profile and link validation. GitHub owns the detailed checkpoints for [website #4](https://github.com/Storm-Harbor/storm-harbor-app/issues/4), [docs #7](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/7) and [profile #4](https://github.com/Storm-Harbor/.github/issues/4).

### Stage 01 — Initial Go configuration contracts

**Outcome:** later consumers can pin a single published owner for offline configuration meaning and rejection rules.

- [x] Implement `recovery.stormharbor.io/v1alpha1`, kind `RecoveryConfiguration`, independently of provider account structures.
- [x] Validate organization/project shape and equality against a caller-supplied trusted context; document that validation does not authenticate that context.
- [x] Keep secrets as scoped logical references without credential values or vault paths.
- [x] Separate exercise/restore, create/existing, placement and target retention/expiry/cleanup intent.
- [x] Reject unknown/duplicate/case-mismatched fields, anchors/aliases, multiple documents, unsafe types and oversized/deep YAML without echoing values.
- [x] Deliver JSON Schema, synthetic fixtures, typed safe errors and ten presentation error catalogs.
- [x] Pass combined Go race/statement coverage (98.3%), schema/scope/lifecycle/catalog checks, formatting, vet, build, Windows cross-build and module metadata checks.
- [x] Manually merge the producer and publish immutable [v0.1.0-alpha.1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/releases/tag/v0.1.0-alpha.1) from accepted commit `315acfaaf13e7606525476764f204566682c2ab6`.
- [x] Download the published module through the public Go proxy, verify its origin commit/tag and match the signed Go checksum database record.

**Exit evidence:** [producer #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1), accepted [main CI](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/actions/runs/43242888926) and the immutable prerelease. Go native profiles measure statement coverage and function-level statement coverage; independent branch/line coverage is not claimed.

Published module checksum: `h1:wVFOP9bwkppvtz6VfAXn1dMaOZMT59NsKLH5xNq8O4A=`. Published `go.mod` checksum: `h1:axlq5Hn6XXSUaB3YG3cngnETpHVTh23/TNxZroiKU30=`.

### Stage 02 — Offline Go CLI

**Outcome:** a maintainer can validate a recovery declaration locally or in a credential-free CI job using an immutable executable release.

First source slice: [Journey #11](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/11) → [Epic #12](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/12) → [CLI #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14) and [documentation #13](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/13). The [CLI guide](getting-started/offline-cli.md) documents candidate commands and the raw-input fingerprint boundary. Checkboxes remain open until their actual merge/publication/acceptance gates pass; semantic normalization and execution-plan digests remain future producer contracts.

- [x] Create the owning runtime issue and acceptance scenarios in `storm-harbor-control-plane` before implementation; import the published producer version, not another repository's source or `@main`.
- [x] Define and document the actual CLI interface, stable exit codes and safe structured error/output contracts. Initial commands are for version/configuration inspection and offline validation; command names are finalized in that delivery.
- [x] Bound file size/reads, accept an explicit expected scope, reject ambiguous inputs and retain the strict parser's safe failures.
- [x] Distinguish a locally declared expected scope from authoritative online authorization. Passing local scope flags can never grant cloud or product access.
- [x] Keep offline validation free of cloud, vault and GitHub mutations; avoid raw document dumps in default output.
- [ ] Define versioned normalization and configuration digests without reading secret values; freeze schema/contract versions into future planning inputs.
- [ ] Produce checksummed executable artifacts with provenance and a declared OS/architecture support matrix. Test each claimed platform; cross-compilation alone is not operational acceptance.
- [x] Verify invalid configuration, cross-scope references, unsafe lifecycle, secret leakage canaries, cancellation/file errors and deterministic exit behavior.
- [ ] Publish the reviewed CLI release and its real usage examples; update Getting Started alongside the executable.

**Exit:** an adopter can run documented validation against synthetic exercise/restore YAML with no provider credentials. `plan`, `run`, `resume`, `cancel` and cleanup operations remain separate later capabilities until their contracts and execution gates are implemented.

**Source evidence:** manually merged [PR #15](https://github.com/Storm-Harbor/storm-harbor-control-plane/pull/15), main commit `9ad80e4274281200e95ed0158e48469e6e6b5b87`, passing [native Go acceptance](https://github.com/Storm-Harbor/storm-harbor-control-plane/actions/runs/43413904985), published [CLI guide](getting-started/offline-cli.md) and [documentation deployment](https://github.com/Storm-Harbor/storm-harbor-control-plane/actions/runs/43413904510). These checked source outcomes do not establish an immutable binary release. The [release runbook](operations/releases.md) describes the remaining signing/publication gates; semantic normalization remains a separate future producer contract.

### Stage 03 — Customer-owned private recovery repository

**Outcome:** the initial adopter has a private source of recovery intent that consumes Storm Harbor without exposing projects or forking its engine.

- [ ] After Stage 02, create the approved private repository in the customer's organization, using a generic structure such as `<customer-org>/disaster-recovery`.
- [ ] Confirm repository owners, effective product scope, branch protection/review rules and the selected source-repository allowlist.
- [ ] Add synthetic starting configurations and customer-owned runbooks/policies. Real adopter names, project IDs and bindings remain private; public templates stay synthetic.
- [ ] Add a credential-free validation workflow pinned to the published CLI/artifact identity and the approved configuration schema.
- [ ] Reference existing infrastructure and application repositories without copying their private source or operational evidence into the public project.
- [ ] Define logical source/provider/placement/secret aliases and customer-controlled external bindings; exclude secrets, state files, backups and private knowledge from Git.
- [ ] Record environment/region targets and separate exercise/restore declarations; installation and cloud-execution capabilities remain disabled until later gates pass.
- [ ] Agree external evidence storage, emergency-access ownership and lifecycle/retention policy.
- [ ] Choose whether to enable a repository dashboard and define its audience, hosting, disclosure policy and retention. Private repository visibility does not establish private website access.

**Exit:** private declarations pass the offline workflow; repository access and ownership are verified. This stage precedes live exercises and does not itself authorize App installation privileges, HCP applies or cloud provisioning.

Suggested layout, not an implemented automatic discovery convention:

```text
disaster-recovery/
  recovery/configurations/    # versioned exercise/restore declarations
  recovery/policies/          # reviewed requirements; engine contract delivered later
  recovery/journeys/          # dependency order/runbooks; journey contract delivered later
  docs/ownership.md
  docs/emergency-access.md    # procedure and aliases, never emergency keys
  .github/CODEOWNERS
  .github/workflows/          # released validation first; execution added after its gates
```

### Stage 04 — GitHub App and trusted scope

**Outcome:** an organization can install the integration for selected repositories, with authoritative scope and safe event handling.

- [ ] Define App hosting/ownership, organization installation, selected-repository permissions, token lifetime, key storage/rotation and uninstall behavior.
- [ ] Keep the App signing key outside ordinary repository/job configuration; use minimum short-lived installation tokens.
- [ ] Authenticate webhook delivery, validate signature/timestamp/delivery identity as applicable, deduplicate and resist replay before scheduling work.
- [ ] Resolve effective organization/project scope through authenticated installation identity plus authoritative membership/policy; YAML is never the authority.
- [ ] Separate installation tokens, `GITHUB_TOKEN`, workflow OIDC identity, human approval and provider/vault authorization.
- [ ] Enforce repository access and source commit selection; reject removed repositories, revoked installations and conflicting scope.
- [ ] Test one actor across two organizations/projects, colliding aliases, concurrent context switches and revoked access with no foreign data or effects.
- [ ] Document actual installation/onboarding steps only when the App endpoint and operational ownership exist.

**Exit:** a selected private repository triggers only its authorized, auditable work; uninstall/revocation prevents new access. There is no mandatory corporate identity service.

### Stage 04A — Optional assisted analysis and manifest drafting

**Outcome:** after explicit organization/project opt-in, a bot analyzes selected repository revisions and proposes recovery manifests in the private recovery repository. It does not apply them or fetch recovery secrets.

- [ ] Publish the scoped analysis/proposal contract separately from `RecoveryConfiguration`; the current alpha does not accept an AI profile under `spec.ai`.
- [ ] Let each project opt in and choose an approved agent adapter, including planned Copilot, OpenAI Codex, Claude Code or customer-owned agent integrations. Agent selection is separate from model/provider/reasoning selection; none is mandatory or implemented merely by this plan.
- [ ] Freeze the approved agent identity/version, runner placement, repository allowlist, data destination, tool capabilities, limits and structured proposal contract. A selected agent cannot widen the bot's effective permissions or inherit unrelated local credentials/plugins.
- [ ] Verify equivalent scoped analysis/proposal acceptance per adapter and fail closed when its required capability, entitlement or policy is unavailable. A disabled/unavailable agent leaves manual onboarding available and never triggers an unapproved agent/provider fallback.
- [ ] Select repositories, branches/commits, file/path classes, byte/context limits and permitted questions through reviewed installation/project configuration.
- [ ] Read authorized Terraform code, application/deployment manifests, dependency declarations and existing runbooks. Exclude raw state, environment/credential files, private keys, secret material and unapproved knowledge; redact/classify before any model egress.
- [ ] Use explicit read-only tool allowlists. Repository text is untrusted input and cannot enable shell execution, widen repository access, fetch vault data or change the tool policy.
- [ ] Establish model destination, data retention/residency, egress authorization, authentication, entitlement and cost/time ceilings before analysis. Repository read permission does not itself authorize exporting source to a hosted model.
- [ ] Discover services, infrastructure ownership and dependencies with provenance tied to source revision/file, confidence and unresolved questions; never invent backup availability, RTO/RPO, destination ownership or production credentials.
- [ ] Produce candidate YAML and a companion findings/missing-input report. Incomplete required inputs stay explicitly non-executable proposals rather than passing validation with invented values.
- [ ] Validate complete candidates with the released Go schema/CLI. Keep structural validation distinct from actual provider capability, binding readiness and execution authorization.
- [ ] Publish only through a separate, constrained proposal capability into a branch/PR of the private recovery repository. The analysis agent cannot write source application/infrastructure repositories, merge directly to main or start recovery.
- [ ] Let the operator manually register/confirm vault bindings and required destination/objective choices outside Git. Secret values are entered into the approved vault path/UI, never committed to the proposed manifest or included in model context.
- [ ] Require human review of provenance, assumptions, missing fields, target, policy and secret references before promotion. Execution later requires the independent exact-plan gates of Stage 08.
- [ ] Test injected repository instructions, embedded secrets, cross-project reads, disabled AI, denied egress, unsupported models, incomplete manifests, malformed output and attempted unauthorized write/tool use.

**Exit:** a permitted analysis yields a reviewed, validated proposal with explicit gaps and zero cloud/vault mutations or credential disclosure. Ordinary manual YAML onboarding and deterministic recovery work when this feature is disabled.

#### Model and reasoning policy

The [agent adapter policy](architecture/ai-providers.md#customer-selected-analysis-agents) permits customer-selected Copilot, Codex, Claude Code and private agents. Choosing an agent does not imply a particular model, a hosted destination or unlimited autonomous tool use. Each adapter must disclose and enforce its actual capabilities before activation.

The owner preference is an OpenAI GPT model for initial deep repository analysis with **`xhigh`** reasoning, while keeping provider, model and reasoning separately configurable. This is an opt-in project default for that task, not a mandatory platform provider or permission to run inference.

- Model selection must resolve to an explicit approved identifier and supported capabilities for the chosen integration; freeze actual model/adapter/version and effective reasoning into proposal provenance.
- An unavailable model or unsupported `xhigh` setting is an actionable configuration failure. Do not silently downgrade effort, switch providers, expose private input to a fallback or select an unreviewed moving default.
- Compare representative manifest-generation accuracy, missing-input detection, security failures, latency and cost before making the preset operational. The deep-analysis preference does not force maximum effort for deterministic extraction or every later task.
- As a current direct-OpenAI evaluation candidate, `gpt-6.1-sol` supports `xhigh`; its tool calls use the Responses API. This does not establish availability under a particular customer's Copilot catalog, account or policy. Verify again at implementation/activation time. See the [official model reference](https://developers.openai.com/api/docs/models/gpt-6.1-sol) and [reasoning guide](https://developers.openai.com/api/docs/guides/reasoning).

#### Copilot, reusable workflows and authentication

Copilot is an optional AI integration, not the Storm Harbor GitHub App itself. Its SDK exposes model discovery/selection and reasoning controls for supported models; effective options must be checked at runtime against the allowed catalog and account policy. The product's planned Go adapter must keep its SDK/runtime version and tool capabilities explicit. See [Copilot SDK compatibility](https://docs.github.com/en/copilot/how-tos/copilot-sdk/troubleshooting/compatibility).

For a supported organization-billed Actions flow, current GitHub documentation describes an explicit Copilot request permission plus enabled organization policy. For other services it documents installation authentication with separate Copilot organization permission/policy. These are additional gates, not privileges implied by ordinary repository access. User-authenticated and bring-your-own-provider-key paths have their own account/billing requirements. See [server-to-server authentication](https://docs.github.com/en/copilot/how-tos/copilot-sdk/auth/server-to-server-tokens).

Reusable workflows must accept bounded typed inputs, pin the released analysis executable, validate caller/repository/ref and maintain the read-versus-proposal-write separation. No automatic organization-policy enablement, personal-credential fallback, stored cloud key or model call is performed by this documentation delivery.

### Stage 05 — Workflow execution and runner identity

**Outcome:** a reviewed event or schedule starts a bounded, pinned Go executor in the permitted execution environment.

- [ ] Map reviewed commits, manual requests, allowed events and schedules to the exact configuration revision and execution mode.
- [ ] Pin workflow actions/executors by immutable revision/digest; validate artifact provenance and minimum job permissions.
- [ ] Support declared GitHub-hosted and customer self-hosted execution targets, with real acceptance for each supported platform/capability.
- [ ] Isolate private runner groups/identities, reject untrusted fork/PR code from privileged runners and prevent credential/workspace reuse.
- [ ] Validate OIDC issuer, audience, expiry and permitted repository/workflow/environment claims against customer-controlled trust policy.
- [ ] Use a reviewed local identity/broker when direct provider federation is unsupported; do not promise universal native OIDC support.
- [ ] Bound job duration, concurrency, cancellation and retries; distinguish diagnostic IDs from authorization and durable idempotency.
- [ ] Emit sanitized structured execution/audit events from the first worker. A successful job start is not recovery success.

**Exit:** a synthetic job runs on an allowed runner and unauthorized identity/event/runner combinations fail before credential access or mutation. Bounded Actions jobs do not constitute continuous reconciliation.

### Stage 06 — External vaults and bindings

**Outcome:** an authorized step resolves only its approved binding and credential while respecting customer location and disclosure rules.

- [ ] Version the provider/placement/secret binding and resolver ports without adding raw secrets or sensitive vault paths to the core declaration.
- [ ] Scope every lookup by authenticated organization/project, purpose, resource, operation and approved plan/version.
- [ ] Exchange temporary workload identity for minimum provider/vault capability, with explicit expiration, audience and revocation semantics.
- [ ] Prefer dynamic database credentials where the actual integration supports them; separately rotate/version any unavoidable static password.
- [ ] Keep sensitive execution on an isolated customer runner when a credential or data cannot leave that environment.
- [ ] Prevent disclosure through logs, exceptions, command arguments, subprocess environments, caches, files, artifacts and AI input. Masking alone is insufficient.
- [ ] Handle expired/revoked/wrong-scope identities, denied secrets, partial resolution and vault unavailability without an unsafe fallback.
- [ ] Ensure vault/KMS/broker access and emergency identity survive the modeled failure domain; never make recovery depend exclusively on the system being rebuilt.

**Exit:** real sandbox resolver acceptance proves narrow access, isolation, expiration and cleanup under failure/cancellation. A temporary vault-access token does not make the stored secret itself temporary.

### Stage 07 — Read-only inventory and readiness

**Outcome:** the system can explain what is protected, where its evidence comes from and what blocks recovery.

- [ ] Version source/provider discovery ports and normalized service/environment/resource/dependency/recovery-point projections.
- [ ] Read only explicitly authorized GitHub, infrastructure-state, cloud, cluster and backup sources; keep sensitive state outside public projections and prompts.
- [ ] Preserve stable identity, provenance, observation time, confidence, freshness and explicit ownership for each relationship.
- [ ] Distinguish desired state, managed state, actual state, backup integrity and business recovery health; discovery never silently takes ownership.
- [ ] Define drift/staleness/unknown results, bounded pagination and dependency traversal with effective scope enforced before data crosses an adapter.
- [ ] Inventory backup age, integrity, decryptability, location and restore compatibility; backup existence alone is not recoverability.
- [ ] Add bounded timeouts, cancellation, retry classification/backoff and dependency failure handling; deterministic denial is not retried.
- [ ] Test partial/unavailable sources, guessed/colliding IDs, scoped counts/bulk/traversal, stale evidence and concurrent access.

**Exit:** a private adopter can inspect a safe, sourced readiness assessment without granting mutation permissions. Missing or stale mandatory evidence produces an explicit blocker/unknown state.

### Stage 08 — Policy, planning and durable orchestration

**Outcome:** every future mutation belongs to an independently authorized immutable plan and a resumable run.

- [ ] Publish journey/DAG, policy, approval, run-state, checkpoint and evidence contracts before their consumers.
- [ ] Evaluate the intersection of installation, project, environment and run policies; repository configuration cannot relax an organizational ceiling.
- [ ] Resolve actual source/target identities, prevent source mutation through alternate aliases and reject unsupported data/network/service mappings before approval.
- [ ] Freeze configuration/schema, source commits, adapters/controllers, recovery points, target identity, allowed actions and cost/time ceilings into the approved plan digest.
- [ ] Bind approval to the authenticated actor, exact plan/target and expiry; changed inputs or revoked permission invalidate it.
- [ ] Persist scoped idempotency, leases/concurrency, transitions and external operation IDs atomically with bounded retention and conflict behavior.
- [ ] Reconcile an unknown provider outcome before retrying mutation; a failed status write never justifies blindly repeating an accepted external action.
- [ ] Define health gates, cancellation, retry classes, compensation and independently authorized cleanup. Preserve retained targets and pre-existing resources.
- [ ] Store allowlisted evidence outside the recovery failure domain with integrity/provenance and reauthorized retrieval.
- [ ] Test replay, duplicates, concurrent jobs, partial persistence/provider failures, revoked approval, kill switch and safe resume with no unintended effects.

**Exit:** deterministic synthetic/provider-sandbox runs can plan, pause, resume, cancel and retain evidence without bypassing scope or repeating unsafe mutations. AI is not required for any approved run.

### Stage 08A — Repository dashboard and historical reports

**Outcome:** each adopter recovery repository can publish a central, read-only dashboard of its authorized history, without requiring the full operational console. See the [dashboard adoption and publication guide](getting-started/recovery-dashboard.md).

- [ ] Publish a versioned, scoped report projection and reporting/publication policy before implementing its Go generator, template or reusable workflow. The initial alpha accepts no dashboard configuration fields.
- [ ] Provide overview, chronological run history, per-run phase timeline, gap-analysis checklist and dependency/journey diagrams. Clearly separate exercise from retained restore, planned topology from observations, and measured objectives from unknown values.
- [ ] Preserve accepted snapshots for every recorded run, including failed, aborted, partial and inconclusive outcomes. Version subsequent corrections and identify missing evidence instead of removing failed attempts or reporting success.
- [ ] Link gaps to source/evidence, severity, owner, remediation and verification. A closed checklist item requires evidence; neither AI output nor a dashboard checkbox authorizes execution.
- [ ] Generate diagrams from scoped versioned inputs, with escaped labels, safe links and accessible text/table alternatives. Apply Storm Harbor colors/layout patterns without copying private console code, telemetry or customer data.
- [ ] Export only audience-approved summary fields. Keep raw evidence, logs, state, credentials, private knowledge and signed access URLs out of the static artifact, search index, downloadable JSON and diagram source.
- [ ] Verify the site's actual access controls, organization entitlement and complete audience before publication. GitHub Pages private publishing requires Enterprise Cloud; otherwise use approved private hosting or a protected local report. Public publishing is an explicit synthetic/disclosure-approved option, never an automatic fallback.
- [ ] Preserve authoritative external evidence with integrity and retention outside the modeled failure domain. Rebuild the dashboard from that evidence; Pages and Actions artifacts are not the only historical record.
- [ ] Use a pinned credential-free rendering step over a reviewed export, separate evidence-reader and publisher permissions, immutable assets and serialized publication that cannot lose concurrent runs. Deny untrusted PR deployment and prevent remote diagram scripts or repository content from widening tools/egress.
- [ ] Display snapshot time, freshness, provenance and publication health. Failed evidence export or website deployment does not change the recovery result; a stale dashboard never claims live readiness.
- [ ] Test scope/audience denial, full-artifact secret canaries, injected diagram content, history retention, concurrency, missing evidence, failed publication, mobile/keyboard access and a successful rebuild from external evidence.

**Exit:** an authorized viewer can follow a synthetic run and its gaps from the repository overview to its historical summary and diagram, while unauthorized viewers cannot retrieve any private artifact. Stage 03 can prepare a scaffold; historical execution reporting waits for the authoritative Stage 07–08 contracts. The dashboard is optional and cannot approve, start or resume recovery.

### Stage 09 — GCP target foundation

**Outcome:** an approved plan creates a new isolated GCP target or verifies an existing one through the correct infrastructure owner.

- [ ] Map logical provider/placement aliases to permitted organization/folder, project, billing, region/zone, network and data-residency bindings.
- [ ] Treat a GCP project as a global container; resource regions/zones are separate placement choices.
- [ ] Separate project creation, billing association, API enablement, service IAM, quota/capacity and readiness checks. Enabled APIs do not imply permissions.
- [ ] Verify the actual credential consumer/quota project and required long-running-operation read/list/cancel/delete capabilities where applicable to each service.
- [ ] Establish narrow bootstrap trust before target creation without workload self-elevation; freeze provider/version inputs.
- [ ] Where Terraform/HCP owns resources, use its reviewed VCS/plan/apply path and a dedicated recovery root/workspace/state binding. Do not redirect production state or imperatively modify those same managed resources.
- [ ] For a different infrastructure executor, declare one explicit owner per resource. Terraform/HCP is an integration choice, not a mandatory product dependency.
- [ ] Compare actual target identity with protected sources, verify compatibility, enforce hard action/cost ceilings and record durable run ownership.
- [ ] Poll/reconcile creation safely, handle partial failure and verify readiness before restore. Do not treat budget alerts as a spending limit.
- [ ] Test create/existing with both exercise and retained intent in safe, explicitly approved sandboxes.

**Exit:** the adapter proves target identity, permissions, readiness and infrastructure ownership. No exercise is enabled until the required run/evidence/cleanup controls also pass.

### Stage 10 — Isolated recovery exercise

**Outcome:** the initial adopter can execute and verify a recovery rehearsal without modifying protected production resources.

- [ ] Select `exercise`, a new target or existing isolated sandbox, positive expiry and cleanup restricted to proven run-created resources.
- [ ] Use synthetic/sanitized data by default; production backup use requires explicit access, isolation, encryption and residency approval.
- [ ] Verify the recovery point, rebuild infrastructure through its owner and deploy reviewed immutable application images/revisions through application-owned workflows.
- [ ] Restore data into the isolated target and test integrity, schema/connectivity, dependencies and business-level behavior. Liveness alone is not success.
- [ ] Define the measured incident/start and business-recovery completion points for RTO, and the verified recovery-point age for RPO; report unknown/missing measurements explicitly.
- [ ] Enforce time, quota, cost and action ceilings, exact-plan approval and kill switch before/between unsafe actions.
- [ ] Reauthorize cleanup and verify real run ownership; preserve existing sandboxes, shared assets and any resource not independently proven run-owned.
- [ ] Exercise cancellation, expiry, denied/failed cleanup, partial restore, provider/network failure and orphan detection. Provide a separate authorized janitor path.
- [ ] Retain sanitized signed/integrity-checked evidence outside the failure domain, including failures and cleanup outcomes.
- [ ] When the repository dashboard is enabled, publish the approved exercise summary, gaps, measured outcomes and diagram; record publication failure separately from recovery/cleanup outcomes.
- [ ] Complete at least three consecutive isolated exercises meeting their declared objectives before leaving mock mode for that supported scenario.

**Exit:** functional recovery objectives and cleanup are measured, not inferred from a green workflow. An exercise never redirects production traffic automatically.

### Stage 11 — Retained restore and emergency operation

**Outcome:** a separately approved incident recovery can restore into a retained target and, when authorized, move service traffic.

- [ ] Require `restore`, explicit target/recovery-point selection, `retained`/`retain` lifecycle and no automatic expiry cleanup.
- [ ] Establish independent incident permissions and approval for the exact plan, target, data, identities and destination network.
- [ ] Validate integrity and business readiness before a separately approved cutover; application owners retain release, migration and traffic authority.
- [ ] Define rollback/compensation and forensic retention without automatically deleting partially restored or pre-existing resources.
- [ ] Keep approved signed plans, artifacts, backups, runbooks and independent evidence available when GitHub, the primary network or original control plane is unavailable.
- [ ] Deliver a verified local/self-hosted emergency executor with independent time-bounded identity, current policy checks and durable audit; unavailability is not a permission bypass.
- [ ] Model regional, account/project, identity and provider-wide outages explicitly. Another project/region does not by itself solve a provider-wide failure.
- [ ] Test revocation, stale plans, inaccessible KMS/vault, unavailable source APIs, interrupted cutover and external-operation reconciliation in approved exercises.
- [ ] When reporting is enabled, retain the approved restore/cutover summary and phase history under the repository's disclosure policy, without exposing incident-sensitive details.

**Exit:** authorized retained recovery and cutover have scenario-specific evidence and runbooks. No real incident automation is inferred from Stage 10 alone.

### Stage 12 — Resource/controller extensions and Kubernetes

**Outcome:** organizations can extend recovery safely without changing the deterministic core or requiring Kubernetes everywhere.

- [ ] Publish versioned resource definitions, schema/admission, status and controller lifecycle contracts with namespaced kind ownership.
- [ ] Define conversion/backward compatibility, migration fixtures and frozen controller/resource versions for approved runs.
- [ ] Isolate controllers and restrict their capabilities, credentials, egress and lifecycle; custom YAML cannot select an unrestricted executable or grant access.
- [ ] Implement continuous reconciliation as an explicit optional mode, with leases, cancellation, bounded retries, safe shutdown and independently measured recovery health.
- [ ] Deliver native Kubernetes CRDs and a Go operator/controller through the Kubernetes adapter, with correct admission, RBAC and observed status.
- [ ] Integrate existing GitOps/Flux ownership rather than racing another reconciler for the same resources.
- [ ] Test admission denial, version conversion, duplicate reconciles, drift, controller crash/restart, partial writes, revoked scope and source protection.

**Exit:** a custom kind/controller and a Kubernetes recovery scenario meet their declared lifecycle/security acceptance. A Kubernetes-style resource envelope alone is not an installed CRD.

### Stage 13 — Additional clouds and multicloud mappings

**Outcome:** supported providers and concrete source/target pairs have explicit capability and acceptance evidence.

- [ ] Define a versioned capability registry for discovery, export/import, provision, restore, validation, cleanup and identity integration.
- [ ] Add AWS account/region, Azure tenant/subscription, DigitalOcean team/project and OCI tenancy/compartment mappings behind adapters; keep those shapes out of the core.
- [ ] Record discovery/simulation/exercise/retained-restore support separately for each provider and service, including unsupported combinations.
- [ ] Verify data format, encryption/KMS, IAM, networking/DNS, service semantics, regions, egress costs and recovery-point portability for each cross-cloud mapping.
- [ ] Plan provider-wide outages with independent executor, identity, backups, decryption and evidence paths.
- [ ] Reuse scope, idempotency, approval, ownership and cancellation contracts; do not infer portability from a generic resource name.
- [ ] Run provider sandbox and end-to-end acceptance per supported pair with sanitized fixtures and explicit resource/cost approvals.

**Exit:** users can choose only mappings whose declared capabilities and scenario evidence support their requested mode.

### Stage 14 — Optional private AI/RAG/CLI

**Outcome:** a project can use its own intelligence service without disclosing its corpus or expanding execution authority.

- [ ] Publish a scoped integration contract for endpoint, identity, data classes, input/output schemas, allowed tasks/tools, time/byte/cost limits and retention/egress.
- [ ] Keep retrieval within the project's permitted environment and apply authorization filters before retrieval; return only approved minimal results.
- [ ] Prevent cross-project cache/vector/file/job leakage, implicit shared memory and automatic provider fallback with private content.
- [ ] Treat retrieved content and model output as untrusted proposals; reject prompt-injected tool/permission changes and invalid references.
- [ ] Execute CLI integrations in bounded workers with fixed executable/argument schemas, not an unrestricted shell.
- [ ] Retain safe provenance and structured decisions without raw secrets, prompts, corpora or unrestricted internal reasoning.
- [ ] Test colliding identifiers, concurrent scope switches, revoked membership, exfiltration canaries, injected instructions, output validation and deterministic operation with AI disabled/unavailable.

**Exit:** the project-owned plugin can explain/propose permitted work while deterministic policy, approval and orchestration remain authoritative. This stage does not block ordinary recovery.

### Stage 15 — Go APIs and operational console

**Outcome:** a user can manage authorized recovery through an API/console using the same contracts as CLI/workflows.

The Go APIs belong to the primary control-plane repository. The separate monitoring-center application is **Harbor Lens**, in `harbor-lens`; the website stays in `storm-harbor-app`. See the [Harbor Lens plan](SPECS-HARBOR-LENS.md) for sandbox, private observation, Harbor groups and approved operation. The first synthetic frontend may proceed without an unpublished API; live integration retains all producer/security gates.

- [ ] Publish transport-neutral domain contracts and explicit versioned OpenAPI/safe error mappings before API consumers.
- [ ] Implement configurable human identity and authoritative organization/project access, including revocation and protected route/resource policies.
- [ ] Preserve secure session/MFA behavior while migrating transitional authentication/BFF runtime boundaries to Go; do not import a mandatory private identity service.
- [ ] Implement scoped inventory, configuration, plans, approvals, runs, evidence and admin operations, including files/exports/counts/bulk/traversal.
- [ ] Make the guided UI edit/explain the same versioned YAML instead of maintaining a second workflow definition.
- [ ] Maintain all configured UI locale catalogs, RTL/accessibility and clear delivered/planned capability labels; public Pages output continues to exclude authenticated/API routes.
- [ ] Validate real user journeys against authoritative backend outcomes using isolated runtime dependencies, not dashboard mocks as proof of recovery.
- [ ] Measure required coverage dimensions independently under the owning guidance; keep unmeasured Go branch/line or legacy application gaps explicit rather than claiming compliance from native statement coverage.

**Exit:** an organization/project-isolated user journey reaches real backend outcomes with appropriate contract, integration, browser and security evidence. Static website publication does not qualify the authenticated application.

### Stage 16 — Operational qualification and supported releases

These controls are developed from Stage 02 onward. This stage is the final qualification of the supported capability, not a deferred introduction of security or observability.

- [ ] Operate the control plane, identity, backups/decryption and evidence outside the modeled failure domains; verify degraded/emergency paths.
- [ ] Distinguish startup, liveness, dependency readiness, worker readiness and business recovery health. Bound shutdown and leave unfinished work resumable.
- [ ] Emit structured sanitized production events and separate security audit decisions; control verbosity and exclude unrestricted request/provider bodies.
- [ ] Define bounded-label Prometheus/OTLP signals for queues/controllers, dependencies, denied operations, recovery outcomes and cleanup; protect metrics access and avoid customer identifiers/high-cardinality secrets.
- [ ] Publish checksums/provenance/SBOM as applicable, pin action/runtime/dependency inputs, verify vulnerabilities and prevent release-tag/artifact replacement.
- [ ] Validate rate/cost/resource limits, state-store durability, backups, replay/concurrency, upgrades, rollback, load and dependency outage behavior.
- [ ] Maintain operator runbooks, support boundaries, incident response, vulnerability disclosure and compatibility/release policy.
- [ ] Qualify only supported scenarios/platforms with executed acceptance and honest coverage measurements; report missing/skipped scenarios separately.
- [ ] Obtain the required owner/customer approvals before enabling scheduled or real recovery operations; no release label enables mutation implicitly.

**Exit:** a supported release has a bounded support matrix, reproducible evidence and exercised operational ownership. Alpha availability is not production qualification.

### Stage 17 — Community, funding and CNCF readiness

This is an independent track. Sponsorship never grants product authorization, private data access or implicit support/SLA entitlements.

- [x] Publish contribution, conduct, security, maintainer and informational funding guidance.
- [ ] Complete owner-controlled receiving-account/onboarding verification and agree contribution currency, recurrence, amounts and terms.
- [ ] Validate hosted payment flows in test mode without real charges; publish only verified live destinations when authorized.
- [ ] If future product entitlements require payment events, deliver signed/replay-bounded durable webhook acceptance and idempotent reconciliation separately.
- [ ] Collect real implementation, governance, licensing, adopter and recovery evidence for CNCF readiness.
- [ ] Check current foundation eligibility/application and legal requirements before submission; affiliation is claimed only after actual acceptance.
- [ ] Track funding activation separately in [community #3](https://github.com/Storm-Harbor/.github/issues/3) and the [readiness guide](community/cncf-readiness.md).

## 5. Cross-cutting gates for every executable delivery

Before enabling an operation, its owning issue must prove the applicable items below. They apply to retries, background work, callbacks, caches and extensions as well as initial requests.

- [ ] Effective scope comes from authenticated identity and authoritative membership/policy; it is enforced before persistence, retrieval or provider calls.
- [ ] Reads/writes/counts/existence checks/bulk/traversal, jobs, caches, idempotency, files and evidence preserve the same organization/project boundary.
- [ ] Revocation, wrong scope, colliding IDs, concurrent context changes and explicit administrative operations have negative acceptance scenarios with no foreign data/effects.
- [ ] Secrets/private knowledge remain outside default output and public artifacts; resource references and metadata follow customer disclosure policy.
- [ ] Timeouts, cancellation, bounded retries/backoff and unknown-outcome reconciliation are owned at adapter boundaries; deterministic denial is not retried.
- [ ] The exact mutation plan, target, adapter/controller versions and recovery point are approved; altered inputs invalidate approval.
- [ ] Source protection and deletion ownership use actual identities and authoritative records, not aliases or labels alone.
- [ ] Cleanup/expiry is independently authorized, tested under partial failure and disabled for retained targets unless a separate explicit action is approved.
- [ ] The kill switch, isolated runner placement, external vault/KMS and independent emergency/evidence paths meet the scenario's failure model.
- [ ] Evidence proves business/data recovery, RTO/RPO and cleanup outcomes; process uptime, a valid YAML or a green job alone is insufficient.
- [ ] Tests use synthetic data and isolated Testcontainers for applicable owned dependencies; provider acceptance uses explicitly approved sandboxes, never implicit production execution.
- [ ] Required coverage and executed/required/passed scenarios are recorded without treating unmeasured dimensions or mocked dependencies as integrated proof.

## 6. Acceptance scenarios to carry into owning issues

| Area | Required examples |
| --- | --- |
| Configuration | Valid exercise/restore; create/existing; malformed/unknown/duplicate fields; oversized/deep documents; type/alias attacks; safe errors |
| Scope | One actor in two organizations/projects; guessed/colliding IDs; switched/concurrent context; revoked membership; scoped traversal/counts/bulk/jobs |
| Identity/vault | Wrong issuer/audience/workflow/repository; expired/revoked credentials; denied alias/version; unavailable vault/KMS; no persistent secret material |
| Planning | Unsupported provider mapping; same real source/target under different aliases; stale recovery point; changed plan; expired/revoked approval |
| Execution | Duplicate/replayed event; concurrent lease; cancellation/kill switch; timeout; provider outcome unknown; status-store failure; restart/resume |
| Target/cleanup | New/existing GCP target; wrong parent/billing/quota project; denied API/IAM; partial creation; failed cleanup; pre-existing/shared/retained resource protection |
| Recovery | Backup integrity/decryptability; application revision and dependency health; business validation; measured RTO/RPO; cutover/rollback gates |
| Resilience | GitHub/source/primary control-plane/region/provider unavailable; independent identity, backups and evidence; no permission bypass |
| AI/extensions | Cross-project retrieval denial; prompt injection; bounded CLI/tool schema; no private fallback egress; operation with AI disabled; controller conversion/crash |
| Repository dashboard | Actual host access control; no private artifacts via direct URLs/search/JSON/diagrams; all terminal run outcomes; scoped gap evidence; concurrent updates; stale/failed publication; rebuild from independent evidence |
| Console/operations | Real authorized API outcomes, accessibility/RTL, safe exports/evidence, protected metrics, graceful shutdown, upgrade/rollback and supported-platform evidence |

## 7. Delivery and progress rules

1. Before runtime implementation, create the owning native Journey/Epic and repository delivery issues with concrete scope, dependencies and focused-hour estimates. Future stages without owning issues are planned work, not active implementations.
2. Every delivery has assignee, native Feature/Bug type, labels, Priority, Size, Estimate, Workspace, Product and Status in the maintainer project. Public contributors follow public issue/release links; maintainer project access is separate.
3. Use the issue branch, English Conventional Commits, synthetic fixtures and one versioned owner for shared semantics. Complete producer checks, manual merge and immutable publication before consumer implementation.
4. Record acceptance scenarios and required metrics before coding. Reuse equivalent passing evidence for the same candidate; repeat checks when relevant source/dependency/toolchain changes, failures or unresolved concerns require it.
5. Prepare a concrete PR with exact validation and deployment/compatibility boundaries. The user reviews and merges manually; no automatic merge is enabled.
6. After merge, verify the delivered revision, deployment/release and public or private acceptance appropriate to that outcome. Reconcile issue closure and project fields only after all gates are satisfied.
7. Keep one current owning-issue checkpoint with timestamp, revision, accepted decisions, evidence, actual blockers and next bounded action. Complete Journey/Epic records only after their descendants and acceptance gates are complete.
8. Update this ledger, Getting Started and executable examples with each accepted capability. Do not mark planned commands, policies, controllers or cloud adapters as available because their design text exists.

### Immediate sequence

The initial producer is released. The next runtime delivery is Stage 02, followed by Stage 03 private adopter configuration. App/identity/runners/vaults and read-only discovery then unlock an approved GCP target and isolated exercise. Retained restore, custom controllers, additional clouds, AI and console capabilities follow their own acceptance gates and can proceed independently where their dependencies permit.

The public plan deliberately uses generic adopter placeholders. The owner-approved private repository name, internal workspace/project bindings and customer-specific adoption decisions are maintained in private notes and later in the customer's private repository.
