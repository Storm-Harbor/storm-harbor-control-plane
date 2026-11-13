# Offline Go CLI

**Source implementation merged; binary installation requires a verified release.** The first CLI is tracked by [delivery #14](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/14) and [adoption documentation #13](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/13), under [Journey #11](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/11) and [Epic #12](https://github.com/Storm-Harbor/storm-harbor-control-plane/issues/12). Use the [release/verification guide](../operations/releases.md) and [actual release list](https://github.com/Storm-Harbor/storm-harbor-control-plane/releases) to establish binary availability. Manual merge or a temporary CI artifact does not substitute for immutable publication, signed provenance and native acceptance.

The primary repository is `storm-harbor-control-plane`; the executable is `storm`. It consumes immutable `storm-lib-recovery-contracts v0.1.0-alpha.1` without copying its source. Product website and future monitoring application remain separate.

## Build the reviewed source candidate

Use Go **1.26.11**, review the owning PR/revision and build from this repository:

```sh
go build -trimpath -o storm ./cmd/storm
./storm version
./storm help
```

On Windows choose `-o storm.exe` and run `./storm.exe`. Default build metadata is `dev`/`unknown`; CI records the exact source revision separately. Candidate workflow artifacts are temporary CI outputs, not signed immutable releases. Only claim a supported platform after native executable acceptance. Cross-compilation alone is insufficient.

## Validate a synthetic declaration

```sh
./storm validate --file examples/exercise.yaml \
  --organization org-example --project project-example

./storm validate --file examples/restore.yaml \
  --organization org-example --project project-example --output json
```

Both examples use the [published YAML contract](yaml.md). All three input flags are required. Choose the expected organization/project independently of the file; copying scope from YAML only proves self-consistency. Even an independently selected **local expectation is not authenticated identity, membership or online authorization**. The CLI makes no cloud, vault, GitHub or AI request and does not write configuration or execute recovery.

Only a regular input file is accepted. Symbolic links as input files, directories, special files and stdin are unsupported. The read limit is **64 KiB**, with the producer's strict depth/type/field rules. The operation has a five-second deadline; interruption cancels the caller's wait. A stalled filesystem operation cannot be forcibly interrupted portably, so the worker may remain blocked until the OS returns or the process exits; this CLI is not a persistent file-processing service.

Flags support `--flag value` and `--flag=value`. Unknown or duplicate flags, positional extras and unsupported commands fail safely. There is no `plan`, `run`, `resume`, `cancel` or secret-management command in this slice.

## Structured results

`--output json` selects `cli.stormharbor.io/v1alpha1`. Success goes to stdout; errors go to stderr with a nonzero exit code. The success envelope contains `command`, `status`, `boundary`, `configurationApiVersion` and `inputSHA256`. It omits supplied names, scope IDs, secret aliases, paths and document contents.

`inputSHA256` hashes **exact input bytes**. Formatting or comments change it. It is not canonical semantic normalization, an approved-plan digest, an authentication proof or permission to disclose the input. Future shared normalization/planning semantics need their own published producer contracts.

Errors expose only a stable code, module, category, localization key and allowlisted safe arguments. Recovery-contract violations retain the producer's `contracts.*` keys. CLI boundary errors use `cli.*` keys. Human-readable presentation is English in this first technical CLI; no locale selector is implemented. A `SCOPE_MISMATCH` category describes a local comparison, not an authoritative access-control decision. Parser excerpts, raw OS errors and supplied flag/file values are not echoed.

| Exit | Meaning |
| --- | --- |
| `0` | Valid against local expectation, version or help |
| `2` | Invalid configuration/scope or oversized document |
| `3` | Unreadable or unsupported input file |
| `64` | Unsupported command or invalid arguments |
| `70` | Safely handled unexpected internal failure |
| `74` | Output could not be written |
| `124` | Validation deadline expired |
| `130` | Operation interrupted/cancelled |

## Qualification and next steps

`make check-go` checks formatting, vet, one combined race/statement-coverage run, actual executable scenarios, build and module integrity. The coverage threshold is 97% of native statements; independent branch/line measurements are not claimed. CI runs native executable acceptance on Linux, macOS and Windows, and qualifies only the actual OS/architecture observed in those jobs. Artifacts include build metadata and the source revision; immutable release checksums, signed provenance and installation are separate publication gates.

After the reviewed CLI release is available, prepare the approved private customer recovery repository and its credential-free validation workflow. The GitHub App, opt-in agent adapters, external bindings, orchestration, dashboards and live exercises remain later deliveries in the [master plan](../SPECS-CONTROL-PLANE.md).
