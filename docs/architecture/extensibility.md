# Resource extensibility

Status: accepted architecture direction; the registry, runtime API, controllers and installable CRDs are not implemented yet.

## Resource definitions

Storm Harbor will let organizations define resource kinds and reconcile them with Go controllers. This follows the Kubernetes resource/controller model while covering resources outside clusters. A Kubernetes CustomResourceDefinition is installed only through the Kubernetes integration; it is not the universal execution engine.

Each definition declares a unique group/kind, served and storage versions, organization/project scope, structural JSON Schema, immutable fields, status schema and controller capability references. Definitions cannot shadow built-in kinds. Unknown fields and incompatible references are rejected before a plan is accepted.

| Resource field | Contract |
| --- | --- |
| `apiVersion`, `kind` | resolve an installed, versioned definition |
| `metadata` | stable identity, effective scope, generation and concurrency version |
| `spec` | user intent, validated before admission; never rewritten by a controller |
| `status` | observations, conditions, evidence references and `observedGeneration` |

Definitions and controller packages are versioned independently. Published versions include compatibility fixtures and migration/conversion rules. A recovery run freezes the resource revisions, schema versions and controller package digests; upgrades cannot silently change an approved plan.

## Controller contract

A Go controller observes current state, computes a bounded difference and requests an idempotent operation through authorized adapters. A resource definition alone grants neither cloud permissions nor permission to execute a controller.

Required behavior:

- scope every read, watch, queue message and mutation to the effective organization/project;
- declare readable/writable kinds, allowed target classes and required permissions;
- use durable leases, concurrency control and idempotency keys;
- distinguish transient retryable errors from permanent rejected conditions;
- check current authorization, deadlines, approvals and kill switch before mutations;
- emit independent `Accepted`, `Progressing`, `Ready`, `Degraded` and evidence conditions;
- preserve cleanup and evidence after partial failure or cancellation;
- bound finalizers and define audited recovery from a stuck finalizer;
- avoid deleting infrastructure merely because a declarative resource was removed.

Controller packages need provenance, immutable digests and an installation allowlist. Tenant extensions execute in isolated workers with bounded credentials, CPU, memory, time, network destinations and tool capabilities. Loading untrusted Go shared libraries into the control-plane process is not the extension model. Transport and packaging remain separate implementation decisions.

## Kubernetes integration

The Go Kubernetes Operator will reconcile native CRDs through the Kubernetes API, status subresources, scoped RBAC and admission validation. It reports Kubernetes outcomes to the external journey; it does not own provider-wide recovery or bootstrap the cluster that hosts it.

The existing [CRD proposals](../kubernetes/crds-and-operator.md) describe `DisasterRecoveryPlan`, `RecoveryTarget` and `RecoveryRun`. They are examples, not installed API resources. User-defined Kubernetes kinds require namespace isolation, upgrade fixtures, admission tests and explicit controller authorization.

## Acceptance before enabling extensions

Validate schema admission and version migration, cross-tenant denial, controller idempotency after crashes, event replay, lease loss, cancellation, finalizer deadlines, secret redaction and restoration of the resource registry outside the protected failure domain.
