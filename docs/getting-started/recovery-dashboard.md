# Recovery repository dashboard

**Accepted design; not implemented yet.** Each customer-owned recovery repository may have its own central read-only dashboard, built from reviewed reporting inputs and published through an approved host. This is separate from the public Storm Harbor website and the future operational console. Follow [Stage 08A](../SPECS-CONTROL-PLANE.md#stage-08a-repository-dashboard-and-historical-reports) for implementation and acceptance gates.

## What the dashboard shows

| View | Intended content |
| --- | --- |
| Overview | Protected service aliases, objectives, readiness gaps, latest assessment and observation time |
| Run history | Every recorded exercise/restore, timestamp, approved scenario, result and evidence availability |
| Run details | Ordered phase attempts, approvals, durations, validations, cancellation/failure and cleanup outcomes |
| Gap analysis | Requirement, finding, severity, owner, remediation reference and evidence required for closure |
| Recovery diagram | Reviewed dependencies and journey order, with observed phase outcomes for the selected run |
| Trends | Comparable measured RTO/RPO, recurring findings and changes between accepted snapshots |

Keep exercises and retained restores visibly distinct. A successful workflow or technical health probe is not proof of business recovery. Missing or stale evidence is shown explicitly. Trend calculations must not silently discard failures, change measurement definitions or combine incomparable scenarios.

The visual direction follows [stormharbor.io](https://stormharbor.io/): Storm Harbor branding, blue/ink/paper colors, clear status cards and connected diagrams. The initial reference tokens are blue `#326ce5`, ink `#111820` and paper `#f4f6f8`; a versioned reusable template and accessible status treatments are future work. Color must not be the only way to convey an outcome. Browser assets are HTML/CSS/JavaScript; the planned report generator follows the Go 1.26.x runtime baseline.

## Evidence to published history

```joint
flowchart LR
  R[Authorized run or assessment]
  E[Independent evidence store]
  P[Scoped reviewed report export]
  G[Pinned Go report generator]
  D[Private repository dashboard]
  R --> E --> P --> G --> D

  classDef source fill:#f4f6f8,stroke:#326ce5,color:#111820,stroke-width:2px;
  classDef report fill:#111820,stroke:#84a7ef,color:#ffffff,stroke-width:2px;
  class R,E source;
  class P,G,D report;
```

Text alternative: an authorized run writes evidence to an independent store. A reviewed export selects fields for a known audience. A pinned generator renders that export into the repository's protected website.

1. Record execution and assessment facts in their authoritative scoped stores, including unsuccessful outcomes and individual attempts.
2. Validate evidence provenance, integrity, completeness, audience and disclosure policy before export. Published summaries omit protected fields rather than relying on browser-side hiding.
3. Render an immutable summary for the selected run/assessment revision and update the history index without losing concurrent runs. Retain prior snapshots; corrections identify what they supersede.
4. Publish only after artifact and host-access checks. Record report generation and deployment outcomes separately from the recovery outcome.
5. Show the report generation time, observation time, source revisions, freshness and unavailable evidence. A static snapshot does not claim live status.

The exact summary, gap and diagram schemas, generator commands, workflow inputs and reporting-policy YAML will be delivered by their owning contracts. Do not add `spec.dashboard`, `spec.reporting` or guessed commands to the current strict `RecoveryConfiguration` alpha. The dashboard will consume approved projections; it will not become another recovery definition or grant execution permissions.

## Publication and privacy

GitHub Pages is one optional static host. A private repository does **not** automatically make its website private. GitHub currently requires **GitHub Enterprise Cloud** for private Pages publication; verify the organization's eligibility and the actual site visibility before enabling it. See [GitHub's site visibility documentation](https://docs.github.com/en/enterprise-cloud@latest/pages/getting-started-with-github-pages/changing-the-visibility-of-your-github-pages-site).

| Publication choice | Required boundary |
| --- | --- |
| Private GitHub Pages | Eligible organization, privately published project site and verified access for the complete report audience |
| Customer-controlled private host | Identity/access enforcement protects HTML, assets, JSON, downloads and diagrams before delivery |
| Protected local report | An authorized operator opens a generated package under the organization's local access and retention policy |
| Explicit public demonstration | Synthetic or expressly approved fields and diagrams only; separate artifact and deployment authorization |

If the chosen host cannot enforce the required audience, publication stops. It must not fall back to a public site. A JavaScript login screen, an obscure URL, `robots.txt` or a custom domain does not protect files already delivered by a public host.

Static hosting commonly grants an audience access to the whole artifact. If projects or users require different access, generate separately authorized artifacts/hosts or use an appropriately scoped authenticated service; client-side filtering is insufficient. Permission to read the configuration repository is not an independent grant to inspect every evidence object.

Exclude secret values, raw logs/state, sensitive resource identifiers, internal endpoints, private knowledge, provider tokens and signed download URLs. Check HTML, source maps, JSON, search indexes, labels, tooltips, SVG metadata and downloadable files. Diagram rendering must escape untrusted labels/links, reject executable content and avoid sending private topology to third-party rendering services. Detailed evidence links require reauthorization at the evidence service and contain no embedded credentials.

## History and failure behavior

The external evidence store remains authoritative. Its retention, integrity and availability must survive the modeled disaster independently of GitHub and the recovery target. Pages deployments, Git history and temporary Actions artifacts are not sufficient as the only historical store. A dashboard can be rebuilt from approved evidence after loss or host migration.

The report index must account for every recorded run, including partial, failed, aborted and inconclusive outcomes. Missing evidence appears as a gap. If reporting is unavailable, the execution journal still records the attempt and later reconciliation can project it; website availability does not determine recovery success. Publication failures retain the previous report with an explicit stale/publication state when possible and an independent operator alert.

Keep complete retained history in the evidence/index store and paginate or archive dashboard views within hosting limits. Retention changes and corrections are explicit policy decisions rather than deletion of inconvenient results.

## Adoption sequence

1. In [Stage 03](../SPECS-CONTROL-PLANE.md#stage-03-customer-owned-private-recovery-repository), choose the audience, host, disclosure policy and external evidence owner. A synthetic scaffold can be prepared alongside the private recovery repository.
2. Publish the reporting contracts after inventory/readiness and run/evidence semantics exist in Stages 07–08. Add a pinned generator/template and a separate publishing workflow with only the required capabilities.
3. Validate synthetic history, gap closure, diagrams, accessibility, complete-artifact privacy, unauthorized access, concurrent updates and rebuild behavior before enabling publication.
4. After each exercise or retained restore, project its approved summary, phase outcomes and findings. Keep reporting read-only; execution, approvals and private evidence access remain independently authorized.

No dashboard, customer Pages site or publication workflow is installed by this documentation change. [Getting Started](index.md) describes the capabilities available today; [the full plan](../SPECS-CONTROL-PLANE.md) tracks the remaining work.
