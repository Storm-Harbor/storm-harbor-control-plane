# CLI alpha releases

The source CLI is merged and has passed native Linux, macOS and Windows acceptance. An installable binary requires a separate immutable release, checksums and verified signed build provenance. Consult the [published releases](https://github.com/Storm-Harbor/storm-harbor-control-plane/releases) for actual availability; a source merge, candidate CI artifact or this runbook is not evidence that a binary version has been published.

The manual `Release CLI alpha` workflow is prepared in the owning [CLI delivery #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14). Its workflow changes require the user's manual PR merge before dispatch from `main`. No release has been activated by this documentation change.

## Publication contract

| Gate | Required behavior |
| --- | --- |
| Dispatch | Explicit maintainer invocation from canonical repository `main`; serialized runs; bounded alpha-version input |
| Source | The exact commit passed successful `Go CLI` and `Documentation` push workflows on main |
| Identity | Signing uses GitHub workload identity; ordinary build/PR jobs keep read-only permissions and cannot sign or publish |
| Native platforms | Linux amd64, macOS arm64 and Windows amd64 must match the actual Go host; no unsupported cross-build is published |
| Executable | The versioned binary passes all nine executable scenarios, including expected version and exact source revision |
| Package | Go build tool creates deterministic tar.gz/zip packages, bounded allowlisted contents, legal notices and a digest manifest; no existing package is overwritten |
| Provenance | Both binary and archive receive signed attestations; publisher verifies repository, signer workflow, main ref, source commit and hosted runner identity |
| Publication | Verify every expected archive/bundle, generate checksums, assemble a draft, attach all seven assets and then publish one immutable prerelease |

GitHub distinguishes signed build provenance from release integrity. The workflow uses [artifact attestations](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations) for the build and [immutable release attestations](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases) for the published tag/assets. A checksum alone does not establish either identity.

The release workflow requests signing permissions only in build jobs, and release-write permission only in the verified publisher. It uses no provider credential, customer environment, GitHub App key, vault, cloud resource, AI agent or recovery execution.

## Maintainer sequence

1. Manually merge the reviewed workflow/packaging PR. Wait for Go and documentation checks on the exact new `main` commit and verify its publication.
2. Confirm repository release immutability is enabled and choose an unused alpha tag. Existing tags/drafts/releases cause a refusal; neither tags nor assets are overwritten. After artifact verification, the publisher atomically creates the new tag at the reviewed source SHA, verifies that target before creating the draft and again after immutable publication, and never moves an existing tag.
3. Dispatch from main, for example:

   ```sh
   gh workflow run release-cli.yml --repo Storm-Harbor/storm-harbor-control-plane \
     --ref main -f version=v0.1.0-alpha.1
   ```

4. Inspect the completed workflow, native binary acceptance, attestation verification, published immutable flag and exact release/source commit. A skipped, failed or partially completed job is not a released version.
5. Verify a freshly downloaded package and exercise it before adoption. Keep the native delivery issue/umbrella open while publication, signature or acceptance gates remain pending.

The Go package tool is `go run ./cmd/storm-package`; it is a maintainer build tool, not a recovery command or installer. Packages contain `storm`/`storm.exe`, `build-info.txt`, `build-manifest.json`, the Apache license and dependency/toolchain notices. Their manifest records version, source revision, target, Go version and each included file's size/SHA-256, without filesystem source paths. Only the pipeline's qualified native matrix establishes platform availability.

## Verify before use

The following is a Linux example for an **already published** version. Select the corresponding archive/bundle from the release for another supported platform. GitHub CLI must support `attestation verify`, `release verify` and `release verify-asset`; the pipeline checks this before publication.

```sh
REPO=Storm-Harbor/storm-harbor-control-plane
VERSION=v0.1.0-alpha.1
ARCHIVE="storm_${VERSION}_linux_amd64.tar.gz"
BUNDLE="storm_${VERSION}_linux_amd64.sigstore.jsonl"

gh release verify "$VERSION" --repo "$REPO"
gh release download "$VERSION" --repo "$REPO" \
  --pattern "$ARCHIVE" --pattern "$BUNDLE" --pattern SHA256SUMS
gh release verify-asset "$VERSION" "$ARCHIVE" --repo "$REPO"
sha256sum --ignore-missing --check SHA256SUMS
SOURCE_REVISION="$(gh api "repos/$REPO/releases/tags/$VERSION" --jq .target_commitish)"
gh attestation verify "$ARCHIVE" --repo "$REPO" --bundle "$BUNDLE" \
  --signer-workflow "$REPO/.github/workflows/release-cli.yml" \
  --source-ref refs/heads/main --source-digest "$SOURCE_REVISION" \
  --deny-self-hosted-runners

mkdir "storm-$VERSION"
tar -xzf "$ARCHIVE" -C "storm-$VERSION"
"./storm-$VERSION/storm" version --output json
```

The reported version/revision must match the verified release. Use [offline validation](../getting-started/offline-cli.md) with independently supplied local expected scope. Verification does not authorize a cloud action or prove recovery readiness.

## Failure and recovery

Before a draft exists, a failed attempt publishes no release; a new explicit dispatch can retry only if the version remains unused. If a draft/tag or partial upload exists, stop and inspect its source and assets. The workflow deliberately does not delete drafts, move tags, overwrite assets or silently resume unknown publication outcomes. Keep any incomplete draft unpublished until its exact provenance and asset set are reconciled under maintainer review.

After publication, verify the immutable release and its assets. Corrections require a new version. Private adopter repository creation and validation-workflow adoption follow a published/qualified CLI; they do not run as part of this release workflow. Canonical semantic normalization, cloud execution, GitHub App installation, AI adapters, dashboards and the future monitoring application retain their separate delivery gates.
