# GitHub App and executor architecture

**Accepted design; execution components planned.** GitHub holds reviewed desired state. A customer installs the Storm Harbor GitHub App in its organization and grants access only to selected repositories. A private central repository declares recovery intent and references authorized application/infrastructure repositories. Installation does not copy those repositories into Storm Harbor or expose them publicly.

## Execution flow

1. A reviewed commit, manual request, approved event or bounded schedule identifies the configuration repository and immutable commit.
2. The integration validates the event, installation and repository access. Future webhook handling must verify signatures, deduplicate deliveries and resist replay; event payloads do not establish product authorization.
3. A GitHub Actions workflow runs an immutable Go executor release/digest with minimum job permissions. App code uses short-lived installation tokens; the App's signing key remains in a controlled integration secret store with rotation.
4. An identity adapter establishes organization/project scope from authenticated installation membership and policy. The executor validates YAML against that independently established context.
5. The run reads authorized provider/backup state, plans, obtains independent approval, requests scoped temporary credentials, executes bounded phases and stores sanitized evidence.

An installation token for reading a repository, a GitHub Actions `GITHUB_TOKEN`, a workflow OIDC identity and a cloud/vault authorization are different credentials and trust decisions. Never use one as implicit proof of the others.

## Runner and hosting choices

| Execution choice | Intended use and boundary |
| --- | --- |
| GitHub-hosted runner | Bounded jobs where customer policy permits code/data/network access from GitHub infrastructure; available OS/architecture must match the executor release |
| Customer self-hosted runner | Private networks, local vaults or restrictions on where sensitive work may execute; isolate trusted jobs and use disposable runners where practical |
| Optional customer-hosted persistent controller | Future continuous reconciliation and operation when GitHub job scheduling is unavailable; not required for the initial bounded workflow model |

Self-hosting is not automatic isolation: separate runner groups/identities, do not execute untrusted pull-request code with privileged access, prevent workspace/credential reuse and reauthorize every operation. Use a compatible published Go binary for each supported OS/architecture; GitHub runner availability alone does not prove an adapter is portable.

The GitHub App integration itself needs an explicitly operated webhook/token service, whether hosted by the project or customer. Using Actions for execution reduces required persistent agents but does not remove App hosting, key management or audit responsibilities. The public Pages website is static documentation/presentation, not that service.

## GitOps and continuity

The model follows GitOps: compare declared intent with actual state and reconcile through deterministic controllers. A workflow performs bounded reconciliation on events/schedules; it is not a continuously running Flux controller. Kubernetes/Flux can be downstream adapters without becoming mandatory for every installation.

GitHub is authoritative for desired state and source revisions. Cloud APIs own observed resources; Terraform's backend owns managed state; backup stores own recovery points; vaults own secret values; the evidence store owns results. [Federated authority](federated-source-of-truth.md) prevents the control plane becoming an unreviewed copy of all those systems.

Freeze and retain signed plans, required artifacts and emergency runbooks outside GitHub and the modeled failure domain. A future emergency local/self-hosted executor must verify provenance, current policy and independent emergency identity before using that bundle. Do not turn GitHub unavailability into a permission bypass.

See [GitHub OIDC](https://docs.github.com/en/actions/concepts/security/openid-connect) for the workload-identity mechanism and [external secrets](external-secrets.md) for its trust constraints.
