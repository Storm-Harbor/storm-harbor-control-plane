# AI providers: Codex Cloud and Gemini

## Purpose and boundary

AI in Storm Harbor Control Plane is an **advisory** layer for interpreting data that has already been discovered. It receives no cloud authority, does not run Terraform, does not apply manifests, and does not approve recovery.

Allowed use cases:

| AI task type | Minimum input | Structured output |
| --- | --- | --- |
| `inventory-classification` | asset metadata and labels | service, type, suggested criticality, and confidence |
| `dependency-inference` | redacted technical relationships | suggested edges, confidence, and explanation |
| `recovery-plan-review` | redacted effective plan | gaps, severity, and mitigation |
| `evidence-summary` | manifest and phase results | executive summary and findings |

Every suggestion retains provenance, confidence, and `PROPOSED` status until human confirmation or an independent deterministic rule validates it.

## Provider-neutral architecture

```joint
flowchart LR
  P[Pipeline or control plane]
  D[Data policy and redaction]
  R[Profile resolver]
  S[Structured prompt contract]
  O[Codex Cloud adapter]
  G[Gemini Vertex AI adapter]
  A[Gemini API adapter]
  V[Schema and safety validation]
  H[Human review]
  E[AI evidence]

  P --> D --> R --> S
  S --> O
  S --> G
  S --> A
  O --> V
  G --> V
  A --> V
  V --> H --> E

  classDef pipeline fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  classDef policy fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef provider fill:#eeeeee,stroke:#777777,color:#1d1d1b,stroke-width:2px;
  classDef proof fill:#eeeeee,stroke:#666666,color:#1d1d1b,stroke-width:2px;
  class P pipeline;
  class D,R,S policy;
  class O,G,A provider;
  class V,H,E proof;
```

## Separate settings

Do not conflate quality, reasoning depth, and AI responsibility:

| Dimension | Values | Function |
| --- | --- | --- |
| `provider` | `codex-cloud`, `gemini-vertex-ai`, `gemini-api` | endpoint, identity, and adapter contract |
| `modelTier` | `efficient`, `balanced`, `frontier`, `custom` | portable capability/cost/latency profile |
| `model` | resolved or custom ID | exact model sent to the provider |
| `reasoningLevel` | `low`, `medium`, `high` | portable reasoning depth shared by adapters |
| `taskType` | four allowlisted types | limited responsibility of the call |
| `mode` | `advisory`, `evaluation` | assisted production use or offline evaluation |

The tier must not silently select “the newest model.” A versioned registry resolves tier → model; the run freezes the model, registry version, and reasoning level in the execution plan and evidence bundle.

## Initial registry

| Provider | Efficient | Balanced | Frontier |
| --- | --- | --- | --- |
| Codex Cloud | `gpt-5.6-luna` | `gpt-5.6-terra` | `gpt-5.6-sol` |
| Gemini | `gemini-3.5-flash-lite` | `gemini-3.6-flash` | `gemini-3.1-pro-preview` |

The IDs are initial defaults, not a permanent contract. Changes require evaluations and a new registry version. Preview models must be blocked in real recovery unless an explicit policy permits them.

OpenAI documentation positions Luna for high volume, Terra for balance, and Sol for maximum capability; the GPT-5.6 family exposes `none`, `low`, `medium`, `high`, `xhigh`, and `max` effort levels. Storm Harbor Control Plane initially uses the portable intersection `low|medium|high`. [OpenAI model guidance](https://developers.openai.com/api/docs/guides/latest-model/)

Current Gemini models use `thinking_level`, with support depending on the model; the portable layer also restricts values to `low|medium|high`. [Gemini thinking](https://ai.google.dev/gemini-api/docs/thinking)

## Reasoning-level translation

```joint
flowchart TB
  L[Portable reasoningLevel]
  LOW[low]
  MED[medium]
  HIGH[high]
  OE[OpenAI reasoning.effort]
  GT[Gemini thinking_level]

  L --> LOW
  L --> MED
  L --> HIGH
  LOW --> OE
  MED --> OE
  HIGH --> OE
  LOW --> GT
  MED --> GT
  HIGH --> GT

  classDef portable fill:#f7d7da,stroke:#e30613,color:#1d1d1b,stroke-width:2px;
  classDef provider fill:#f4f4f4,stroke:#c0000d,color:#1d1d1b,stroke-width:2px;
  class L,LOW,MED,HIGH portable;
  class OE,GT provider;
```

For `custom`, the adapter checks its capabilities and rejects an invalid combination before sending content. Never send `thinking_level` and `thinking_budget` to Gemini simultaneously.

## Configuration example

```yaml
ai:
  enabled: false
  mode: advisory
  defaultProfile: codex-balanced
  allowedTaskTypes:
    - inventory-classification
    - dependency-inference
    - recovery-plan-review
    - evidence-summary
  profiles:
    - name: codex-balanced
      provider: codex-cloud
      modelTier: balanced
      model: gpt-5.6-terra
      reasoningLevel: medium
    - name: gemini-balanced
      provider: gemini-vertex-ai
      modelTier: balanced
      model: gemini-3.6-flash
      reasoningLevel: medium
  dataPolicy:
    redactSecrets: true
    allowSourceCode: false
    allowInfrastructureState: false
    maxPromptBytes: 65536
    retention: provider-zero-data-retention
  enforcement:
    default: advisory
    blockRecoveryOnProviderFailure: false
```

## Identity and secrets

- `codex-cloud`: token/API credential resolved at runtime through a secret reference;
- `gemini-vertex-ai`: prefer OIDC/WIF and a short-lived GCP identity;
- `gemini-api`: API key only in a secret manager and only when Vertex AI is unsuitable;
- no token enters YAML, GitHub input, logs, or artifacts;
- the current mock job neither requests nor reads any of these credentials.

## Data policy

Before the call:

1. reduce data to the minimum required for the `taskType`;
2. remove secrets, tokens, kubeconfig files, Terraform state, and personal data;
3. classify it as `metadata-only` or `redacted-config`;
4. apply a byte limit and field allowlist;
5. calculate a digest of the redacted input;
6. require JSON Schema output;
7. validate references to existing assets/findings;
8. send it for human review.

Do not persist internal reasoning. Record the redacted input, structured final output, usage, latency, effective model, policy decision, and digests.

## Provider failure

AI unavailability must never block an already approved emergency recovery. `required-for-planning` may prevent a plan from being published before an incident, but the frozen executable run must remain independent of the provider. During incidents, an AI failure creates a finding and triggers a deterministic fallback.

## Evaluations before activation

- classification and dependency-edge accuracy;
- false positives in critical findings;
- schema compliance;
- secret-leakage canaries;
- grounding in supplied asset IDs;
- latency, tokens, and cost per task type;
- comparison by provider/tier/reasoning level;
- regression testing when the registry changes models.
