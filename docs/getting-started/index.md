# Getting started

Start here to follow Storm Harbor while it is being built. The intended experience is a private configuration repository, an organization-installed GitHub App and versioned Go executors running in GitHub Actions or customer-owned runners. The configuration describes intent; independent identity, policy and provider checks authorize each operation.

Use the [complete plan and checklists](../SPECS-CONTROL-PLANE.md) to follow every implementation/adoption stage, its dependencies and its acceptance gates.

## What you can use today

| Capability | Availability |
| --- | --- |
| Public website and English documentation | Published |
| Credential-free GCP-shaped simulations | Available in this repository; see [simulation pipelines](../operations/pipelines.md) |
| Provider-neutral Go configuration validation | Immutable [v0.1.0-alpha.1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/releases/tag/v0.1.0-alpha.1), verified through the public Go proxy/checksum database; offline library only |
| Offline Go CLI | [Source candidate and exact usage](offline-cli.md); no published binary release yet |
| GitHub App installation, published executor binaries, vault resolution, project provisioning and real recovery | Planned; no installation link or production recovery command is available yet |
| Continuous controller, custom resource registry, Kubernetes CRDs and private AI adapters | Planned |
| Per-repository recovery dashboard, historical reports, gap checklist and diagrams | Accepted design; [adoption/publication guide](recovery-dashboard.md), not implemented |

The [Journey #3](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/3) and [Epic #4](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/4) own this foundation. GitHub issues and linked PRs carry current delivery status; this guide describes the capability boundary. Configuration validation is not proof that recovery works.

## Follow the adoption journey

1. **Choose your boundary.** Identify the organization and project authorized to manage recovery. These product identifiers are separate from a cloud project, account, subscription or region.
2. **Prepare a private central repository.** Keep recovery configurations and reviewed policies there. Reference only explicitly authorized application/infrastructure repositories. They stay in the customer's organization.
3. **Describe a safe target.** Start with an exercise. Choose a new target or an existing isolated sandbox and a logical placement alias. Set an expiry and restrict cleanup to resources proven to belong to the run.
4. **Prepare external bindings.** Keep credentials, sensitive provider identifiers and vault paths in customer-controlled systems. YAML uses scoped aliases only. Decide whether execution must stay on an isolated private runner.
5. **Review the declaration.** Use the [YAML guide](yaml.md), [policy and journey guide](policies-and-journey.md) and [GCP walkthrough](gcp-targets.md). The Go producer validates structure and scope; future adapters must validate actual identities and capabilities too.
6. **Adopt execution when delivered.** Install the future App for selected repositories, pin the executor release/digest, validate and approve an immutable plan, then run an isolated exercise and inspect evidence. These execution steps are the accepted design, not available automation today.
7. **Publish an optional private dashboard.** Choose approved hosting and audience controls, then generate run history, gaps and diagrams from safe evidence projections. Repository visibility and website visibility require separate verification; see the [dashboard guide](recovery-dashboard.md).

No mandatory AI, Kubernetes cluster or Storm Harbor hosted service is required by the design. The customer chooses cloud adapters, identity, evidence storage and optional private AI/RAG/CLI.

## Suggested repository layout

This is an adoption convention, not a directory discovery feature implemented by the validator:

```text
private-recovery-config/
  recovery/
    configurations/
      service-exercise.yaml
      service-restore.yaml
    policies/              # reviewed requirements; policy engine is planned
    journeys/              # recovery order/runbooks; journey resource is planned
  docs/
    ownership.md
    emergency-access.md    # procedures and aliases only; never recovery keys
  .github/
    CODEOWNERS
    workflows/             # future pinned executor integration
```

Store operational evidence in the customer's approved evidence store, outside the modeled failure domain. Do not commit evidence, state files, backups, credentials or private knowledge to a public Storm Harbor repository. Even a private configuration repository is not a credential store.

## Reading order

- [YAML structure and examples](yaml.md): the exact initial contract and its limits.
- [Offline Go CLI](offline-cli.md): reviewed source builds, local validation, output/exit behavior and publication gates.
- [Policies and recovery journey](policies-and-journey.md): what must happen before and after validation.
- [GCP targets](gcp-targets.md): creating a new project versus restoring into an existing target.
- [Recovery repository dashboard](recovery-dashboard.md): history, gap-analysis checklists, diagrams and private publication.
- [GitHub App and executor architecture](../architecture/github-execution.md): events, runners and ownership.
- [External secrets and recovery identity](../architecture/external-secrets.md): vaults, temporary access and emergency operation.
- [Accepted decisions](../architecture/accepted-decisions.md) and [roadmap](../roadmap.md): why the design looks this way and the next delivery boundaries.
