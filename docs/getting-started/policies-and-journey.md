# Policies and recovery journey

The owner accepted this strategy. The initial Go producer implements offline configuration validation; the policy engine, approval service, executor and journey/DAG resources are still planned. Treat the following as operational requirements and implementation guidance, not as fields accepted by `RecoveryConfiguration` today.

## Policy responsibilities

| Policy | Owner and required enforcement |
| --- | --- |
| Repository access | Customer installs the future GitHub App for selected repositories; repository scope alone does not authorize cloud operations |
| Product scope | Authenticated organization/project membership and operation permission, rechecked at protected boundaries |
| Provider binding | Explicit permitted accounts/projects, regions, networks and supported source-to-target mappings |
| Secret access | Customer vault grants the minimum secret/version/operation for the selected step and workload |
| Exercise safety | Disposable target/resource ownership, TTL, independent source protection, cost and quota ceilings |
| Restore safety | Approved recovery point and target, retained lifecycle, separate traffic/cutover approval |
| Evidence and AI | Customer storage, access, retention and egress rules; no raw credentials or implicit knowledge sharing |

Installation policy sets a ceiling. A repository cannot relax it. Evaluate the effective intersection of installation, project, environment and run policy; a feature toggle is never authorization. Human approval binds the exact plan digest, source revisions, target, allowed actions and expiry. Changes to any of those inputs invalidate approval.

An attacker who changes a workflow must not be able to grant that workflow access to the vault or approve its own production action. Use external trust/approval enforcement for sensitive operations. GitHub environment approvals can add a gate but do not replace independent vault/provider policy.

## Ordered recovery phases

1. **Validate:** parse the pinned schema, establish trusted scope and reject forbidden capabilities.
2. **Discover:** read authorized sources, backups and provider state; bind revisions/digests and verify recovery-point integrity.
3. **Preflight:** check source protection, actual target identity, APIs/IAM/quotas, secret availability, region/network compatibility, cost and emergency access.
4. **Plan:** construct dependency order and bounded actions; freeze configuration, schemas, adapters, controller images and source commits.
5. **Approve:** authorize that exact plan and target with an expiry; revalidate before obtaining mutation credentials.
6. **Provision and restore:** create or verify the isolated target, restore data, apply infrastructure and deploy approved application revisions through their established owners.
7. **Verify:** test service behavior, dependencies, data integrity, RTO/RPO and permitted traffic. A green process or HTTP liveness check alone is insufficient.
8. **Retain evidence and finish:** record sanitized outcomes outside the failure domain. Exercises clean only independently proven run-owned resources; retained restores remain available under the approved policy.

Each phase needs bounded timeout/cancellation, classified retries and durable idempotency. On partial failure, record the last authoritative result before retrying. Cancellation or a kill switch stops new unsafe actions; it cannot erase side effects already accepted by a provider. Cleanup needs its own authorization and failure handling, including a separate janitor for abandoned exercise resources. No blanket cleanup applies to retained restores.

## Infrastructure and application ownership

If Terraform/HCP already owns infrastructure, Storm Harbor requests a reviewed run through that owner. It must not issue parallel imperative mutations against the same resources. Separate bootstrap identity/foundation, target infrastructure and application release phases. Give a new exercise its own reviewed root/workspace/state boundary rather than pointing production state at a temporary project.

Record the exact Terraform configuration revision, workspace binding, provider lockfile, plan digest and apply authorization without exporting sensitive state/plan contents into evidence. Creating a project may require a bootstrap identity before the target project exists; its workload must never self-elevate. Application workflows retain ownership of immutable images/revisions, migrations and traffic changes.

Optional AI/RAG/CLI can explain failures or propose a typed plan. It cannot authorize a phase, read arbitrary vault paths or override source/target protection. The same approved journey must remain executable with AI disabled.
