# Simulation pipelines

## Purpose

The pipelines are an acceptance harness for contracts and guardrails. They imitate recovery with synthetic data, produce verifiable evidence, and **do not** run `terraform`, `gcloud`, `kubectl`, or `flux`.

They all have:

- `workflow_dispatch`, a default-on `dry_run` switch, and the literal confirmation `SIMULATE`;
- `contents: read`, without `id-token: write`;
- inputs passed through the environment, never interpolated directly into the shell;
- a pinned and verified runtime;
- explicit, ordered phases;
- the shared `scripts/simulate-recovery.mjs` engine;
- a JSON artifact with declared retention;
- an operator-readable Job Summary.

## Main flow

```joint
flowchart TD
  M[Manual dispatch]
  C{Confirmation equals SIMULATE?}
  S{Select scenario}
  F[Full project loss]
  G[GCP project loss]
  T[Terraform drift]
  I[IAM recovery]
  D[Database recovery]
  K[GKE recovery]
  X[FluxCD bootstrap]
  E[Structured evidence]
  R[Job summary and artifact]
  Z[Reject without credentials]

  M --> C
  C --> S
  C --> Z
  S --> F
  S --> G
  S --> T
  S --> I
  S --> D
  S --> K
  S --> X
  F --> E
  G --> E
  T --> E
  I --> E
  D --> E
  K --> E
  X --> E
  E --> R

  classDef entry fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef gate fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef scenario fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  classDef proof fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef reject fill:#fdebed,stroke:#b0000c,color:#1d1d1b,stroke-width:2px;
  class M entry;
  class C,S gate;
  class F,G,T,I,D,K,X scenario;
  class E,R proof;
  class Z reject;
```

## Scenario catalog

| Workflow | Specific guardrails | Simulated phases |
| --- | --- | --- |
| `disaster-recovery-simulation.yml` | exact routing and per-repository concurrency | delegates one scenario |
| `simulate-gcp-project-loss.yml` | `-dr-sandbox` target | context, billing/labels, project, APIs, KMS, backend, handoff |
| `simulate-terraform-drift.yml` | allowlisted scenario and zero mutations | source SHA, backend metadata, refresh-only plan, classify, policy |
| `simulate-iam-recovery.yml` | snapshot within RPO and no private keys | catalog, normalize, redact, integrity, dependency order, restore |
| `simulate-database-recovery.yml` | new `-dr-restore` target and recovery point within RPO | catalog, checksum, isolated restore, connectivity, schema/query |
| `simulate-gke-recovery.yml` | `-dr-sandbox` cluster | target UID, network, quota, cluster, WI, Operator, policy/health |
| `simulate-flux-bootstrap.yml` | relative path without traversal | cluster, revision, controllers, sources, Kustomize, Helm, health |
| `simulate-ai-analysis.yml` | provider, tier, reasoning, task type, and data class in closed enums | data policy, model resolution, structured analysis, validation, human review |
| `simulate-full-recovery.yml` | dependency order and `always` cleanup | composes six workflows and validates every artifact |

Select `all-pipelines` in the main workflow to execute the full six-stage recovery chain and the AI-analysis pipeline in one manual dry run.

## Full project loss

```joint
flowchart LR
  F[Foundation]
  I[IAM]
  T[Terraform]
  D[Database]
  K[GKE]
  X[FluxCD]
  A[Aggregate evidence]
  C[Cleanup verification]

  F --> I --> T --> D --> K --> X --> A
  F --> C
  I --> C
  T --> C
  D --> C
  K --> C
  X --> C

  classDef phase fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef evidence fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef cleanup fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  class F,I,T,D,K,X phase;
  class A evidence;
  class C cleanup;
```

The aggregate job downloads component artifacts, rejects a missing scenario, any result other than `PASSED`, or any external mutation, preserves order, and creates a manifest with its own digest.

## Mock-engine contract

Main inputs:

| Variable | Rule |
| --- | --- |
| `DR_CONFIRMATION` | must be `SIMULATE` |
| `DR_DRY_RUN` | defaults to `true`; any other value fails closed because live adapters are unavailable |
| `DR_SCENARIO` | closed enum |
| `DR_PHASES` | nonempty list without duplicates |
| `DR_TARGET` | suffix and format compatible with the scenario |
| `DR_*_AGE_MINUTES` | nonnegative integer within the target RPO |
| `DR_GITOPS_PATH` | relative and without path traversal |
| `DR_EVIDENCE_DIRECTORY` | artifact directory; internal runner value |

Output JSON:

```json
{
  "schemaVersion": "1.0.0",
  "run": {
    "mode": "simulation",
    "scenario": "database-recovery",
    "result": "PASSED"
  },
  "objectives": {
    "targetRpoMinutes": 15,
    "actualRpoMinutes": 5
  },
  "guardrails": {
    "cloudIdentityRequested": false,
    "externalMutations": 0
  },
  "phases": [],
  "integrity": {
    "algorithm": "sha256",
    "signature": "not-signed-mock"
  }
}
```

`not-signed-mock` is intentional and prevents the scaffold from being mistaken for production evidence.

### Simulate an AI provider

```bash
DR_CONFIRMATION=SIMULATE \
DR_DRY_RUN=true \
DR_AI_PROVIDER=codex-cloud \
DR_AI_MODEL_TIER=balanced \
DR_AI_REASONING_LEVEL=medium \
DR_AI_TASK_TYPE=recovery-plan-review \
DR_AI_DATA_CLASSIFICATION=metadata-only \
node scripts/simulate-ai-provider.mjs
```

The same workflow accepts `gemini-vertex-ai` or `gemini-api`. In mock mode, it resolves the model and validates every policy without reading credentials or calling the provider.

## Run locally

```bash
DR_CONFIRMATION=SIMULATE \
DR_DRY_RUN=true \
DR_SCENARIO=gke-recovery \
DR_TARGET=sample-gke-dr-sandbox \
DR_PHASES=target-identity,private-network,cluster,operator-install,policy-health \
node scripts/simulate-recovery.mjs
```

The file will be created at `artifacts/evidence/gke-recovery.json`.

## Evolution toward real adapters

Replacing the mock with real execution must preserve the same phase and evidence contract while adding:

1. read-only planning and a serialized diff;
2. an approval environment for `mode=recovery`;
3. OIDC/WIF only in the mutating job;
4. action images and actions pinned by digest/SHA;
5. a policy check immediately before each mutation;
6. timeout, classified retry, and an idempotency key;
7. validation independent of the executor;
8. cleanup that is always executed and evidenced;
9. KMS signing and upload to immutable storage.

!!! warning "Do not enable scheduling yet"
    Periodic execution should be introduced only after budget, quotas, isolation, cleanup under failure, the kill switch, and emergency stop have been proven in disposable environments.
