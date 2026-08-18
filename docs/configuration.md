# Configuration plan

## Decision

Use a versioned YAML document as the primary contract, validated by JSON Schema. Environment variables serve only to bootstrap the process and provide emergency overrides for a run; secret values belong in a secret manager.

An environment-variable-only configuration looks simple, but becomes difficult to validate and audit across multiple repositories, projects, phases, and feature toggles. YAML keeps relationships explicit; JSON Schema provides validation and allows invalid plans to be rejected before execution.

`dr.config.yaml` is the global Storm Harbor Control Plane contract, even though it uses the `apiVersion`/`kind` convention. It does not automatically become a Kubernetes object. The control plane validates and freezes this plan, then materializes only the Kubernetes portion as `DisasterRecoveryPlan`, `RecoveryTarget`, and `RecoveryRun` after the corresponding CRDs are installed on the API server.

## Configuration layers

1. **Internal defaults:** safe values; mutating capabilities start disabled.
2. **Repository plan:** `dr.config.yaml`, reviewed as infrastructure code.
3. **Installation policy:** allowlists and organizational ceilings.
4. **Runtime overrides:** a small set, valid for one run.
5. **Secret resolution:** values fetched through workload identity, never persisted or printed.

The installation policy is a ceiling: a repository cannot enable a capability prohibited by the organization.

## AI profiles

`spec.ai` keeps the provider, tier, exact model, reasoning level, and task type separate. Configuration starts with `enabled: false` and `mode: advisory`. The matrix and translation rules are documented in [AI providers](architecture/ai-providers.md).

Model IDs are not secrets; credentials are resolved at runtime. Terraform state, kubeconfig files, tokens, and secrets are prohibited in prompts. The control plane records only redacted input, structured output, the effective model, metrics, and digests.

## Allowed variables

| Variable | Purpose | Secret |
| --- | --- | --- |
| `DR_CONFIG_PATH` | Path to the plan; defaults to `dr.config.yaml` | No |
| `DR_RUN_ID` | Correlation identifier supplied by the orchestrator | No |
| `DR_LOG_LEVEL` | `error`, `warn`, `info` or `debug`; defaults to `info` | No |
| `DR_DRY_RUN` | Forces non-mutating execution; must default to `true` | No |
| `DR_GITHUB_TOKEN` | Short-lived GitHub App installation token | Yes |
| `DR_ID_TOKEN` | Short-lived OIDC token when supplied by the runtime | Yes |

Do not create environment variables for every YAML field. Do not place GCP service-account JSON keys, IAM snapshots, encryption material or kubeconfigs in repository variables.

## Feature toggles

Feature toggles control code paths, not authorization. Each operation requires:

1. enabled in the repository plan;
2. permitted by the GitHub App installation policy;
3. supported by the installed adapter;
4. authorized by workload identity;
5. allowed for the selected execution mode;
6. accepted by scenario guardrails.

Unknown toggles must fail schema validation. New toggles should default to `false` and document their required permissions.

## IAM snapshot contract

An IAM snapshot should contain normalized declarative data and recovery metadata:

- project/folder/organization resource identifiers;
- service-account identities, excluding private keys;
- roles, members, conditions and policy etags;
- Workload Identity pools, providers and bindings;
- schema version, source commit, creation time and recovery dependency order;
- SHA-256 digest and KMS key version used for envelope encryption.

The vault path in configuration is a reference, not a secret. Snapshot reads should be limited to the recovery workload identity. Snapshot writes should use a separate backup identity where practical, preventing a compromised recovery process from silently replacing trusted backups.

## Validation and lifecycle

- Validate syntax and schema on every pull request.
- Resolve repositories to immutable commit SHAs for each recovery run.
- Reject duplicate source roles and projects not matching the configured simulation labels.
- Record the effective configuration after redacting secret references.
- Sign evidence and retain it independently from the infrastructure being recovered.
- Introduce scheduled execution only after cleanup, quotas, budget alerts and kill-switch behavior are tested manually.

## Guardrail evaluation order

1. Validate the configuration schema and installation policy ceiling.
2. Check the global kill switch before requesting workload identity.
3. Resolve the source repositories to immutable commits.
4. Prove that the target project is allowlisted and carries disposable labels.
5. Estimate changes, deletions and cost; reject values above configured ceilings.
6. Confirm that the scenario cannot alter organization or folder IAM.
7. Take and verify the required recovery points before mutation.
8. Execute one phase at a time and stop on a failed health gate.
9. Run cleanup regardless of the scenario outcome.
10. Sign and retain redacted evidence outside the recovered failure domain.

Database simulations must restore into a new `-dr-restore` instance and validate connectivity, schema checksum and read-only queries. GKE simulations must target a `-dr-sandbox` cluster and require private networking and Workload Identity before FluxCD is bootstrapped.
