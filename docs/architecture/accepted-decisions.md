# Accepted execution and recovery decisions

The product owner accepted these decisions during the November 2026 design conversation. This record preserves that authorization; no additional strategy review is required to prepare implementation. Repository PRs still follow the [manual merge process](https://github.com/Storm-Harbor/.github/blob/main/specs/SPECS-DELIVERY.md). Acceptance of a design is not evidence that its runtime exists.

| Decision | Rationale and consequence |
| --- | --- |
| GitHub App + private central configuration repository | Reuse reviewed GitHub desired state and selected repository access without making customer projects public |
| Bounded Go executors in hosted or self-hosted runners | Start without requiring a permanent recovery agent; private execution remains a customer choice |
| Federated domain authority | GitHub declares intent; providers, state backends, backups, vaults and evidence retain their authoritative data |
| Logical secret aliases + external bindings | Keep credentials and optionally sensitive identifiers outside GitHub; use scoped temporary identity and independent policies |
| Exercise/restore distinct from target create/reuse and placement | Support disposable drills and retained recovery without conflating cloud containers with regions |
| GCP first target design; provider-independent core | Validate an initial concrete path while preserving AWS, Azure, DigitalOcean and OCI adapter boundaries |
| Respect infrastructure and application owners | Terraform/HCP governs its resources/state/apply approvals; application workflows govern immutable releases and traffic |
| English adoption docs alongside implementation | Getting Started and YAML examples expose the exact current boundary; GitHub holds live delivery status |
| Optional dashboard for each recovery repository | Publish historical exercise/restore summaries, gap checklists and diagrams using the Storm Harbor identity; independently verify the site's audience and retain authoritative evidence outside GitHub |
| Primary implementation repository: `storm-harbor-control-plane` (2026-11-13) | Keep CLI, GitHub App/bots, orchestration, APIs, adapters, controllers/CRDs and internal packages together; sites and user-facing apps remain separate. Preserve the published contract module; no new core repository is the default |
| Product website and Harbor Lens (2026-11-13) | The website stays in `storm-harbor-app` at `stormharbor.io`. Our separate `harbor-lens` application starts with a synthetic sandbox; private observation, authorized Harbor groups and approved operation follow published core contracts. No live monitoring capability is claimed |
| Customer-selected opt-in analysis agent (2026-11-13) | Plan bounded adapters for Copilot, Codex, Claude Code or customer agents with the same permitted repositories, data policy, proposal contract and review gates; agent choice never expands bot permissions |

The approved reference assessment uses neighboring infrastructure, application, AI, IoT and Go-service patterns as design inputs only. Private project IDs, workspace bindings, account details and operational state are not public examples or product dependencies. The explicitly excluded workspace is outside this delivery.

Delivery: [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3) → [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4) → [documentation #5](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/5) and [producer library #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1). The first slice delivers these records and offline validation. App/CLI/provider consumers wait for a manually merged, validated immutable producer release and separate owning issues.
