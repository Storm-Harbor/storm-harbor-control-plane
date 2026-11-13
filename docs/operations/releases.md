# CLI alpha releases

The source CLI is merged and has passed native Linux, macOS and Windows acceptance. An installable binary requires a separate immutable release, checksums and verified signed build provenance. Consult the [published releases](https://github.com/Storm-Harbor/storm-harbor-control-plane/releases) for actual availability; a source merge, candidate CI artifact or this runbook is not evidence that a binary version has been published.

The manual `Release CLI alpha` workflow is prepared in the owning [CLI delivery #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14). Its workflow changes require the user's manual PR merge before dispatch from `master`. No release has been activated by this documentation change.

## Publication contract

| Gate | Required behavior |
| --- | --- |
| Dispatch | Explicit maintainer invocation from canonical repository `master`; serialized runs; bounded alpha-version input; GitHub CLI 2.x >= 2.93 |
| Source | The exact commit passed successful `Go CLI` and `Documentation` push workflows on master |
| Identity | Signing uses GitHub workload identity; ordinary build/PR jobs keep read-only permissions and cannot sign or publish |
| Native platforms | Linux amd64, macOS arm64 and Windows amd64 must match the actual Go host; no unsupported cross-build is published |
| Executable | The versioned binary passes all nine executable scenarios, including expected version and exact source revision |
| Package | Go build tool creates deterministic tar.gz/zip packages, bounded allowlisted contents, legal notices and a digest manifest; no existing package is overwritten |
| Provenance | Both binary and archive receive signed attestations; publisher verifies repository, signer workflow, master ref, source commit and hosted runner identity |
| Draft | Verify every expected archive/bundle, generate checksums and attach all seven assets to an unpublished draft |
| Publication | An authenticated maintainer verifies administrative immutability policy and all signed draft assets, seals the draft and verifies the resulting immutable release |

GitHub distinguishes signed build provenance from release integrity. The workflow uses [artifact attestations](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations) for the build and [immutable release attestations](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases) for the published tag/assets. A checksum alone does not establish either identity.

The workflow requests signing permissions only in build jobs and release-write permission only in its verified draft job. It never publishes the draft. GitHub's [immutable-release settings endpoint](https://docs.github.com/en/rest/repos/repos#check-if-immutable-releases-are-enabled-for-a-repository) requires administrative read access, which `GITHUB_TOKEN` cannot request. The authenticated maintainer checks that policy immediately before publication; no administrative token is stored in CI. A repository variable or permission-denied response is not policy proof.

Use a maintained GitHub CLI 2.x version at least 2.93 for all verification. This minimum includes the upstream fix for [authorization-header forwarding during attestation/release verification](https://github.com/cli/cli/security/advisories/GHSA-8xvp-7hj6-mcj9). The build workflow checks its installed client before any verification. No provider credential, customer environment, GitHub App key, vault, cloud resource, AI call or recovery execution belongs to publication.

## Maintainer sequence

1. Manually merge the reviewed workflow/packaging PR into `develop` and validate it. Then manually merge a separate same-repository `develop` → `master` promotion PR. Wait for Go and documentation push checks on the exact new `master` commit and verify its documentation publication.
2. Using an authenticated maintainer account with administrative read access, confirm repository release immutability is enabled and choose an unused alpha tag. Existing tags/drafts/releases cause a refusal; neither tags nor assets are overwritten. After artifact verification, the publisher atomically creates the new tag at the reviewed source SHA, verifies that target before creating the draft and again after immutable publication, and never moves an existing tag.
3. Dispatch from master, for example:

   ```sh
   gh workflow run release-cli.yml --repo Storm-Harbor/storm-harbor-control-plane \
     --ref master -f version=v0.1.0-alpha.1
   ```

4. Inspect the completed workflow, all native binary/attestation checks and the unpublished draft. CI prepares assets; successful CI is not publication. The draft must have exactly seven expected assets, the reviewed source and an unchanged tag.
5. Verify and seal the draft using the maintainer sequence below. A skipped, failed or partially completed stage is not a released version.
6. Verify a freshly downloaded package and exercise it before adoption. Keep the native delivery issue/umbrella open while publication, signature or acceptance gates remain pending.

The Go package tool is `go run ./cmd/storm-package`; it is a maintainer build tool, not a recovery command or installer. Packages contain `storm`/`storm.exe`, `build-info.txt`, `build-manifest.json`, the Apache license and dependency/toolchain notices. Their manifest records version, source revision, target, Go version and each included file's size/SHA-256, without filesystem source paths. Only the pipeline's qualified native matrix establishes platform availability.

## Maintainer verification and sealing

Run this only after the reviewed preparation workflow succeeds, with GitHub CLI 2.x >= 2.93 and the maintainer's normal authenticated account. Administrative access is used only to read repository policy; publication needs contents write. Keep credentials in the normal credential manager, never in commands, source, logs or CI secrets. The following shell example assumes GNU checksum tools; use `shasum -a 256 -c` on macOS.

```sh
set -euo pipefail
REPO=Storm-Harbor/storm-harbor-control-plane
VERSION=v0.1.0-alpha.1
SOURCE_REVISION="$(gh api "repos/$REPO/git/ref/heads/master" --jq .object.sha)"
[[ "$SOURCE_REVISION" =~ ^[a-f0-9]{40}$ ]] || exit 1
[[ "$VERSION" =~ ^v0\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)-alpha\.[1-9][0-9]*$ ]] || exit 1
STORM_VERIFY_DIR="$(mktemp -d)"
gh release view "$VERSION" --repo "$REPO" \
  --json assets,isDraft,isPrerelease,tagName,targetCommitish > "$STORM_VERIFY_DIR/draft.json"
jq -e --arg tag "$VERSION" --arg sha "$SOURCE_REVISION" \
  '.isDraft == true and .isPrerelease == true and .tagName == $tag and .targetCommitish == $sha and (.assets | length) == 7' \
  "$STORM_VERIFY_DIR/draft.json" > /dev/null
# Keep downloaded assets in their own directory; metadata is not an asset.
mkdir "$STORM_VERIFY_DIR/assets"
gh release download "$VERSION" --repo "$REPO" --dir "$STORM_VERIFY_DIR/assets"
cd "$STORM_VERIFY_DIR/assets"
expected=("storm_${VERSION}_linux_amd64.tar.gz" "storm_${VERSION}_darwin_arm64.tar.gz" "storm_${VERSION}_windows_amd64.zip")
bundles=("storm_${VERSION}_linux_amd64.sigstore.jsonl" "storm_${VERSION}_darwin_arm64.sigstore.jsonl" "storm_${VERSION}_windows_amd64.sigstore.jsonl")
test "$(find . -mindepth 1 -maxdepth 1 -type f | wc -l)" -eq 7
for index in 0 1 2; do
  test -f "${expected[$index]}" && test ! -L "${expected[$index]}"
  test -f "${bundles[$index]}" && test ! -L "${bundles[$index]}"
  gh attestation verify "${expected[$index]}" --repo "$REPO" \
    --bundle "${bundles[$index]}" --signer-workflow "$REPO/.github/workflows/release-cli.yml" \
    --source-ref refs/heads/master --source-digest "$SOURCE_REVISION" --deny-self-hosted-runners
done
test -f SHA256SUMS && test ! -L SHA256SUMS
sha256sum --check SHA256SUMS
# Administrative policy read is a maintainer action, outside GITHUB_TOKEN CI.
test "$(gh api "repos/$REPO/immutable-releases" --jq .enabled)" = true
test "$(gh api "repos/$REPO/git/ref/tags/$VERSION" --jq .object.sha)" = "$SOURCE_REVISION"
gh release edit "$VERSION" --repo "$REPO" --draft=false --latest=false
test "$(gh release view "$VERSION" --repo "$REPO" --json isImmutable --jq .isImmutable)" = true
test "$(gh api "repos/$REPO/git/ref/tags/$VERSION" --jq .object.sha)" = "$SOURCE_REVISION"
gh release verify "$VERSION" --repo "$REPO"
for asset in "${expected[@]}" "${bundles[@]}" SHA256SUMS; do
  gh release verify-asset "$VERSION" "$asset" --repo "$REPO"
done
```

If policy access is denied/disabled, source/tag differs, assets are missing or signatures fail, stop with the draft unpublished. Do not ignore errors, publish a mutable release or broaden CI credentials. After sealing, corrections require a new version. Independently download/inspect the published package and run executable acceptance before adopter workflows consume it.

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
  --source-ref refs/heads/master --source-digest "$SOURCE_REVISION" \
  --deny-self-hosted-runners

mkdir "storm-$VERSION"
tar -xzf "$ARCHIVE" -C "storm-$VERSION"
"./storm-$VERSION/storm" version --output json
```

The reported version/revision must match the verified release. Use [offline validation](../getting-started/offline-cli.md) with independently supplied local expected scope. Verification does not authorize a cloud action or prove recovery readiness.

## Failure and recovery

Before a draft exists, a failed attempt publishes no release; a new explicit dispatch can retry only if the version remains unused. If a draft/tag or partial upload exists, stop and inspect its source and assets. The workflow deliberately does not delete drafts, move tags, overwrite assets or silently resume unknown publication outcomes. Keep any incomplete draft unpublished until its exact provenance and asset set are reconciled under maintainer review.

After publication, verify the immutable release and its assets. Corrections require a new version. Private adopter repository creation and validation-workflow adoption follow a published/qualified CLI; they do not run as part of this release workflow. Canonical semantic normalization, cloud execution, GitHub App installation, AI adapters, dashboards and the future monitoring application retain their separate delivery gates.
