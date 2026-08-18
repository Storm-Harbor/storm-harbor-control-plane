# RTO, RPO, and evidence

## Measurement

RTO and RPO require measurement points defined for each service.

```joint
flowchart LR
  I[Incident or exercise declared<br/>T0]
  S[Recovery starts<br/>T1]
  D[Dependencies restored<br/>T2]
  H[Service health gate passes<br/>T3]
  B[Business validation passes<br/>T4]

  I --> S --> D --> H --> B

  classDef event fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef technical fill:#fbe6e8,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef business fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  class I,S event;
  class D,H technical;
  class B business;
```

- **Technical recovery time:** `T3 - T0`.
- **Business recovery time:** `T4 - T0`.
- **Execution time:** `T3 - T1`; operationally useful, but not a replacement for RTO.
- **Actual RPO:** `T0 - timestamp of the applied recovery point`.

Approval or diagnostic pauses must not disappear from the metric. Record intervals separately and present both total time and a breakdown.

## Evidence bundle

Each run produces a manifest and phase documents:

```text
evidence/<run-id>/
├── manifest.json
├── effective-plan.redacted.json
├── policy-decisions.json
├── inventory-snapshot.json
├── phases/
│   ├── foundation.json
│   ├── database.json
│   └── kubernetes.json
├── validations/
│   ├── service-health.json
│   └── data-integrity.json
└── signatures/
    └── manifest.sig
```

The manifest records the schema version, tenant, run, mode, scenario, source revisions, execution-plan digest, timestamps, objectives, actuals, result, artifact digests, and executor version. Secret values and personal data are redacted before persistence.

## Integrity and retention

1. normalize documents with deterministic serialization;
2. calculate the SHA-256 digest of each artifact;
3. build the manifest with those digests;
4. sign the manifest with a KMS key outside the target;
5. write it to a versioned bucket with retention;
6. validate signature and completeness when querying or auditing.

Logs help diagnosis, but are not sufficient evidence on their own. Evidence must link the input, policy decision, actor/identity, action, target, timestamps, and outcome.

## Exercise outcome

| Status | Condition |
| --- | --- |
| `PASSED` | health/business gates pass, RTO/RPO are met, and the bundle is intact |
| `PASSED_WITH_FINDINGS` | service recovered within objectives, with nonblocking gaps |
| `FAILED` | recovery or validation fails, or an objective is exceeded |
| `ABORTED` | kill switch, cancellation, or revoked approval |
| `INCONCLUSIVE` | evidence or a critical source is unavailable; never equivalent to success |
