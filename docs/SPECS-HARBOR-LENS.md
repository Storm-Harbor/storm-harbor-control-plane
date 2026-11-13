# Harbor Lens — Plan and acceptance checklist

Harbor Lens is Storm Harbor's own operational frontend and public reference application. The owner accepted this direction on 2026-11-13, replacing the earlier Harbor Commons proposal. OpenLens and Terraform dashboards inspire the experience; they are not adopted runtimes, forks or authentication systems. Gitea is not the selected application.

This is accepted architecture, not a claim of implemented production capabilities. GitHub owns live status: [Journey #16](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/16) → [Epic #17](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/17) → [specification #18](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/18) and [sandbox #1](https://github.com/Storm-Harbor/harbor-lens/issues/1).

## Ownership

| Repository | Responsibility |
| --- | --- |
| `storm-harbor-control-plane` | Go APIs, CLI, GitHub App/bots, orchestration, collectors, adapters, controllers and contract ownership |
| `harbor-lens` | Separate browser application, presentation and sandbox; future consumer of published core APIs |
| `storm-harbor-app` | Existing multilingual product website at `stormharbor.io`; stays in its repository |
| Customer recovery repository | Private YAML intent and optional historical report site; no fork of the engine |
| Existing contract module | Immutable `RecoveryConfiguration v1alpha1` compatibility dependency |

Backend components target Go 1.26.x. Browser HTML/CSS/JavaScript and development/build tools are separate. No new core repository is needed for a collector or Lens API. The first UI slice proceeds independently because it uses synthetic local fixtures and consumes no unpublished API.

## Frontend technology and visual decisions

The owner refined the frontend decision on 2026-11-13: use the same technologies as the reference console, with Storm Harbor's own identity. The selected foundation is React 19, Vite 8, JavaScript/JSX, React Router, Tailwind CSS 4 and shadcn/ui `radix-nova` with Radix primitives, Lucide icons, i18next/react-i18next and Recharts. Use Vitest, Testing Library and Playwright for complementary validation. Pin resolved dependencies, Node 26.10.0 and Yarn 4.18.0 strict Plug'n'Play with a repository-local cache and lockfile. Node is a build/test tool, not the production backend.

Do not adopt the reference console's private packages, corporate branding, business rules, accounts, authentication service or credentials. The user selected its technologies only. Introduce state libraries only for an actual application need; do not carry legacy dependencies into the sandbox.

Recharts owns metrics/time-series presentation. Advanced dependency/journey graphs are planned through the open-source JointJS Community/Core edition, consistent with the documentation's existing diagram direction. [JointJS licensing](https://www.jointjs.com/license) distinguishes free MPL-2.0 Community from commercial JointJS+; no paid extension, proprietary React integration or purchase is approved/required by this delivery. Verify each capability and preserve license/notices before adoption. The first sandbox topology is an accessible synthetic projection, not an implemented advanced graph editor. Future pan/zoom, bounded expansion, filters and node inspection require explicit scenarios; each real edge needs source/freshness provenance and scoped access.

Colors follow the existing product website's canonical tokens: Storm blue `#326ce5`, dark blue `#2554b8`, light blue `#84a7ef`, ink `#111820`, panel `#151e28` and paper `#f4f6f8`. Preserve accessible contrast and semantic state colors. The product website stays in its existing repository; this decision changes Lens, not that site's stack or publication.

## Company SSO, social sign-in and mandatory MFA

The owner explicitly requested company-owned SSO, GitHub and Google sign-in, and mandatory MFA on 2026-11-13. These are planned protected-account capabilities, not implemented by the synthetic sandbox.

| Entry | Planned adapter |
| --- | --- |
| Company-selected identity provider | OIDC first; SAML 2.0 through a separately qualified adapter |
| Google | Validated OIDC authorization-code flow |
| GitHub | Supported OAuth 2.0 web flow with minimal identity scopes; distinct from bot installation/repository grants |

Core Go handles callbacks, token validation/exchange, policy, session and factor lifecycle. The browser receives a bounded same-origin session using `Secure`, `HttpOnly`, appropriately scoped `SameSite` cookies; provider tokens, secrets and pending MFA credentials never enter local/session storage, URLs, logs or the frontend bundle. Validate issuer, audience, signature, expiry, state, nonce/PKCE as applicable, exact redirect destinations and configured provider metadata. Follow [OIDC Core](https://openid.net/specs/openid-connect-core-1_0.html), [OAuth security BCP](https://www.rfc-editor.org/rfc/rfc9700.html), [Google OIDC](https://developers.google.com/identity/openid-connect/openid-connect), [GitHub authorization](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps) and [SAML browser SSO](https://docs.oasis-open.org/security/saml/v2.0/saml-profiles-2.0-os.pdf).

Authentication does not grant organization/project/Harbor membership, provider IAM or repository access. Bind identities to stable provider subject/issuer identifiers, not email alone; prevent automatic account linking or tenant enrollment based only on matching domains/emails. Company federation configuration is an authorized administrative operation with reviewed issuer/metadata/destination and claim mapping.

Require verified MFA before every protected application session. An upstream login, configured provider MFA setting or unverified claim is not proof that the current authentication met policy. Accept upstream assurance only when validated and explicitly trusted by company policy; otherwise require a Lens-managed second factor. Plan WebAuthn/passkeys with user verification and explicit assurance rules, TOTP and hashed single-use recovery codes. Protect enrollment, replacement, reset and recovery with expiry, rate limits, audit and independent proof. No SMS/email-only fallback or silent MFA bypass is introduced. Pending enrollment/challenge sessions cannot read customer data. Reauthenticate for sensitive operations and invalidate sessions after revocation, factor reset or policy changes.

Automated collectors use separately scoped short-lived workload identities, not reusable human sessions or exemptions. Emergency identity/recovery procedures must survive provider outage without disabling the MFA requirement. The unauthenticated synthetic sandbox contains no customer connection or protected operational capability; a sandbox switch cannot create a real authenticated session.

## Harbors and federation

A **Harbor** is a registered control-plane instance or managed environment boundary, bound to authenticated organization/project scope and explicit capabilities. It is not necessarily a Kubernetes cluster, cloud account or public endpoint. Final resource/identifier contracts remain to be designed.

A **Harbor group** is a logical federation of explicitly authorized Harbors. Membership does not merge permissions or grant cross-organization access. Reads, counts, traversal, cache, evidence and actions enforce each member's effective scope before data crosses a boundary.

```joint
flowchart LR
  Lens[Harbor Lens] --> API[Scoped Go API]
  API --> A[Authorized Harbor A]
  API --> B[Authorized Harbor B]
  A --> Cloud[Cloud and Kubernetes adapters]
  B --> IaC[Infrastructure and metrics adapters]
```

This is the planned live flow. The initial sandbox has no API, identity service or collector.

## Three distinct modes

| Mode | Data and effects | Gate |
| --- | --- | --- |
| Sandbox | Synthetic inventory, metrics, diagrams and simulated runs; no external effects | Visible synthetic identity and isolated state |
| Observation | Customer-selected sources, bounded read-only inventory/metrics | Published contracts, authenticated scope, provider read grants and collection/egress/retention policy |
| Operation | Exact approved plans executed through core orchestration | Independent policy, frozen plan, approval and scoped workload identity |

Inspired by [Stripe sandboxes](https://docs.stripe.com/sandboxes), isolation covers identity, storage, caches, jobs, evidence and external destinations. A mode switch cannot grant authority. Synthetic outcomes are not readiness evidence. A real cloud drill is a separately authorized, budgeted operation with target lifecycle and cleanup controls. This analogy introduces no payment dependency.

## Views and private data

The intended views cover authorized Harbor inventory, cloud/Kubernetes/databases, desired/managed/observed state, Terraform ownership/drift, freshness and unknown states, accessible topology, YAML journeys, target environment/region, exercise versus restore, create versus reuse, phase outcomes, cancellation, gaps and history. Actual measured RTO/RPO must remain distinct from targets and simulations.

Optional opt-in AI agents keep the same permitted repositories, data policy and proposal review gates. Agent selection never grants authority or introduces a shared knowledge corpus.

Customers explicitly select repositories, sources and capabilities. GitHub App access does not grant provider IAM or permission to export code to an AI provider. Planned Go collectors may run inside customer networks, including outbound-only/self-hosted execution. Secrets stay in customer-selected vaults; prefer short-lived workload identity. Static credentials require separately managed rotation/lifecycle.

Never put raw Terraform state, kubeconfig, tokens, credentials, raw logs, prompts or knowledge bases into browser fixtures, public Pages, logs or frontend storage. Connections require destination/protocol allowlists, SSRF/DNS/private-network and egress policy, bounded polling, timeout/cancellation, retry, cardinality, bytes, cost and retention. Sensitive/high-cardinality customer identifiers cannot become metric labels. Reauthorize downloads, caches, jobs and evidence after membership/policy changes.

The separate [repository dashboard](getting-started/recovery-dashboard.md), Stage 08A, remains a historical evidence projection with independently enforced hosting audience. A private repository alone does not establish private Pages hosting.

## Contracts before live consumers

Publish reviewed contracts for Harbor registration/capabilities/revocation; organization/project and Harbor group membership; connections/workload identity/data policy; mode-bound identities/stores/destinations; allowlisted inventory/health/freshness/topology/metrics; safe errors/pagination/cancellation; and exact-plan operations/idempotency/history/evidence.

Core owns these semantics. Producer validation, manual merge, immutable publication and acceptance precede live consumers. Current `RecoveryConfiguration recovery.stormharbor.io/v1alpha1` accepts none of `spec.lens`, `spec.harbors`, `spec.mode`, `spec.metrics` or `spec.ai`. This plan is not executable YAML; do not add speculative fields to the alpha validator.

## Delivery checklist

### L0 — Recorded scope

- [x] Accept our own separate Harbor Lens application; preserve website ownership.
- [x] Register native Journey/Epic and owning issues in Project #1.
- [x] Create `harbor-lens` without live integrations.
- [ ] Manually merge/publish this specification and verify its public URL.

### L1 — Synthetic sandbox

Revised estimate after frontend, identity and hosting refinements: specification 8 focused hours plus frontend 24 focused hours; not elapsed work or production qualification.

- [ ] Build original responsive overview, inventory, metrics, topology, recovery and history views.
- [ ] Keep a persistent sandbox/synthetic label; live connections, observation and operation stay unavailable.
- [ ] Provide an explicitly fake sign-in journey with GitHub, Google or company SSO selection and a simulated mandatory MFA step; collect no passwords, real codes, tokens or IdP data. Session-only demo access grants no operational authority. Keep real authentication in L2.
- [ ] Provide synthetic Harbor creation and simulation without credentials, network APIs or browser persistence.
- [ ] Label client-side scope filtering as demonstration, never backend authorization.
- [ ] Maintain locale parity, keyboard navigation, mobile menu, bounded safe labels and accessible graph alternatives.
- [ ] Exercise the actual built artifact: negative paths, isolation, run history, cancellation/reset, i18n and responsive behavior.
- [ ] Measure statements, branches, functions and lines independently; target 100% each and record gaps honestly.
- [ ] Prepare the source PR and publish the public synthetic demo preview at `lens.stormharbor.io` through Cloud Run, following the owner's final hosting decision on 2026-11-13. Use a pinned non-root static Nginx container, immutable image digest and independently qualified Terraform/WIF/edge producers; Node stays build-only. Keep `stormharbor.io` in its existing repository. Manual merge, infrastructure applies, release/deployment checks and a live URL check precede publication claims. GitHub Pages is not the Lens application host. Live integrations remain separate deliveries.

### L2 — Authoritative contracts and identity

- [ ] Publish reviewed registration, membership, mode, projection and operation contracts.
- [ ] Enforce scope before storage, retrieval, counts, graphs, exports and cache access.
- [ ] Test colliding IDs, switched/concurrent scope, revocation, jobs, idempotency and group membership.
- [ ] Separate sandbox/observation/operation identity, storage, destinations and evidence.
- [ ] Publish configurable company SSO, Google/GitHub sign-in, safe server session and mandatory MFA contracts before login consumers.
- [ ] Qualify OIDC/GitHub adapters and later SAML; test callback mix-up, forged/expired tokens, issuer/audience/redirect denial and account-linking attempts.
- [ ] Test mandatory MFA enrollment/challenge/recovery, replay/rate limits, verified upstream assurance, factor revocation and session invalidation; no protected access before assurance.

### L3 — Read-only observation

- [ ] Implement scoped Go connection/collector ports and initial adapters.
- [ ] Add customer-controlled workload identity, endpoint/egress policy and collection budgets.
- [ ] Handle denied, stale and partial sources, timeout/cancellation and bounded retries without inventing healthy metrics.
- [ ] Validate isolated dependencies/safe provider sandboxes; prove no mutation or foreign data exposure.

### L4 — Harbor groups

- [ ] Implement explicit membership, per-member projections and access revalidation.
- [ ] Preserve independent failure domains, quotas and freshness during partial outages.
- [ ] Require a separately owned administrative contract for any cross-organization aggregation.

### L5 — Approved operations

- [ ] Connect core planning/approval/execution only after producer publication gates.
- [ ] Freeze scope, source revision, target, journey/schema/controller versions and exact-plan digest.
- [ ] Test replay, concurrency, cancellation, partial failure, rollback/compensation and evidence retention.
- [ ] Keep observation grants distinct from mutation; browser state and AI output never authorize actions.

### L6 — Public reference and production qualification

- [ ] Define synthetic public workloads and private deployment bindings; begin with GCP Kubernetes/Cloud SQL.
- [ ] Add a second cloud and tested source/target compatibility before claiming multicloud recovery coverage.
- [ ] Cover Kubernetes, databases, IaC, backup, secrets, networking, regions and recovery targets through explicit scenarios.
- [ ] Recover Lens with independent artifacts, identity, backups/evidence and a qualified CLI while Lens is unavailable.
- [ ] Meet production API/frontend coverage dimensions and applicable E2E scenarios; report deployment and recovery evidence separately.

## Acceptance evidence

| Boundary | Required evidence |
| --- | --- |
| Browser sandbox | Actual artifact, no live destinations/persisted credentials, safe input, keyboard/mobile/i18n, synthetic history |
| Scope/federation | Absent foreign rows/counts/graphs/cache/evidence, revocation and concurrent-scope scenarios |
| Collection | Read-only capability, bounded egress/cardinality, cancellation, stale/partial results and retention |
| Operation | Independent exact-plan approval, scoped idempotency, replay/concurrency denial and immutable evidence |
| Resilience | Lens outage cannot prevent emergency CLI, vault access or evidence recovery |

Design references only: [OpenLens](https://github.com/MuhammedKalkan/OpenLens), [Terraboard](https://github.com/camptocamp/terraboard), [Stripe sandboxes](https://docs.stripe.com/sandboxes). No runtime or code from these projects is adopted.
