# YAML structure and examples

The initial provider-neutral contract is `recovery.stormharbor.io/v1alpha1`, kind `RecoveryConfiguration`. Its owner is [storm-lib-recovery-contracts](https://github.com/Storm-Harbor/storm-lib-recovery-contracts). Immutable [v0.1.0-alpha.1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/releases/tag/v0.1.0-alpha.1) is published and its public Go checksums are verified. Implementation and release evidence are recorded in [issue #1](https://github.com/Storm-Harbor/storm-lib-recovery-contracts/issues/1).

This contract validates intent offline. It does not execute the journey, resolve aliases, fetch secrets or provision a project. The legacy `dr.config.example.yaml` in this documentation repository remains a separate, GCP-shaped **simulation** contract. Neither document becomes a Kubernetes object merely by using `apiVersion` and `kind`.

## First exercise

Every identifier below is synthetic. Download the [exercise YAML](examples/exercise.yaml) or read it here:

```yaml
apiVersion: recovery.stormharbor.io/v1alpha1
kind: RecoveryConfiguration
metadata:
  name: service-exercise
  scope:
    organizationId: org-example
    projectId: project-example
spec:
  purpose: exercise
  source:
    name: primary-environment
    scope:
      organizationId: org-example
      projectId: project-example
  target:
    name: exercise-environment
    providerRef: gcp-recovery
    provisioning: create
    placementRef: exercise-europe
    lifecycle: ephemeral
    cleanup: run-owned
    ttlSeconds: 3600
  secrets:
    - name: recovery-database
      scope:
        organizationId: org-example
        projectId: project-example
      purpose: restore-database
```

| Field | Meaning |
| --- | --- |
| `metadata.name` | Name of this declaration, within its authorized scope |
| `metadata.scope` | Storm Harbor organization/project; must match independently authenticated context |
| `spec.purpose` | `exercise` or `restore`; does not grant operation permission |
| `spec.source` | Scoped alias for the protected source environment |
| `target.name` | Logical target alias, distinct from source |
| `target.providerRef` | Scoped alias for a provider binding and its capabilities |
| `target.provisioning` | `create` requests a new target; `existing` selects an independently verified existing target |
| `target.placementRef` | Alias for environment/region/zone/network/data-residency choices resolved by the future adapter |
| `target.lifecycle` | `ephemeral` for an exercise; `retained` for restore |
| `target.cleanup` | `run-owned` for an exercise; `retain` for restore |
| `target.ttlSeconds` | Positive exercise expiry request; omitted or zero for restore |
| `spec.secrets` | Optional scoped aliases and bounded purposes; no secret value or vault path |

Cloud-specific account identifiers and actual region names are deliberately absent from this core contract. Bindings map aliases to those details under customer control. A provider alias is not evidence that a corresponding cloud adapter is implemented. Placement and provider references inherit the configuration's scope.

## A retained restore

The [restore YAML](examples/restore.yaml) uses an existing target and a different provider alias to demonstrate a provider-independent shape. It does not claim AWS recovery support or cross-cloud service compatibility. The target section is:

```yaml
spec:
  purpose: restore
  # The complete document must also include its scoped source.
  target:
    name: recovery-environment
    providerRef: aws-recovery
    provisioning: existing
    placementRef: recovery-region
    lifecycle: retained
    cleanup: retain
```

This excerpt is not a complete configuration. Restore can also request `create`; an exercise can use `existing`. For an existing exercise sandbox, cleanup must preserve the sandbox/project itself and all pre-existing resources. Even for a new target, deletion requires independent proof of run ownership and current authorization. TTL alone does not run a janitor or authorize deletion.

## Validation boundaries

The alpha Go producer's `DecodeYAML(document, trustedScope)` accepts exactly one document up to 64 KiB, with depth limited to 32. It rejects unknown or duplicate fields, case-mismatched field names, anchors/aliases, custom tags, multiple documents and invalid scalar types. Use strings for identifiers and integer seconds for expiry. Errors return stable codes/localization keys and safe field names, never the supplied YAML values.

The caller must establish `trustedScope` from authenticated identity, authoritative membership and policy. Copying scope from the same YAML into that argument defeats the boundary. Schema validation checks shape; Go validation additionally checks scope equality, duplicate secret aliases and equal source/target aliases. Neither detects different aliases resolving to the same actual cloud resource: runtime adapters must verify that before mutation.

The canonical JSON Schema is versioned in the producer's `schemas/recovery-configuration.schema.json`. Its `$id` is a logical identifier, not a deployed schema download endpoint. Pin the file and Go module to the same immutable release, currently `v0.1.0-alpha.1`. No executable recovery CLI command is published yet; the [stage checklist](../SPECS-CONTROL-PLANE.md) describes its delivery boundary.

## Keep planned resources separate

Policies, dependency DAGs, approvals, backup selection, custom resource definitions and AI profiles are planned contracts. They are discussed in [policies and recovery journey](policies-and-journey.md). Do not add `spec.policy`, `spec.journey`, `spec.ai` or a vault path to this alpha resource: unknown fields are rejected. Their owning contracts, versions and migrations must be delivered before use.
