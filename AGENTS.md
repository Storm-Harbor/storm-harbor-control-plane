# Storm Harbor technical documentation and simulations

Read the shared [engineering guidance](https://github.com/Storm-Harbor/.github/blob/main/AGENTS.md) and its three specifications. In the full workspace, the same guidance is available at `../../.github/AGENTS.md`.

Technical documentation always remains in English. Preserve the distinction between accepted architecture, mock fixtures and implemented runtime support. Go 1.26.x is the official language for new executable components; the Node.js 26.6.0 scripts are transitional safe simulations.

Use `docs/SPECS-CONTROL-PLANE.md` for the complete stage/dependency and acceptance checklist. Checked items require concrete release/deployment evidence; live lifecycle and ownership remain in native GitHub issues. Keep private adopter names/bindings out of the public plan.

The checked-in schema/example is a GCP-shaped simulation contract. Provider-neutral resources and optional private AI/RAG/CLI require a versioned migration, not a provider enum edit. User-defined resource/controller and Kubernetes CRD documentation does not mean those runtimes are installed or implemented.

Use the owning issue, Storm Harbor Project #1, issue branch and a user-reviewed/manually merged PR. The current base is `main`. Relevant checks are `node --run test:disaster-recovery` and `python -m mkdocs build --strict` with the pinned toolchain/dependencies. Publishing workflows are separate from manual, read-only simulation workflows. Never enable real cloud mutation as part of documentation publication.

For navigation/template/style changes, also run the pinned Playwright suite with `npm ci`, `npx playwright install chromium` and `npm run test:docs-browser` after the strict build. It reads the actual static artifact under the configured Pages origin, including instant navigation and mobile submenus. Browser-only tooling is not a production runtime.
