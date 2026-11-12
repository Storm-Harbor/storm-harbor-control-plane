# External secrets and recovery identity

**Accepted design; vault/identity adapters planned.** Keep configuration, execution identity and secret values separate. YAML carries scoped logical references. Sensitive vault paths, provider account mappings and binding metadata can remain outside GitHub in customer-controlled systems. The future resolver selects those bindings only after authentication, effective product authorization and step-level policy checks.

## Access sequence

1. Establish the authenticated installation/workload and effective organization/project scope independently of YAML.
2. Validate the declaration, actual source/target, current policy and approved immutable plan.
3. Present a short-lived workload identity to the customer's trust broker/provider where supported. Validate issuer, audience, expiry, installation/repository identity, immutable workflow identity and permitted branch/environment claims. Subject claims must match the customer's explicit trust policy.
4. Exchange it for the minimum operation/resource capability for this step. OIDC federation support and claim semantics are provider-specific; use a reviewed broker or local identity adapter when a provider cannot support the desired exchange directly.
5. Resolve only the authorized secret/version/purpose, use it for the bounded step and revoke/expire access promptly. Prefer dynamically generated database credentials when the actual database/vault integration supports them.
6. Record an allowlisted audit outcome and non-sensitive binding/version provenance. Never record secret bytes, raw provider errors, prompts containing credentials or unrestricted environment dumps.

A temporary access token for Secret Manager does not make a stored static database password temporary. Static values need independently enforced versioning/rotation and narrower runtime handling. Avoid permanent cloud keys in GitHub secrets; customer-owned recovery keys, if unavoidable, belong in the approved external emergency path with separate rotation/control.

## Prevent credential leakage

Prefer direct use in process memory; zero/bound temporary material where practical without claiming complete erasure from managed memory. If an integration requires a file, create it in an isolated temporary location with restrictive permissions, exclude it from artifacts/caches and remove it on all exit paths. Avoid command-line secret arguments, shell tracing, unrestricted child-process environments and relying solely on log masking.

If a secret cannot leave the customer's environment, the sensitive step runs on an isolated customer runner with local vault access. Only approved sanitized outcomes leave that boundary. A private runner must not run untrusted repository content with its privileged workload identity.

The App signing key and bootstrap/emergency identities are also sensitive credentials, even though they are not declared in recovery YAML. Separate their storage/rotation from ordinary job access. Repository/workflow compromise must not independently change the vault trust policy, approve a privileged restore or overwrite trusted backups.

## Survive the modeled failure

Vault, trust broker, KMS/decryption access, backups, evidence and emergency execution must remain reachable outside the failure domain being recovered. Recovery cannot depend exclusively on an identity service or network that the same plan is rebuilding. Define a customer-controlled emergency access path with narrow permissions, expiration, independent approval and durable audit. Test it in isolated exercises without committing recovery keys or secret values.

## Required future acceptance

Reject wrong issuer/audience/repository/workflow/scope, expired identities, revoked membership, unapproved targets and cross-project alias resolution. Prove redaction across logs/artifacts/AI, no credential persistence on failed/cancelled jobs, isolated concurrent runs, bounded credential lifetime and safe operation when vault/identity is unavailable. Test real provider semantics in safe sandboxes before claiming support.

References: [GitHub OIDC](https://docs.github.com/en/actions/concepts/security/openid-connect), [Vault JWT/OIDC authentication](https://developer.hashicorp.com/vault/docs/auth/jwt), [Secret Manager authentication](https://docs.cloud.google.com/secret-manager/docs/authentication).
