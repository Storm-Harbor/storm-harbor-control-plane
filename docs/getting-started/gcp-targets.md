# GCP exercise and restore targets

GCP is the first provider design example. This guide does not provision resources: the adapter and cloud acceptance suite are planned. The core schema remains provider-independent; other adapters will map their account/subscription/compartment boundaries explicitly.

## Choose purpose and destination separately

| Choice | New target | Existing target |
| --- | --- | --- |
| Exercise | Request a new isolated project; ephemeral resources, expiry, proven run-owned cleanup | Use an authorized isolated sandbox; preserve pre-existing project/resources |
| Restore | Request a new retained recovery project | Restore into an explicitly approved retained environment |

A GCP project is a global resource container, not a region. The provider binding resolves organization/folder, billing and target project creation/reuse. The placement binding resolves resource regions/zones, network, data residency and service-specific restrictions. Different resources in one project may have different location scopes. The product's organization/project IDs are not these GCP identifiers.

For a regional outage, choose an available permitted recovery region and prove backup/network compatibility. A provider-wide outage also requires an independent provider/control-plane/identity path and tested data portability; creating another GCP project is not sufficient.

## New-project preflight

Before an approved create operation, the future GCP adapter must:

1. Resolve customer-owned provider/placement aliases and verify the caller's effective product scope.
2. Verify a permitted parent organization/folder, available project ID and project-count quota; record run ownership without exposing private identifiers publicly.
3. Check project creation permission separately from billing-account association and required billing availability. Permissions must be explicitly granted to a bootstrap workload; the target cannot grant its own initial authority.
4. Confirm Resource Manager, Service Usage and required service APIs/permissions on their actual consumer/quota projects. An enabled API does not imply IAM permission, and the credential quota project can differ from the target project.
5. Verify service quotas/capacity, organization policies, allowed locations, networking, encryption, external secret access and the reviewed budget ceiling. Budget alerts alone are not a hard spending limit.
6. Use the existing infrastructure owner for creation and changes. If HCP Terraform is the source of execution, create a separately reviewed DR workspace/root/state binding with explicit manual apply policy; do not repurpose production state.
7. Persist operation identity and ownership before retries. Poll supported long-running operations with required read/list/cancel/delete permissions as applicable to that service, bounded timeout and cancellation. Never blindly retry a create with an unknown outcome.
8. Verify actual target identity against the protected source, then prove readiness before restore/deployment.

Do not infer that a `create` field grants IAM, attaches billing, starts HCP runs or installs an adapter. These are separate planned actions with separate policies and acceptance evidence.

## Restore, verify and clean up

An exercise restores a verified recovery point into its isolated target, deploys approved immutable application revisions and performs service/data checks before reporting measured recovery objectives. Synthetic/sanitized data is the default; any sensitive production backup requires explicit access, residency and isolation approval. Restoration must not redirect production traffic automatically.

For `run-owned` cleanup, record which real resources were created by the run and independently verify that ownership at deletion time. Labels alone are insufficient deletion authority. Existing sandboxes and shared resources must survive. New-project deletion needs proof that the project itself is run-owned and no unrelated resources have become dependencies. Test expiry, denied cleanup, partial failures and orphan detection before scheduled exercises.

A retained restore has no automatic cleanup timer. Data validation, availability and traffic cutover have independent approval gates. The evidence must distinguish a restored resource from a fully recovered business service.

## Authoritative references

- [Creating and managing projects](https://docs.cloud.google.com/resource-manager/docs/creating-managing-projects)
- [Cloud Billing access control](https://docs.cloud.google.com/billing/docs/how-to/billing-access)
- [Enable and disable services](https://docs.cloud.google.com/service-usage/docs/enable-disable)
- [Workload Identity Federation with deployment pipelines](https://docs.cloud.google.com/iam/docs/workload-identity-federation-with-deployment-pipelines)

Provider permission/capability matrices and sandbox evidence will be versioned with the adapter. These references do not constitute a tested GCP implementation.
